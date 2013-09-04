package compositor

import (
	"fbdev"
	"image"
	"image/draw"
	log "github.com/cihub/seelog"
	"time"
	"toolkit"
)

type Compositor struct {
	fb 					*fbdev.Framebuffer
	Elms				[]*toolkit.Element
	ElmsCnt				int
	InvMsgPipe 			chan int64
	
	MouseWait			chan bool
	CompositorRelease	chan bool
	
	lastInv		time.Time
}


func Init(fb *fbdev.Framebuffer) (*Compositor) {

	c := &Compositor{
		fb: 				fb,
		Elms: 				make([]*toolkit.Element, 128),
		InvMsgPipe: 		make(chan int64, 128),
		MouseWait: 			make(chan bool, 1),
		CompositorRelease:	make(chan bool, 1),
		ElmsCnt: 			0,
		lastInv:			time.Now(),
	} 

	c.CompositorRelease<- false
	
	return c
}


func (c *Compositor) Compose() {

    for {
        <-c.InvMsgPipe
				
		<-c.MouseWait

		for i := 0; i < c.ElmsCnt; i++ {
			RenderElement(c.fb, c.Elms[i], c.Elms[i].X, c.Elms[i].Y, c.Elms[i].Width, c.Elms[i].Height)
		}

		flush(c)
		
		c.CompositorRelease<- false
    }   
}

func flush(c *Compositor) {
	copy(c.fb.Mem, c.fb.MemOffscreen)
}

func (c *Compositor) RegisterElement(e *toolkit.Element) {
	c.Elms[c.ElmsCnt] = e
	c.ElmsCnt++
}

func RenderElement(fb *fbdev.Framebuffer, o *toolkit.Element, x, y, width, height int) {
	
	rect := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x+width, Y: y+height},
	}
	
	log.Debugf("Rendering: %v  at: %v", o.Bounds(), rect)

	draw.Draw(fb, rect, o, o.Bounds().Min, draw.Src)
}