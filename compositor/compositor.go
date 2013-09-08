package compositor

import (
	"fbdev"
	"image"
	"image/draw"
	log "github.com/cihub/seelog"
	"toolkit"
	"container/list"
)

type Compositor struct {
	fb 					*fbdev.Framebuffer
	Elms				*list.List
	InvMsgPipe 			chan int64
	
	MouseWait			chan bool
	CompositorRelease	chan bool
}


func Init(fb *fbdev.Framebuffer) (*Compositor) {

	c := &Compositor{
		fb: 				fb,
		Elms: 				list.New(),
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

		for v := c.Elms.Back(); v != nil; v = v.Prev() {
			e := v.Value.(*toolkit.Element)
			RenderElement(c.fb, e, e.X, e.Y, e.Width, e.Height)
		}

		flush(c)
		
		c.CompositorRelease<- false
    }   
}

func flush(c *Compositor) {
	copy(c.fb.Mem, c.fb.MemOffscreen)
}

func (c *Compositor) RegisterElement(e *toolkit.Element) {
	c.Elms.PushFront(e)
}

func RenderElement(fb *fbdev.Framebuffer, o *toolkit.Element, x, y, width, height int) {
	
	rect := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x+width, Y: y+height},
	}
	
	log.Debugf("Rendering: %v  at: %v", o.Bounds(), rect)

	draw.Draw(fb, rect, o, o.Bounds().Min, draw.Src)
}