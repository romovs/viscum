// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	log "github.com/cihub/seelog"
	"gfx"
)

type clickHandler func()

type Button struct {
	Element
	parent			*Window
	fb				*fbdev.Framebuffer
	wasClicked		bool
	Click			clickHandler
}


func (win *Window) Button(ms *mouse.Mouse, clickHnd clickHandler, x, y, w, h, z int) (*Button) {

	but := &Button{
		parent: 	win,
		fb: 		win.fb,
		wasClicked: false,
		Click:		clickHnd,
	}
	
	win.Children[win.ChildrenCnt] = but
	win.ChildrenCnt++

	but.Element = Element{
		X: 			x,
		Y: 			y,
		ScreenX:	win.Element.X+x,
		ScreenY:	win.Element.Y+y,
		Width: 		w,
		Height:		h,
		InvMsgPipe: win.InvMsgPipe,
	}
	
	ms.RegisterMouse(but.mouse, &but.Element.ScreenX, &but.Element.ScreenY, w, h)
	
	win.Draw()
	
	return but
}


func (but *Button) Draw() {
	log.Debug("    Drawing Button")
	
	gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width, 80, 130, 0, 0)	
	gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, 0)	
}


// mouse handler
func (but *Button) mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	if (flags & mouse.BTN_FLAG_LEFT_CLICK) != 0 {
		log.Debug("Button ms handler: BTN_FLAG_LEFT_CLICK")
		but.wasClicked = true
		// visualise the click
		gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width,  49, 80, 0, 0)	
		gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, 0)	
 	} else if but.wasClicked && (flags & mouse.BTN_FLAG_LEFT_CLICK) == 0 && (flags & mouse.BTN_FLAG_LEFT_HOLD) == 0 {
		log.Debug("Button ms handler: L btn release")
		but.wasClicked = false
		gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width, 80, 130, 0, 0)	
		gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, 0)	
		but.Click()
	} else {
		log.Debug("Button ms handler: do nothing...")
		but.wasClicked = false
	}
}