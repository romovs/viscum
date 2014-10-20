// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package compositor

import (
	"viscum/fbdev"
	"container/list"
	"viscum/toolkit/base"
	"viscum/gfx"
)

type Compositor struct {
	fb 					*fbdev.Framebuffer
	WindowList			*list.List
	InvMsgPipe 			chan int64
	MouseWait			chan bool		// mouse-compositor synchronization
	CompositorRelease	chan bool		// mouse-compositor synchronization
	mousePointer		*base.Element
}


func CreateCompositor(fb *fbdev.Framebuffer) (*Compositor) {

	c := &Compositor{
		fb: 				fb,
		WindowList: 		list.New(),
		InvMsgPipe: 		make(chan int64, 128),
		MouseWait: 			make(chan bool, 1),
		CompositorRelease:	make(chan bool, 1),
	} 
	
	c.CompositorRelease<- false
	
	return c
}


func (c *Compositor) Compose() {

    for {
        <-c.InvMsgPipe

		<-c.MouseWait

		// render desktop, windows, and child elements
		for v := c.WindowList.Back(); v != nil; v = v.Prev() {
			e := v.Value.(*base.Element)
			gfx.DrawSrc(c.fb, e, e.X, e.Y, e.Width, e.Height)
		}
		
		// render mouse pointer
		if c.mousePointer != nil {
			gfx.DrawOver(c.fb, c.mousePointer, c.mousePointer.X, c.mousePointer.Y, c.mousePointer.Width, c.mousePointer.Height)
		}

		flush(c)
		
		c.CompositorRelease<- false
    }   
}



func flush(c *Compositor) {
	copy(c.fb.Mem, c.fb.MemOffscreen)
}

func (c *Compositor) RegisterElement(e *base.Element) {
	e.CompRemoveHdnr = c.RemoveElement
	c.WindowList.PushFront(e)
}

func (c *Compositor) RemoveElement(id uint64) {
	for v := c.WindowList.Front(); v != nil; v = v.Next() {
		e := v.Value.(*base.Element)
		if e.Id == id {
			c.WindowList.Remove(v)
			e.Children.Init()			
		}
	}
}

func (c *Compositor) RegisterMousePointer(mousePointer *base.Element) {
	c.mousePointer = mousePointer
}

func (c *Compositor) ActivateWindow(id uint64) {

	// deactive currently active window
	if c.WindowList.Front().Value.(*base.Element).DeactivateHndr != nil { // could be null for Desktop
		c.WindowList.Front().Value.(*base.Element).DeactivateHndr()
	}
	
	for v := c.WindowList.Front(); v != nil; v = v.Next() {
		// move it to the front of the list - this essentially makes it to be the currently active window
		e := v.Value.(*base.Element)
		if e.Id == id {
			c.WindowList.MoveToFront(v)
			break
		}
	}
}
