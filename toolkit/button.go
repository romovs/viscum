// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	log "github.com/cihub/seelog"
	"gfx"
	"toolkit/base"
	"fonts"
	"image"
)

type Button struct {
	base.Element
	parent			*Window
	fb				*fbdev.Framebuffer
	wasClicked		bool
	clickHndr		clickHandler
	txt				string
	style			byte
	icon 			image.Image
}

const(
	BS_TEXT			= 1 << iota
	BS_ICON			= 1 << iota
	BS_ICON_TEXT	= 1 << iota
	
	BS_TOGGLE		= 1 << iota
)

const(
	ICON_TEXT_SPACE	= 5		// horizontal space in pixels between icon and text for BS_ICON_TEXT buttons
)

// style should be one of BS_TEXT, BS_ICON_TEXT, BS_ICON
// icon should be nil for BS_TEXT
// txt should be empty string for BS_ICON
func (win *Window) Button(style byte, ms *mouse.Mouse, icon image.Image, txt string, fnClick clickHandler, x, y, w, h int) (*Button) {

	but := &Button{
		parent: 	win,
		fb: 		win.fb,
		wasClicked: false,
		clickHndr:	fnClick,
		txt:		txt,
		style:		style,
		icon:		icon,
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
		Font:		fonts.Default(),
	}
	
	ms.RegisterMouse(but.Element.Id, but.Mouse, nil, &but.Element.ScreenX, &but.Element.ScreenY, w, h)
	
	but.Draw(false)
	
	return but
}


func (but *Button) Draw(isPushed interface{}) {
	var r, g, b byte
	
	if isPushed.(bool) {
		r, g, b = 49, 80, 0
	} else {
		r, g, b = 80, 130, 0
	}

	gfx.RectFilled(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width, but.Element.Y+but.Element.Height, but.parent.Element.Width,  r, g, b, gfx.A_OPAQUE)	
	gfx.Rect(but.parent.Element.Buffer, but.Element.X, but.Element.Y, but.Element.X+but.Element.Width-1, but.Element.Y+but.Element.Height-1, but.parent.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	

	if but.style & BS_TEXT != 0 {
		w, h, _ := fonts.ExpectedSize(but.Element.Font, but.txt)
		hoffsetTxt := but.Element.Width/2 - int(w)/2 
		voffsetTxt := but.Element.Height/2 - int(h)/2 + 1
		fonts.Render(&but.parent.Element, but.txt, 
					 but.Element.X + hoffsetTxt, but.Element.Y + voffsetTxt, 
					 but.Element.Width, but.Element.Height, but.Element.Font)
	} else if but.style & BS_ICON_TEXT != 0 {
		iw := but.icon.Bounds().Dx()
		ih := but.icon.Bounds().Dy()
		w, h, _ := fonts.ExpectedSize(but.Element.Font, but.txt)
		hoffsetIcon := but.Element.Width/2 - int(int(w)+iw+ICON_TEXT_SPACE)/2 
		voffsetTxt := but.Element.Height/2 - int(h)/2 + 1
		voffsetIcon := but.Element.Height/2 - int(ih)/2
		gfx.DrawOver(&but.parent.Element, but.icon, but.Element.X + hoffsetIcon, but.Element.Y + voffsetIcon, iw, ih)
		fonts.Render(&but.parent.Element, but.txt, 
					 but.Element.X + hoffsetIcon + iw + ICON_TEXT_SPACE, but.Element.Y + voffsetTxt, 
					 but.Element.Width, but.Element.Height, but.Element.Font)
	} else if but.style & BS_ICON != 0 {
		iw := but.icon.Bounds().Dx()
		ih := but.icon.Bounds().Dy()
		hoffsetIcon := but.Element.Width/2 - int(iw)/2 
		voffsetIcon := but.Element.Height/2 - int(ih)/2
		gfx.DrawOver(&but.parent.Element, but.icon, but.Element.X + hoffsetIcon, but.Element.Y + voffsetIcon, iw, ih)
	}
}


// mouse handler
func (but *Button) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {

	if (flags & mouse.F_L_CLICK) != 0 {
		log.Debug("Button ms handler: click")
		but.wasClicked = true
		but.Draw(true)
 	} else if but.wasClicked && (flags & mouse.F_L_RELEASE) != 0 {
		log.Debug("Button ms handler: release")
		but.wasClicked = false
		but.Draw(false)
		but.clickHndr()
	} else if but.wasClicked && (flags & mouse.F_EL_LEAVE) != 0 {
		// release the button if user clicked inside it and then dragged the mouse outside without releasing the mouse button
		log.Debug("Button ms handler: clicked inside. released outside.")
		but.wasClicked = false
		but.Draw(false)
	} 
}