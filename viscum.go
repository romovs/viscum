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
	cmp := compositor.CreateCompositor(fb)
	go cmp.Compose()
	
	// TODO: initialize touchscreen handler

	// initialize mouse handler & create the pointer
	ms, err := mouse.Init("/dev/input/mice", int(fb.Vinfo.Xres), int(fb.Vinfo.Yres), cmp.MouseWait, cmp.CompositorRelease, cmp.WindowList)
	if err != nil {
		log.Critical(err)
		return
	}
	defer ms.Close()
	go ms.Process()
	
	mp := ms.CreatePointer(fb, cmp.InvMsgPipe)
	cmp.RegisterMousePointer(&mp.Element)

	// create desktop
	desk, err := toolkit.CreateDesktop(fb, ms, cmp.InvMsgPipe, 50, 107, 89, 0)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&desk.Element)
	
	// create test app #1
	var win *toolkit.Window;
	win, err = toolkit.CreateWindow(cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 132, 345, 200, 200)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&win.Element)
	win.Button(ms, func () { log.Debug("Button clicked!") }, 70, 80, 60, 20)
	
	// create test app #2
	win, err = toolkit.CreateWindow(cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 500, 500, 100, 100)
	if err != nil {
		log.Critical(err)
		return
	}
	cmp.RegisterElement(&win.Element)
	
	
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