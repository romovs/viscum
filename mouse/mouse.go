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
)

const(
	MAX_ELEMENTS = 128
)

type Mouse struct {
	dev				*os.File
	fb				*fbdev.Framebuffer
	
	xPos			int		// current X position
	yPos			int		// current Y position
	width			int		
	height			int
	
	Cb				[]*Callback
	CbIndex 		uint32
	
	mouseHndr		mouseHandler
	
	leftPressed		bool
	rightPressed	bool
	flags			byte
	
	mouseRelease	chan bool
	compositorWait	chan bool
}

type MouseMovementPacket struct {
	Flags		byte
	X			byte
	Y			byte
}

type Callback struct {
	mouseHndr	mouseHandler
	entity		interface{}
	x			*int
	y			*int
	width		int
	height		int
}

const (
	
	BTN_FLAG_LEFT_HOLD			= 2
	BTN_FLAG_LEFT_HOLD_RELEASE	= 4
	BTN_FLAG_RIGTH_HOLD			= 2
	BTN_FLAG_RIGTH_HOLD_RELEASE	= 4
	BTN_FLAG_LEFT_CLICK			= 8
)


func Init(dev string, fb *fbdev.Framebuffer, msRelease chan bool, compWait chan bool) (*Mouse, error) {

	fd, err := os.OpenFile(dev, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	ms := &Mouse{
		fb: 			fb,
		dev:			fd,
		Cb: 			make([]*Callback, MAX_ELEMENTS),
		CbIndex:		0,
		leftPressed:	false,
		rightPressed:	false,
		flags:			0,
		mouseRelease:	msRelease,
		compositorWait: compWait,
		xPos:			int(fb.Vinfo.Xres)/2,		// start in the middle
		yPos:			int(fb.Vinfo.Yres)/2,
		width:			int(fb.Vinfo.Xres),
		height:			int(fb.Vinfo.Yres),
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
			
		log.Debugf("Mouse Mov. Deltas: %v : %v    Flags: %v", deltaX, deltaY, mmp.Flags)
		
		// TODO: handle overflow

		if mmp.BtnLeft() && !ms.leftPressed {
			ms.leftPressed = true
			ms.flags |= BTN_FLAG_LEFT_CLICK
			log.Debug("Mouse: L btn press")
		} else if mmp.BtnLeft() && ms.leftPressed && (ms.flags & BTN_FLAG_LEFT_HOLD) == 0 {
			ms.flags |= BTN_FLAG_LEFT_HOLD
			ms.flags &^= BTN_FLAG_LEFT_CLICK
			log.Debug("Mouse: L btn hold")
		} else if !mmp.BtnLeft() && ms.leftPressed {	
			ms.flags &^= BTN_FLAG_LEFT_HOLD
			ms.leftPressed = false
			log.Debug("Mouse: L btn release")
		} 
		
		// call mouse handlers whenever mouse pointer enters registered element
		for i := 0; i < int(ms.CbIndex); i++ {
		
			if (ms.flags & BTN_FLAG_LEFT_HOLD) == 0 {
		 		if *(ms.Cb[i].x) < ms.xPos && *(ms.Cb[i].x) + ms.Cb[i].width > ms.xPos &&
				   *(ms.Cb[i].y) < ms.yPos && *(ms.Cb[i].y) + ms.Cb[i].height > ms.yPos {
	
					log.Debugf("Mouse: L HOLD: Within element: %v : %v  -  %v : %v", *(ms.Cb[i].x), *(ms.Cb[i].y), ms.Cb[i].width, ms.Cb[i].height)
	
					ms.Cb[i].mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, ms.flags)
				}
			} else {
		 		if *(ms.Cb[i].x) < oldXpos && *(ms.Cb[i].x) + ms.Cb[i].width > oldXpos &&
				   *(ms.Cb[i].y) < oldYpos && *(ms.Cb[i].y) + ms.Cb[i].height > oldYpos {
	
					log.Debugf("Mouse: Within element: %v : %v  -  %v : %v", *(ms.Cb[i].x), *(ms.Cb[i].y), ms.Cb[i].width, ms.Cb[i].height)
	
					ms.Cb[i].mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, ms.flags)
				}
			}
	    }
	    
	    // process mouse pointer element handler separately from all other windows
		ms.mouseHndr(ms.xPos, ms.yPos, deltaX, deltaY, ms.flags)

		ms.mouseRelease<- true
	}

	return nil
}


type mouseHandler func(int, int, int, int, byte)

func (ms *Mouse) RegisterMouse(fnMouse mouseHandler, x *int, y *int, width, height int) {

	c := &Callback{
		mouseHndr:	fnMouse,
		x:			x,
		y:			y,
		width:		width,
		height:		height,
	}
	
	ms.Cb[ms.CbIndex] = c
	ms.CbIndex++
	
	return
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