// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	log "github.com/cihub/seelog"
	"gfx"
	"toolkit/base"
)


type TitleBarButton struct {
	base.Element
	parent			*Window
	fb				*fbdev.Framebuffer
	wasClicked		bool
	clickHndr		clickHandler
}


func (win *Window) TitleBarButton(ms *mouse.Mouse, fnClick clickHandler, x, y, w, h int) (*TitleBarButton) {

	but := &TitleBarButton{
		parent: 	win,
		fb: 		win.fb,
		wasClicked: false,
		clickHndr:	fnClick,
	}
	
	win.Children.PushFront(but)

	but.Element = base.Element{
		Id:			base.GenerateId(),
		X: 			x,
		Y: 			y,
		ScreenX:	win.Element.X+x,
		ScreenY:	win.Element.Y+y,
		Width: 		w,
		Height:		h,
		InvMsgPipe: win.InvMsgPipe,
	}
	
	ms.RegisterMouse(but.Element.Id, but.Mouse, nil, &but.Element.ScreenX, &but.Element.ScreenY, w, h)
	
	but.Draw()
	
	return but
}


func (but *TitleBarButton) Draw() {
	log.Debug("    Drawing TitleBarButton")
	
	gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width, 80, 130, 0, gfx.A_OPAQUE)	
	gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	
}


// mouse handler
func (but *TitleBarButton) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {

	if (flags & mouse.F_L_CLICK) != 0 {
		log.Debug("TitleBarButton ms handler: click")
		but.wasClicked = true
		gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width,  49, 80, 0, gfx.A_OPAQUE)	
		gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	
 	} else if but.wasClicked && (flags & mouse.F_L_RELEASE) != 0 {
		log.Debug("TitleBarButton ms handler: release")
		but.wasClicked = false
		but.Draw()
		but.clickHndr(false)
	}
}