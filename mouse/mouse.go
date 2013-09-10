// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// /dev/input/mouseN handling.
//
//=====================================================================================================================

package mouse

import (
	"os"
	"bytes"
	"encoding/binary"
	"fbdev"
	log "github.com/cihub/seelog"
	"utils"
	"container/list"
	"toolkit/base"
)

type Mouse struct {
	dev				*os.File
	fb				*fbdev.Framebuffer
	xPos			int				// current X position
	yPos			int				// current Y position
	width			int		
	height			int
	cb				*list.List
	windowList		*list.List		// element list used by the composer
	mouseHndr		mouseHandler
	leftPressed		bool
	rightPressed	bool
	flags			byte
	mouseRelease	chan bool
	compositorWait	chan bool
}

type MouseMovementPacket struct {
	Flags			byte
	X				byte
	Y				byte
}

type Callback struct {
	id				uint64
	mouseHndr		mouseHandler			
	x				*int
	y				*int
	width			int
	height			int
	isMouseIn		bool		// specifies whether mouse was inside the element on previous MMP
	activateHndr	activateHandler
	flags			byte
}

const (
	F_LEFT_HOLD				= 2
	F_LEFT_HOLD_RELEASE		= 4
	F_RIGTH_HOLD			= 2
	F_RIGTH_HOLD_RELEASE	= 4
	F_LEFT_CLICK			= 8
	F_EL_ENTER				= 16
	F_EL_LEAVE				= 32
)


func Init(dev string, screenWidth int, screenHeight int, msRelease chan bool, compWait chan bool, wl *list.List) (*Mouse, error) {

	fd, err := os.OpenFile(dev, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	
	ms := &Mouse{
		dev:			fd,
		cb: 			list.New(),
		windowList:		wl,
		leftPressed:	false,
		rightPressed:	false,
		flags:			0,
		mouseRelease:	msRelease,
		compositorWait: compWait,
		xPos:			screenWidth/2,		// start in the middle
		yPos:			screenHeight/2,
		width:			screenWidth,
		height:			screenHeight,
	}

	return ms, nil
}


func (ms *Mouse) Process() (error) {

	buf := make([]byte, 3)

	mmp := MouseMovementPacket{}

	q := utils.CreateQueue(128)
	
	for {
	   	_, err := ms.dev.Read(buf)
	    if err != nil {
	        ms.dev.Close()
	        return err
	    }
		b := bytes.NewBuffer(buf)
		binary.Read(b, binary.LittleEndian, &mmp)
		q.Push(mmp)
		
		// make sure that mouse events and composer invalidation processed alternately
		<-ms.compositorWait
		
		log.Debugf("Mouse: In queue: %v", q.Count)
		
		mmp = q.Pop().(MouseMovementPacket)
	
		deltaX := int(mmp.X)
		deltaY := -int(mmp.Y)	
		
		if mmp.XSign() {
			deltaX -= 256 
		}
			
		if mmp.YSign() {
			deltaY += 256
		} 
		
		oldXpos := ms.xPos
		oldYpos := ms.yPos
		ms.xPos += deltaX
		ms.yPos += deltaY
			
		log.Debugf("Mouse Mov. Deltas: %v : %v    flags: %v", deltaX, deltaY, mmp.Flags)
		
		// TODO: handle overflow

		if mmp.BtnLeft() && !ms.leftPressed {
			ms.leftPressed = true
			ms.flags |= F_LEFT_CLICK
			log.Debug("Mouse: L btn press")
		} else if mmp.BtnLeft() && ms.leftPressed && (ms.flags & F_LEFT_HOLD) == 0 {
			ms.flags |= F_LEFT_HOLD
			ms.flags &^= F_LEFT_CLICK
			log.Debug("Mouse: L btn hold")
		} else if !mmp.BtnLeft() && ((ms.flags & F_LEFT_HOLD) != 0){	
			ms.flags &^= F_LEFT_HOLD
			ms.leftPressed = false
			log.Debug("Mouse: L btn release after hold")
		} else if !mmp.BtnLeft() && ((ms.flags & F_LEFT_CLICK) != 0){	
			ms.flags &^= F_LEFT_CLICK
			ms.leftPressed = false
			log.Debug("Mouse: L btn release")
		} 
		
		
		set := make(map[interface{}]*Callback, ms.cb.Len())
		
		// call mouse handlers whenever mouse pointer enters registered element
		for v := ms.cb.Front(); v != nil; v = v.Next() {
			cb := v.Value.(*Callback)
			
	 		if (*(cb.x) < ms.xPos && *(cb.x) + cb.width > ms.xPos && 
	 		   *(cb.y) < ms.yPos && *(cb.y) + cb.height > ms.yPos) ||
			   (*(cb.x) < oldXpos && *(cb.x) + cb.width > oldXpos &&	// accout for any element movement
			   *(cb.y) < oldYpos && *(cb.y) + cb.height > oldYpos) {

				log.Debugf("Mouse: Within element: %v : %v  -  %v : %v", *(cb.x), *(cb.y), cb.width, cb.height)

				cb.flags = ms.flags

				// pointer entered the element`s area
				if !cb.isMouseIn {	
					cb.flags |= F_EL_ENTER
					cb.isMouseIn = true
				} 
				
				// if it's a window - add it to a map so we could find later if it needs activation
				// otherwise it's a child element and we execute its mouse handler
				if cb.activateHndr == nil  {	// FIXME not a good way to test...
					cb.mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, cb.flags)
				} else {
					set[cb.id] = cb	
				}

			} else {
				// pointer left the element`s area
				if cb.isMouseIn {		
					cb.mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, ms.flags | F_EL_LEAVE)
					cb.isMouseIn = false
				}
			}
	    }
	    
	   	// find window with highest Z order 
	   	curActiveWinId := ms.windowList.Front().Value.(*base.Element).Id

	    for v := ms.windowList.Front(); v != nil; v = v.Next() {

	     	cb := set[v.Value.(*base.Element).Id]
	    
			if cb != nil {
				// execute activation handler only if it's a click and window isn't active already
				if (cb.flags & F_LEFT_CLICK) != 0 && curActiveWinId != v.Value.(*base.Element).Id {
					cb.activateHndr()
				}
				// execute the mouse handler for this window. after this we don't need to check further
				// as any of the windows at mouse location had lower Z order and we(partialy) 'hidden'
				// by the newly activated window
				cb.mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, cb.flags)
				break
			}
		}

	    // process mouse pointer element handler separately from all other windows
		ms.mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, ms.flags)

		ms.mouseRelease<- true
	}

	return nil
}


type mouseHandler func(int, int, int, int, byte)
type activateHandler func()


func (ms *Mouse) RegisterMouse(id uint64, fnMouse mouseHandler, fnActivate activateHandler, x *int, y *int, width, height int) {

	c := &Callback{
		id:				id,
		mouseHndr:		fnMouse,
		activateHndr:	fnActivate,
		x:				x,
		y:				y,
		width:			width,
		height:			height,
	}
	
	ms.cb.PushBack(c)
	
	return
}

func (ms *Mouse) DeregisterMouse(id uint64) {

	for v := ms.cb.Front(); v != nil; v = v.Next() {
	
		if v.Value.(*Callback).id == id {
			
			log.Debugf("Removing button ")
			ms.cb.Remove(v)
			break
		}
	}
}

func (ms *Mouse) RegisterMousePointer(fnMouse mouseHandler) {
	ms.mouseHndr = fnMouse
}

func (mmp *MouseMovementPacket) BtnRight() (bool) {
	return (mmp.Flags & byte(2)) != 0
}

func (mmp *MouseMovementPacket) BtnLeft() (bool) {
	return (mmp.Flags & byte(1)) != 0
}

func (mmp *MouseMovementPacket) BtnMiddle() (bool) {
	return (mmp.Flags & byte(4)) != 0					// not yet working
}

func (mmp *MouseMovementPacket) XSign() (bool) {
	return (mmp.Flags & byte(16)) != 0
}

func (mmp *MouseMovementPacket) YSign() (bool) {
	return (mmp.Flags & byte(32)) != 0
}

func (mmp *MouseMovementPacket) XOverflow() (bool) {
	return (mmp.Flags & byte(64)) != 0
}

func (mmp *MouseMovementPacket) YOverflow() (bool) {
	return (mmp.Flags & byte(128)) != 0
}


func (ms *Mouse) Close() {
	ms.dev.Close()
}