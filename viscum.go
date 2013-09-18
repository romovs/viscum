// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package main

import (
	"fmt"
	"fbdev"
	"toolkit"
	"mouse"
	"utils"
	"image"
	_ "image/png"
	log "github.com/cihub/seelog"
	"compositor"
	"os"
	"gfx"
	"bufio"
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

	go cmp.Compose()

	// create desktop
	desk, err := toolkit.CreateDesktop(fb, ms, cmp.InvMsgPipe, 50, 107, 89, gfx.A_OPAQUE)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&desk.Element)
	
	// create test app #1
	var win *toolkit.Window;
	win, err = toolkit.CreateWindow(0, cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 132, 345, 250, 200)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&win.Element)
	win.Button(toolkit.BS_TEXT, ms, nil, "Click Me", func (_ bool) { log.Debug("Button clicked!") }, 10, 160, 60, 30)
	
	f, err := os.Open("data/v.png")
    if err != nil {
		log.Critical(err)
		os.Exit(1)
    }
    defer f.Close()
    		
    img, _, err := image.Decode(bufio.NewReader(f))
    if err != nil  {
		log.Critical(err)
		os.Exit(1)
	}
	win.Button(toolkit.BS_ICON_TEXT, ms, img, "No! Me!!!", func (_ bool) { log.Debug("Button clicked!") }, 80, 160, 100, 30)
	win.Button(toolkit.BS_ICON | toolkit.BS_TOGGLE, ms, img, "", func (pushed bool) { log.Debugf("Button toggled! pushed: %t", pushed) }, 190, 160, 30, 30)
	
	win.Label(ms, "Testing...", 100, 50, 60, 20)
	
	// create test app #2
	win, err = toolkit.CreateWindow(0, cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 500, 500, 100, 100)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&win.Element)
	
	// create test app #3
	win, err = toolkit.CreateWindow(toolkit.WS_TITLEBAR_HIDDEN, cmp.ActivateWindow, fb, ms, cmp.InvMsgPipe, 300, 150, 150, 150)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	cmp.RegisterElement(&win.Element)
	

	fmt.Scanln()
}