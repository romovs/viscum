// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (
	"fbdev"
	"mouse"
	"gfx"
)

type Desktop struct {
	Element
	fb			*fbdev.Framebuffer
	
	red			byte
	green		byte
	blue		byte
	alpha		byte
}


func (desk *Desktop) Init(fb *fbdev.Framebuffer, ms *mouse.Mouse, imp chan int64, red, green, blue, alpha byte) (error) {

	desk.fb = fb
	desk.red = red
	desk.green = green
	desk.blue =	blue
	desk.alpha = alpha
	
	desk.Element = Element{
		Width: 		int(fb.Vinfo.Xres),
		Height: 	int(fb.Vinfo.Yres),
		Buffer: 	make([]byte, fb.Vinfo.Xres*fb.Vinfo.Yres*4),
		InvMsgPipe: imp,
		X:			0,
		Y: 			0,
		Z:			0,
	}
	
	gfx.Clear(desk.Element.Buffer, desk.Element.Width, desk.Element.Height, desk.red, desk.green, desk.blue, desk.alpha)

	return nil
}