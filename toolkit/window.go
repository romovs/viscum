// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	log "github.com/cihub/seelog"
	"gfx"
)


type Window struct {
	Element
	fb				*fbdev.Framebuffer
	wasClicked		bool
	titleBarHeight	int
}


func (win *Window) Init(fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, x, y, w, h, z int) (error) {

	win.fb = fb
	win.titleBarHeight = 20
	win.wasClicked = false

	win.Element = Element{
		X: 			x,
		Y: 			y,
		Width: 		w,
		Height:		h,
		Buffer: 	make([]byte, 1024*768*4),
		InvMsgPipe: imp,
	}
	
	ms.RegisterMouse(win.mouse, &win.Element.X, &win.Element.Y, w, h)
	
	gfx.RectFilled(win.Element.Buffer, 0, 0, win.Element.Width, win.Element.Height, win.Element.Width, 255, 0, 0, 0)	
	gfx.Rect(win.Element.Buffer, 0, 0, win.Element.Width-1, win.Element.Height-1, win.Element.Width, 0, 0, 0, 0)	

	// title bar
	gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 0, 0, 255, 0)	

	return nil
}


// mouse handler
func (win *Window) mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	log.Debugf("Click-hold at: [%v:%v]. Window at %v:%v - %v:%v", x, y, win.Element.X, win.Element.Y, win.Element.Width+win.Element.X, win.Element.Height+win.Element.Y)
	
	// drag only if clicked inside titlebar. it's enoguh to check Y position. because X will be anyway inside the window bounds
	if win.Element.Y + deltaY <= y && y <= win.Element.Y + deltaY + win.titleBarHeight{
	
		if (flags & mouse.BTN_FLAG_LEFT_CLICK) != 0 {
			log.Debug("win BTN_FLAG_LEFT_CLICK")
			win.wasClicked = true
		} else if win.wasClicked && (flags & mouse.BTN_FLAG_LEFT_HOLD) != 0 {
			log.Debug("win BTN_FLAG_LEFT_HOLD")
			win.Element.X += deltaX
			win.Element.Y += deltaY
		} else {
			log.Debug("win")
			win.wasClicked = false
		}
	}
}