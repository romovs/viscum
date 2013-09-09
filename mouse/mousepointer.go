// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package mouse

import (
	"fbdev"
	"gfx"
	"time"
	"toolkit/base"
)


type MousePointer struct {
	base.Element
	fb			*fbdev.Framebuffer
}

const (
	WIDTH = 8
)


func (ms *Mouse) CreatePointer(fb *fbdev.Framebuffer, imp chan int64) (*MousePointer) {

	mp := &MousePointer{
		fb: fb,
	}

	mp.Element = base.Element{
		X:				int(fb.Vinfo.Xres/2),
		Y: 				int(fb.Vinfo.Yres/2),
		Width: 			WIDTH,
		Height:			WIDTH,
		Buffer: 		make([]byte, WIDTH*WIDTH*4),
		InvMsgPipe:		imp,
	}
	
	ms.RegisterMousePointer(mp.mouse)

	gfx.SetPixel(mp.Element.Buffer, 0, 0, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 1, 0, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 2, 0, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 3, 0, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 0, 1, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 1, 1, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 2, 1, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 0, 2, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 1, 2, mp.Element.Width, 0, 255, 0, 0)
   	gfx.SetPixel(mp.Element.Buffer, 0, 3, mp.Element.Width, 0, 255, 0, 0)
	return mp
}


func (mp *MousePointer) mouse(x int, y int, deltaX int, deltaY int, flags byte) {
	mp.Element.X += deltaX
	mp.Element.Y += deltaY
	mp.Element.InvMsgPipe <- time.Now().UnixNano()
}
