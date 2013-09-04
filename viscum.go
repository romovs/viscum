// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package main

import (
	"fmt"
	"fbdev"
	"toolkit"
	"mouse"
	"utils"
	_ "image/jpeg"
	_ "image/png"
	log "github.com/cihub/seelog"
	"compositor"
)

func main() {
	defer log.Flush()
	utils.LoadLogConfig("test")
		
	// initialize frame buffer device
	fb, err := 	fbdev.Init("/dev/fb0", "/dev/tty0")
	if err != nil {
		log.Critical(err)
		return
	}
	defer fb.Close()

	// initialize compositor
	cmp := compositor.Init(fb)
	go cmp.Compose()

	// TODO: initialize touchscreen handler

	// initialize mouse handler
	ms, err := mouse.Init("/dev/input/mice", fb, cmp.MouseWait, cmp.CompositorRelease)
	if err != nil {
		log.Critical(err)
		return
	}
	defer ms.Close()
	go ms.Process()

	// create desktop
	desk := new(toolkit.Desktop)
	err = desk.Init(fb, ms, cmp.InvMsgPipe, 50, 107, 89, 0)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&desk.Element)
	
	// create test app #1
	win := new(toolkit.Window)
	err = win.Init(fb, ms, cmp.InvMsgPipe, 132, 345, 200, 200, 1)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&win.Element)
	
	// create test app #2
	win = new(toolkit.Window)
	err = win.Init(fb, ms, cmp.InvMsgPipe, 500, 500, 100, 100, 2)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&win.Element)
	
	// create mouse pointer
	mp := new(toolkit.MousePointer)
	err = mp.Init(fb, ms, cmp.InvMsgPipe)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&mp.Element)
	
	
	/*f, err := os.Open("test.png")
    if err != nil {
		log.Critical(err)
		return
    }
    defer f.Close()
    		
    img, _, err := image.Decode(bufio.NewReader(f))
    if err != nil  {
		log.Critical(err)
		return
	}*/
	
	//fb.Draw(350, 250, img)

	fmt.Scanln()
}