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
	
	Mouse(int, int, int, int, uint16) 
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
	
	return color.RGBA{
		R: e.Buffer[offset + I_R],
		G: e.Buffer[offset + I_G],
		B: e.Buffer[offset + I_B],
		A: e.Buffer[offset + I_A],
	}
}


func (e *Element) Set(x, y int, c color.Color) {
	offset := gfx.GetPixelOffset(x, y, e.Width)
	r, g, b, a := c.RGBA()
	
	// since image.Draw works with RGBA64 internally we need to convert the color value to RGBA before writing 
	// to the framebuffer - hence the div 256
	e.Buffer[offset + I_R] = uint8(r/256)
	e.Buffer[offset + I_G] = uint8(g/256)
	e.Buffer[offset + I_B] = uint8(b/256)
	e.Buffer[offset + I_A] = uint8(a/256)
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
