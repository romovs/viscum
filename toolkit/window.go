// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	log "github.com/cihub/seelog"
	"gfx"
	"math/rand"
	"time"
	"toolkit/base"
)


type Window struct {
	base.Element
	fb				*fbdev.Framebuffer
	wasClicked		bool
	titleBarHeight	int
	Children		[]base.IElement
	ChildrenCnt		int
	cmpWinActHndr	cmpWinActivateHandler
}

type cmpWinActivateHandler func(uint64)


func CreateWindow (fnCmpWinActivate cmpWinActivateHandler, fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, 
					x, y, w, h int) (*Window, error) {

	win := &Window{
		fb: 				fb,
		titleBarHeight: 	20,
		wasClicked: 		false,
		Children: 			make([]base.IElement, 128),
		cmpWinActHndr:		fnCmpWinActivate,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))			// FIXME

	win.Element = base.Element{
		Id:				uint64(r.Int63()),
		X: 				x,				// relative position within the parent element
		Y: 				y,
		ScreenX:		x,				// position within the screen coordinates
		ScreenY:		y,
		Width: 			w,
		Height:			h,
		Buffer: 		make([]byte, w*h*4),
		InvMsgPipe: 	imp,
		DeactivateHndr: win.Deactivate,
	}

	ms.RegisterMouse(win.Element.Id, win.Mouse, win.activate, &win.Element.ScreenX, &win.Element.ScreenY, w, h)
	
	win.Draw()

	return win, nil
}


func (win *Window) Draw() {
	log.Debug("Drawing Window")

	// window
	gfx.RectFilled(win.Element.Buffer, 0, 0, win.Element.Width, win.Element.Height, win.Element.Width, 255, 0, 0, 0)	
	gfx.Rect(win.Element.Buffer, 0, 0, win.Element.Width-1, win.Element.Height-1, win.Element.Width, 0, 0, 0, 0)	

	// title bar
	gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 29, 59, 99, 0)	

	// render children
	for i := 0; i < win.ChildrenCnt; i++ {
		win.Children[i].Draw()
	}
}


func (win *Window) activate() {
	gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width,  55, 109, 181, 0)	
	win.cmpWinActHndr(win.Id)
}


func (win *Window) Deactivate() {
	gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 29, 59, 99, 0)	
}


// mouse handler
func (win *Window) Mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	// drag only if clicked inside titlebar. it's enoguh to check Y position. because X will be anyway inside the window bounds
	if win.Element.Y + deltaY <= y && y <= win.Element.Y + deltaY + win.titleBarHeight{
	
		if (flags & mouse.F_LEFT_CLICK) != 0 {
			log.Debug("Window ms handler: title F_LEFT_CLICK")
			win.wasClicked = true
		} else if win.wasClicked && (flags & mouse.F_LEFT_HOLD) != 0 {
			log.Debug("Window ms handler: title F_LEFT_HOLD")
			win.Element.X += deltaX
			win.Element.Y += deltaY
			win.Element.ScreenX += deltaX
			win.Element.ScreenY += deltaY
			
			// update screen position for all the children
			for i := 0; i < win.ChildrenCnt; i++ {
				win.Children[i].UpdateScreenX(deltaX)
				win.Children[i].UpdateScreenY(deltaY)
			}
		} else {
			log.Debug("Window ms handler: title do nothing...")
			win.wasClicked = false
		}
	}
}