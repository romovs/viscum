// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package mouse

import (
	"fbdev"
	"gfx"
	"time"
	"toolkit/base"
	"os"
	"image"
	"bufio"
	log "github.com/cihub/seelog"
)

type MousePointer struct {
	base.Element
	fb			*fbdev.Framebuffer
}


func (ms *Mouse) CreatePointer(fb *fbdev.Framebuffer, imp chan int64) (*MousePointer) {

	f, err := os.Open("data/cursor/arrow.png")
    if err != nil {
		log.Critical(err)
		return nil
    }
    defer f.Close()
    		
    img, _, err := image.Decode(bufio.NewReader(f))
    if err != nil  {
		log.Critical(err)
		return nil
	}
	
	mp := &MousePointer{
		fb: fb,
	}
	
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	mp.Element = base.Element{
		X:				int(fb.Vinfo.Xres/2),
		Y: 				int(fb.Vinfo.Yres/2),
		Width: 			w,
		Height:			h,
		Buffer: 		make([]byte, w*h*4),
		InvMsgPipe:		imp,
	}
	
	ms.RegisterMousePointer(mp.mouse)

	gfx.DrawSrc(&mp.Element, img, 0, 0, w, h)

	return mp
}


func (mp *MousePointer) mouse(x int, y int, deltaX int, deltaY int, flags uint16) {
	mp.Element.X += deltaX
	mp.Element.Y += deltaY
	mp.Element.InvMsgPipe <- time.Now().UnixNano()
}
