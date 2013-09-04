// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// Graphics functions.
//=====================================================================================================================

package gfx

import (

)


func GetPixelOffset(x, y, width int) (int) {

	return (y * width + x) * 4
}


func GetPixel(data []byte, x, y, width int) (red, green, blue, alpha byte) {		
			
	offset := GetPixelOffset(x, y, width)
	return data[offset+2], data[offset+1], data[offset], data[offset+3] 
}


func SetPixel(data []byte, x, y, width int, red, green, blue, alpha byte) {	
				
	offset := GetPixelOffset(x, y, width)
	data[offset] = blue
	data[offset+1] = green
	data[offset+2] = red
	data[offset+3] = alpha
}


func RectFilled(data []byte, x1, y1, x2, y2, width int, red, green, blue, alpha byte) {

	var xmin, xmax, ymin, ymax int;
	
	if x1 < x2 {
		xmin = x1
		xmax = x2	
	} else {
		xmin = x2
		xmax = x1
	}

	if y1 < y2 {
		ymin = y1
		ymax = y2	
	} else {
		ymin = y2
		ymax = y1
	}
	
	for i := xmin; i < xmax; i++ {
		for j := ymin; j < ymax; j++ {
			SetPixel(data, i, j, width, red, green, blue, alpha)
		}
	}
}


func Rect(data []byte, x1, y1, x2, y2, width int, red, green, blue, alpha byte) {

	// top/bottom
	for x := x1; x <= x2; x++ {
		SetPixel(data, x, y1, width, red, green, blue, alpha)
		SetPixel(data, x, y2, width, red, green, blue, alpha)
	}
	// left/right
	for y := y1; y <= y2; y++ {
		SetPixel(data, x1, y, width, red, green, blue, alpha)
		SetPixel(data, x2, y, width, red, green, blue, alpha)
	}
}


func Clear(data []byte, width, height int, red, green, blue, alpha byte) {

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			SetPixel(data, i, j, width, red, green, blue, alpha)
		}
	}
}