// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"viscum/fbdev"
	"viscum/mouse"
	log "github.com/cihub/seelog"
	"viscum/toolkit/base"
	"viscum/fonts"
)


type Label struct {
	base.Element
	parent			*Window
	fb				*fbdev.Framebuffer
	txt				string
}


func (win *Window) Label(ms *mouse.Mouse, txt string, x, y, w, h int) (*Label) {

	lab := &Label{
		parent: 	win,
		fb: 		win.fb,
		txt:		txt,
	}
	
	win.Children.PushFront(lab)

	lab.Element = base.Element{
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
	
	lab.Draw() 
	
	return lab
}


func (lab *Label) Draw() {
	fonts.Render(&lab.parent.Element, lab.txt, lab.Element.X, lab.Element.Y, lab.Element.Width, lab.Element.Height, lab.Element.Font)
}


// mouse handler
func (lab *Label) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {

	log.Debug("Label ms handler: do nothing...")

}
