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
	Children		[]IElement
	ChildrenCnt		int
}


func (win *Window) Init(fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, x, y, w, h, z int) (error) {

	win.fb = fb
	win.titleBarHeight = 20
	win.wasClicked = false
	win.Children = make([]IElement, 128)

	win.Element = Element{
		X: 			x,				// relative position within the parent element
		Y: 			y,
		ScreenX:	x,				// position within the screen coordinates
		ScreenY:	y,
		Width: 		w,
		Height:		h,
		Buffer: 	make([]byte, w*h*4),
		InvMsgPipe: imp,
	}
	
	ms.RegisterMouse(win.mouse, &win.Element.ScreenX, &win.Element.ScreenY, w, h)
	
	win.Draw()

	return nil
}


func (win *Window) Draw() {
	log.Debug("Drawing Window")

	// window
	gfx.RectFilled(win.Element.Buffer, 0, 0, win.Element.Width, win.Element.Height, win.Element.Width, 255, 0, 0, 0)	
	gfx.Rect(win.Element.Buffer, 0, 0, win.Element.Width-1, win.Element.Height-1, win.Element.Width, 0, 0, 0, 0)	

	// title bar
	gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 0, 0, 255, 0)	

	// render children
	for i := 0; i < win.ChildrenCnt; i++ {
		win.Children[i].Draw()
	}
}


// mouse handler
func (win *Window) mouse(x int, y int, deltaX int, deltaY int, flags byte) {

	// drag only if clicked inside titlebar. it's enoguh to check Y position. because X will be anyway inside the window bounds
	if win.Element.Y + deltaY <= y && y <= win.Element.Y + deltaY + win.titleBarHeight{
	
		if (flags & mouse.F_LEFT_CLICK) != 0 {
			log.Debug("Window ms handler: F_LEFT_CLICK")
			win.wasClicked = true
		} else if win.wasClicked && (flags & mouse.F_LEFT_HOLD) != 0 {
			log.Debug("Window ms handler: F_LEFT_HOLD")
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
			log.Debug("Window ms handler: do nothing...")
			win.wasClicked = false
		}
	}
}