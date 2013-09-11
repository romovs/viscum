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

func Render(txt string, width, height int, font *truetype.Font) (*image.RGBA, error) {

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)
		
	pt := freetype.Pt(0, int(c.PointToFix32(size) >> 8))
		
	strs := strings.Split(txt, "\n")

	for _, s := range strs {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log.Critical(err)
			return nil, err
		}
		pt.Y += c.PointToFix32(size * spacing)
	}

	return rgba, nil
}