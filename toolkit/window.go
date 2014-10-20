// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"viscum/fbdev"
	"viscum/mouse"
	log "github.com/cihub/seelog"
	"viscum/gfx"
	"viscum/toolkit/base"
	"container/list"
	"image/color"
)


type Window struct {
	base.Element
	fb				*fbdev.Framebuffer
	wasClicked		bool
	titleBarHeight	int
	cmpWinActHndr	cmpWinActivateHandler
	closeButton		*TitleBarButton
	tbHidden		bool
	bg				color.RGBA
}

const(
	WS_TITLEBAR_HIDDEN	= 1 << iota
)


func (win  *Window) GetBG() (color.RGBA) {
	return win.bg
}

type cmpWinActivateHandler func(uint64)


func CreateWindow (style byte, fnCmpWinActivate cmpWinActivateHandler, fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, 
					x, y, w, h int) (*Window, error) {

	win := &Window{
		fb: 				fb,
		titleBarHeight: 	20,
		wasClicked: 		false,
		cmpWinActHndr:		fnCmpWinActivate,
		tbHidden:			style & WS_TITLEBAR_HIDDEN != 0,
		bg:					color.RGBA{
			R:	241, 
			G:	240, 
			B:	238,
			A:	gfx.A_OPAQUE,
		},
	}
	
	win.Element = base.Element{
		Id:				base.GenerateId(),
		X: 				x,				// relative position within the parent element
		Y: 				y,
		ScreenX:		x,				// position within the screen coordinates
		ScreenY:		y,
		Width: 			w,
		Height:			h,
		Buffer: 		make([]byte, w*h*4),
		InvMsgPipe: 	imp,
		DeactivateHndr: win.Deactivate,
		Children: 		list.New(),
	}

	ms.RegisterMouse(win.Element.Id, win.Mouse, win.activate, &win.Element.ScreenX, &win.Element.ScreenY, w, h)
	
	if !win.tbHidden {
		win.closeButton = win.TitleBarButton(ms, func(_ bool) {
			log.Debugf("Window %v exiting...", win.Id)
			// deregister the window
			win.Element.CompRemoveHdnr(win.Id)
			// remove mouse handlers
			ms.DeregisterMouse(win.Id)
			for v := win.Children.Front(); v != nil; v = v.Next() {
				ms.DeregisterMouse(v.Value.(base.IElement).GetId())
			}	
		}, w-15, 5, 10, 10)
	}

	win.Draw()
	
	return win, nil
}


func (win *Window) Draw() {
	log.Debugf("Drawing Window %v", win.Id)

	// window
	gfx.RectFilled(win.Element.Buffer, 0, 0, win.Element.Width, win.Element.Height, win.Element.Width, win.bg.R, win.bg.G, win.bg.B, win.bg.A)	
	gfx.Rect(win.Element.Buffer, 0, 0, win.Element.Width-1, win.Element.Height-1, win.Element.Width, 0, 0, 0, gfx.A_OPAQUE)	

	// title bar
	if !win.tbHidden {
		gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 29, 59, 99, gfx.A_OPAQUE)	
	}
	// render children
	for v := win.Children.Front(); v != nil; v = v.Next() {
		v.Value.(base.IElement).Draw()
	}	
}


func (win *Window) activate() {
	win.cmpWinActHndr(win.Id)

	if !win.tbHidden {
		gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 55, 109, 181, gfx.A_OPAQUE)	
		// redraw titlebar buttons
		win.closeButton.Draw()
	}
}


func (win *Window) Deactivate() {
	if !win.tbHidden {
		gfx.RectFilled(win.Element.Buffer, 1, 1, win.Element.Width-1, win.titleBarHeight, win.Element.Width, 29, 59, 99, gfx.A_OPAQUE)	
		// redraw titlebar buttons
		win.closeButton.Draw()
	}
}


// mouse handler
func (win *Window) Mouse(x int, y int, deltaX int, deltaY int, flags uint16) {

	// drag only if clicked inside titlebar. checking Y position is enough, because X will be inside the window bounds anyway
	if !win.tbHidden && win.Element.Y + deltaY <= y && y <= win.Element.Y + deltaY + win.titleBarHeight{
	
		if (flags & mouse.F_L_CLICK) != 0 {
			log.Debug("Window ms handler: click")
			win.wasClicked = true
		} else if win.wasClicked && (flags & mouse.F_L_HOLD) != 0 {
			log.Debug("Window ms handler: drag")
			win.Element.X += deltaX
			win.Element.Y += deltaY
			win.Element.ScreenX += deltaX
			win.Element.ScreenY += deltaY
			
			// update screen position for all the children
			for v := win.Children.Front(); v != nil; v = v.Next() {
				v.Value.(base.IElement).UpdateScreenX(deltaX)
				v.Value.(base.IElement).UpdateScreenY(deltaY)
			}	
		} else if (flags & mouse.F_L_DBL_CLICK) != 0 {
			log.Debug("Window ms handler: double click")
		} else {
			log.Debug("Window ms handler: title do nothing...")
			win.wasClicked = false
		}
	}
}
