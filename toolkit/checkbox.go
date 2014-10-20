// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is districbed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"viscum/fbdev"
	"viscum/mouse"
	"viscum/gfx"
	"viscum/toolkit/base"
	"viscum/fonts"
	"image"
)

type CheckBox struct {
	base.Element
	parent			*Window
	fb				*fbdev.Framebuffer
	clickHndr		clickHandler
	txt				string
	style			byte
	checked			bool
	mouseIn			bool
}


const(
	CBS_LEFTTEXT	= 1 << iota
)

const(
	CB_TEXT_SPACE	= 5		// horizontal space in pixels between the checkbox and the text
)


func (win *Window) CheckBox(style byte, ms *mouse.Mouse, txt string, fnClick clickHandler, x, y int) (*CheckBox) {

	cb := &CheckBox{
		parent: 	win,
		fb: 		win.fb,
		clickHndr:	fnClick,
		txt:		txt,
		style:		style,
		checked:	false,
		mouseIn:	false,
	}
	
	win.Children.PushFront(cb)

	cb.Element = base.Element{
		Id:			base.GenerateId(),
		X: 			x,
		Y: 			y,
		ScreenX:	win.Element.X+x,
		ScreenY:	win.Element.Y+y,
		Width: 		tk.cbChecked.Bounds().Dx(),
		Height:		tk.cbChecked.Bounds().Dy(),
		InvMsgPipe: win.InvMsgPipe,
		Font:		fonts.Default(),
	}
	
	var totalWidth, totalHeight int
	
	if txt != "" {
		w, h, _ := fonts.ExpectedSize(cb.Element.Font, txt)
		totalWidth = cb.Element.Width + ICON_TEXT_SPACE + int(w)
			
		if int(h) > cb.Element.Height {
			totalHeight = int(h)
		} else {
			totalHeight = cb.Element.Height
		}
	} else {
		totalWidth = cb.Element.Width
		totalHeight = cb.Element.Height
	}
	
	ms.RegisterMouse(cb.Element.Id, cb.Mouse, nil, &cb.Element.ScreenX, &cb.Element.ScreenY, totalWidth, totalHeight)
	
	cb.Draw()
	
	return cb
}


func (cb *CheckBox) Draw() {
	var img image.Image

	if cb.checked && cb.mouseIn {
		img = tk.cbCheckedHover
	} else if cb.checked && !cb.mouseIn {
		img = tk.cbChecked
	} else if !cb.checked && cb.mouseIn {
		img = tk.cbUncheckedHover
	} else if !cb.checked && !cb.mouseIn {
		img = tk.cbUnchecked
	}

	gfx.DrawOver(&cb.parent.Element, img, cb.Element.X, cb.Element.Y, cb.Element.Width, cb.Element.Height)

	if cb.txt != "" {
		w, h, _ := fonts.ExpectedSize(cb.Element.Font, cb.txt)
		voffsetTxt := cb.Element.Height/2 - int(h)/2 + 1
		
		// FIXME: temp solution. instead, parent element should expose region redraw method
		bg := cb.parent.GetBG()
		gfx.RectFilled(cb.parent.Element.Buffer, 
						cb.Element.X + cb.Element.Width + ICON_TEXT_SPACE, 
						cb.Element.Y + voffsetTxt, 
						cb.Element.X + cb.Element.Width + ICON_TEXT_SPACE + int(w),
						cb.Element.Y + voffsetTxt + int(h), 
						cb.parent.Element.Width,
						bg.R, bg.G, bg.B, bg.A)
	
		fonts.Render(&cb.parent.Element, cb.txt, 
					 cb.Element.X + cb.Element.Width + ICON_TEXT_SPACE, cb.Element.Y + voffsetTxt, 
					 cb.parent.Element.Width - cb.Element.X - ICON_TEXT_SPACE - cb.Element.Width, cb.Element.Height, 
					 cb.Element.Font)
	}
}


func (cb *CheckBox) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {
	if (flags & mouse.F_L_CLICK) != 0 {
		cb.checked = !cb.checked
		cb.Draw()
		cb.clickHndr(cb.checked)		
	} else if (flags & mouse.F_EL_ENTER) != 0 {
		cb.mouseIn = true
		cb.Draw()
	} else if (flags & mouse.F_EL_LEAVE) != 0 {
		cb.mouseIn = false
		cb.Draw()
	}
}
