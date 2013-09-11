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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))			// FIXME

	lab.Element = base.Element{
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
	
	lab.Draw() 
	
	return lab
}


func (lab *Label) Draw() {
	log.Debug("    Drawing Label")
	
	rgba, err := fonts.Render(lab.txt, lab.Element.Width, lab.Element.Height, lab.Element.Font)
    if err == nil {
		gfx.DrawOver(&lab.parent.Element, rgba, lab.Element.X, lab.Element.Y, lab.Element.Width, lab.Element.Height)
	}
}


// mouse handler
func (lab *Label) Mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	log.Debug("Label ms handler: do nothing...")

}