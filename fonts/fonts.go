// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package fonts

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"io/ioutil"
	log "github.com/cihub/seelog"
	"image"
	"image/draw"
	"strings"
	"image/color"
)


var (
	dpi float64 	= 72 	// screen resolution in Dots Per Inch
	size float64 	= 12	// font size in points
	spacing float64 = 1.5 	// line spacing (e.g. 2 means double spaced)
)

var defaultFont *truetype.Font = nil


func Default() (*truetype.Font) {
	
	if defaultFont == nil {
        fontBytes, err := ioutil.ReadFile("test.ttf")
        if err != nil {
                log.Critical(err)
                return nil
        }
        defaultFont, err = freetype.ParseFont(fontBytes)
        if err != nil {
                log.Critical(err)
                return nil
        }
	}

	return defaultFont
}


func ExpectedSize(font *truetype.Font, s string) (int32, int32, error) {

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)

	scale := size / float64(font.FUnitsPerEm())
        
	prev := font.Index(rune(s[0]))
	width := int32(font.HMetric(font.FUnitsPerEm(), prev).AdvanceWidth)
	for _, char := range s[1:] {
		index := font.Index(char)
		width += int32(font.Kerning(font.FUnitsPerEm(), prev, index) + 
					   font.HMetric(font.FUnitsPerEm(), index).AdvanceWidth)
		prev = index
	}
	width = int32(float64(width) * scale)

	bounds := font.Bounds(font.FUnitsPerEm())
	height := int32(float64(bounds.YMax-bounds.YMin) * scale)

	return width, height, nil
}


func Render(dst draw.Image, txt string, x, y, width, height int, font *truetype.Font) (error) {

	fg := &image.Uniform{color.RGBA{0, 0, 0, 255}}

	mask := image.NewAlpha(image.Rect(0, 0, width, height))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(mask.Bounds())
	c.SetDst(mask)
	c.SetSrc(image.Opaque)
		
	pt := freetype.Pt(0, int(c.PointToFix32(size) >> 8))
		
	strs := strings.Split(txt, "\n")

	for _, s := range strs {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log.Critical(err)
			return err
		}
		pt.Y += c.PointToFix32(size * spacing)
	}

	rect := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x+width, Y: y+height},
	}

    draw.DrawMask(dst, rect, fg, image.ZP, mask, image.ZP, draw.Over)

	return nil
}