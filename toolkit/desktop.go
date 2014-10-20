// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"viscum/fbdev"
	"viscum/mouse"
	"viscum/gfx"
	"viscum/toolkit/base"
)

type Desktop struct {
	base.Element
	fb			*fbdev.Framebuffer
	red			byte
	green		byte
	blue		byte
	alpha		byte
}


func CreateDesktop(fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, red, green, blue, alpha byte) (*Desktop, error) {

	desk := &Desktop{
		fb: 	fb,
		red: 	red,
		green: 	green,
		blue:	blue,
		alpha: 	alpha,
	}
	
	desk.Element = base.Element{
		Width: 		int(fb.Vinfo.Xres),
		Height: 	int(fb.Vinfo.Yres),
		Buffer: 	make([]byte, fb.Vinfo.Xres*fb.Vinfo.Yres*4),
		InvMsgPipe: imp,
		X:			0,
		Y: 			0,
	}
	
	gfx.Clear(desk.Element.Buffer, desk.Element.Width, desk.Element.Height, desk.red, desk.green, desk.blue, desk.alpha)

	return desk, nil
}
