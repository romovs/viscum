// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// All toolkit elements should implement IElement
//=====================================================================================================================


package base

import (
	"gfx"
	"image/color"
	"image"
	"container/list"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const (
	I_A = 3	
	I_R = 2
	I_G = 1
	I_B = 0
)

type IElement interface {
	UpdateScreenX(int)
	UpdateScreenY(int)
	
	ColorModel() color.Model
	Bounds() image.Rectangle
	At(x, y int) color.Color
	Set(x, y int, c color.Color)
	
	Mouse(int, int, int, int, byte) 
	Draw()
	
	GetId() uint64
}

type Element struct {
	Id				uint64
	InvMsgPipe		chan int64
	X				int
	Y				int
	Width			int
	Height			int
	ScreenX			int 
	ScreenY			int
	Buffer			[]byte
	DeactivateHndr	func()
	Children		*list.List
	CompRemoveHdnr	func(uint64)
	Font			*truetype.Font
}


//---------------------------------------------------------------------------------------------------------------------
// draw.Image implementation
//
//---------------------------------------------------------------------------------------------------------------------

func (e *Element) ColorModel() color.Model {
	return color.RGBAModel
}


func (e *Element) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: e.Width, Y: e.Height},
	}
}


func (e *Element) At(x, y int) color.Color {
	offset := gfx.GetPixelOffset(x, y, e.Width)
	
	c := color.RGBA{
		R: e.Buffer[offset + I_R],
		G: e.Buffer[offset + I_G],
		B: e.Buffer[offset + I_B],
	}

	/*if fb.Vinfo.Transp.Length != 0 {
		c.A = fb.Data[offset + I_A]
	}*/
	return c
}


func (e *Element) Set(x, y int, c color.Color) {
	offset := gfx.GetPixelOffset(x, y, e.Width)
	r, g, b, _ := c.RGBA()
	
	e.Buffer[offset + I_R] = uint8(r)
	e.Buffer[offset + I_G] = uint8(g)
	e.Buffer[offset + I_B] = uint8(b)
	
	/*if fb.Vinfo.Transp.Length != 0 {
		fb.Data[offset + I_A] = uint8(a)
	}*/
}


func (e *Element) UpdateScreenX(deltaX int) {
	e.ScreenX += deltaX
}

func (e *Element) UpdateScreenY(deltaY int) {
	e.ScreenY += deltaY
}

func (e *Element) GetId() (uint64) {
	return e.Id
}

/*func (o *Element) Dr(fb *fbdev.Framebuffer, x, y, width, height int) {
	
	rect := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x+width, Y: y+height},
	}
	
	draw.Draw(fb, rect, o, o.Bounds().Min, draw.Src)
}*/