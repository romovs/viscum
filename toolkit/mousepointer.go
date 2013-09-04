// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	"gfx"
	"time"
)


type MousePointer struct {
	Element
	fb			*fbdev.Framebuffer
}

const (
	WIDTH = 8
)


func (app *MousePointer) Init(fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64) (error) {

	app.fb = fb

	app.Element = Element{
		X:				int(fb.Vinfo.Xres/2),
		Y: 				int(fb.Vinfo.Yres/2),
		Width: 			WIDTH,
		Height:			WIDTH,
		Buffer: 		make([]byte, 8*8*4),
		InvMsgPipe:		imp,
	}
	
	ms.RegisterMousePointer(app.mouse)

	gfx.SetPixel(app.Element.Buffer, 0, 0, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 1, 0, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 2, 0, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 3, 0, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 0, 1, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 1, 1, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 2, 1, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 0, 2, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 1, 2, app.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(app.Element.Buffer, 0, 3, app.Element.Width, 0, 255, 0, 0)
	return nil
}


func (app *MousePointer) mouse(x int, y int, deltaX int, deltaY int, flags byte) {
	app.Element.X += deltaX
	app.Element.Y += deltaY
	app.Element.InvMsgPipe <- time.Now().UnixNano()
}
