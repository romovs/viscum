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
	"fonts"
)

type Button struct {
	base.Element
	parent			*Window
	fb				*fbdev.Framebuffer
	wasClicked		bool
	clickHndr		clickHandler
	txt				string
}


func (win *Window) Button(ms *mouse.Mouse, txt string, fnClick clickHandler, x, y, w, h int) (*Button) {

	but := &Button{
		parent: 	win,
		fb: 		win.fb,
		wasClicked: false,
		clickHndr:	fnClick,
		txt:		txt,
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
		Font:		fonts.Default(),
	}
	
	ms.RegisterMouse(but.Element.Id, but.Mouse, nil, &but.Element.ScreenX, &but.Element.ScreenY, w, h)
	
	but.Draw()
	
	return but
}


func (but *Button) Draw() {
	log.Debug("    Drawing Button")
	
	gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width, 80, 130, 0, gfx.A_OPAQUE)	
	gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	

	// since there is no layout manager, just reneder the text starting at (0,0) for now..
	fonts.Render(&but.parent.Element, but.txt, but.Element.X, but.Element.Y, but.Element.Width, but.Element.Height, but.Element.Font)
}


// mouse handler
func (but *Button) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {

	if (flags & mouse.F_L_CLICK) != 0 {
		log.Debug("Button ms handler: click")
		but.wasClicked = true
		// visualise the click
		gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width,  49, 80, 0, gfx.A_OPAQUE)	
		gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	
		fonts.Render(&but.parent.Element, but.txt, but.Element.X, but.Element.Y, but.Element.Width, but.Element.Height, but.Element.Font)
 	} else if but.wasClicked && (flags & mouse.F_L_RELEASE) != 0 {
		log.Debug("Button ms handler: release")
		but.wasClicked = false
		but.Draw()
		but.clickHndr()
	} else if but.wasClicked && (flags & mouse.F_EL_LEAVE) != 0 {
		// release the button if user clicked inside it and then dragged the mouse outside without releasing the mouse button
		log.Debug("Button ms handler: clicked inside. released outside.")
		but.wasClicked = false
		but.Draw()
	} 
}