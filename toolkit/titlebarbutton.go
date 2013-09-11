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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))			// FIXME

	but.Element = base.Element{
		Id:			uint64(r.Int63()),
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
	
	gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width, 80, 130, 0, 0)	
	gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, 0)	
}


// mouse handler
func (but *TitleBarButton) Mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	if (flags & mouse.F_LEFT_CLICK) != 0 {
		log.Debug("TitleBarButton ms handler: clicked inside.")
		but.wasClicked = true
		gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width,  49, 80, 0, 0)	
		gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, 0)	
 	} else if but.wasClicked && (flags & mouse.F_LEFT_CLICK) == 0 && (flags & mouse.F_LEFT_HOLD) == 0 {
		log.Debug("TitleBarButton ms handler: clicked & released inside.")
		but.wasClicked = false
		but.Draw()
		but.clickHndr()
	} else if but.wasClicked && (flags & mouse.F_EL_LEAVE) != 0 {
		log.Debug("TitleBarButton ms handler: clicked inside. released outside.")
		but.wasClicked = false
		but.Draw()
	} else {
		log.Debug("TitleBarButton ms handler: do nothing...")
	}
}