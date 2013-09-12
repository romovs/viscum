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
	"os"
	"gfx"
)

func main() {
	defer log.Flush()
	utils.LoadLogConfig("test")
		
	// initialize frame buffer device
	fb, err := 	fbdev.Init("/dev/fb0", "/dev/tty0")
	if err != nil {
		log.Critical(err)
		os.Exit(1)
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
		os.Exit(1)
	}
	defer ms.Close()
	go ms.Process()
	
	mp := ms.CreatePointer(fb, cmp.InvMsgPipe)
	cmp.RegisterMousePointer(&mp.Element)

	// create desktop
	desk, err := toolkit.CreateDesktop(fb, ms, cmp.InvMsgPipe, 50, 107, 89, gfx.A_OPAQUE)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&desk.Element)
	
	// create test app #1
	var win *toolkit.Window;
	win, err = toolkit.CreateWindow(cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 132, 345, 200, 200)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&win.Element)
	win.Button(ms, "Click Me", func () { log.Debug("Button clicked!") }, 90, 80, 60, 20)
	win.Label(ms, "A button:", 30, 82, 60, 20)
	
	// create test app #2
	win, err = toolkit.CreateWindow(cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 500, 500, 100, 100)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&win.Element)

	fmt.Scanln()
}