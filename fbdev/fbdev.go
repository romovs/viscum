// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// Linux Frame Buffer device initialization & management.
//
// RGB888 and ARGB888 support only.
//=====================================================================================================================

package fbdev

import (
	"unsafe"
	"os"
	"syscall"
	"utils"
	log "github.com/cihub/seelog"
	"image"
	"image/color"
	"gfx"
)



type Framebuffer struct {
	dev				*os.File
	tty				*os.File
	finfo			fixScreenInfo
	Vinfo			varScreenInfo
	Mem				[]byte
	MemOffscreen	[]byte
}

const (
	I_A = 3	
	I_R = 2
	I_G = 1
	I_B = 0
)

// Initialize frame buffer device
func Init(fbdev, tty string) (*Framebuffer, error) {

	var fb = new(Framebuffer)
	var err	error

	fb.tty, err = os.OpenFile(tty, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	// switch to graphics mode
	// this prevents kernel modifying the video ram (vt switching/blanking, cursor, gpm mouse cursor)
	err = ioctl(fb.tty.Fd(), KDSETMODE, unsafe.Pointer(uintptr(KD_GRAPHICS)))
	if err != nil {
        fb.tty.Close()
		return nil, err
	}

	fb.dev, err = os.OpenFile(fbdev, os.O_RDWR, os.ModeDevice)
	if err != nil {
		fb.tty.Close()
		return nil, err
	}
	
	err = ioctl(fb.dev.Fd(), FBIOGET_FSCREENINFO, unsafe.Pointer(&fb.finfo))
	if err != nil {
        fb.dev.Close()
        fb.tty.Close()
		return nil, err
	}
	log.Debug(utils.StructPrint(&fb.finfo))
	
	err = ioctl(fb.dev.Fd(), FBIOGET_VSCREENINFO, unsafe.Pointer(&fb.Vinfo))
	if err != nil {
        fb.dev.Close()
        fb.tty.Close()
		return nil, err
	}
	log.Debug(utils.StructPrint(&fb.Vinfo))
	
	memSize := int(fb.finfo.Smem_len + uint32(fb.finfo.Smem_start & uint64(syscall.Getpagesize() - 1)));

	fb.Mem, err = syscall.Mmap(int(fb.dev.Fd()), 0, memSize, PROT_READ | PROT_WRITE, MAP_SHARED)
	if err != nil {
        fb.dev.Close()
        fb.tty.Close()
		return nil, err
	}

	fb.MemOffscreen = make([]byte, memSize)

	return fb, nil
}


func (fb *Framebuffer) Close() {
	// restore text mode
	ioctl(fb.tty.Fd(), KDSETMODE, unsafe.Pointer(uintptr(KD_TEXT)))
	fb.tty.Close()
	syscall.Munmap(fb.Mem)
	fb.dev.Close()
}


func ioctl(fd uintptr, cmd uintptr, data unsafe.Pointer) error {
        _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, uintptr(data))
        if errno != 0 {
            return os.NewSyscallError("IOCTL", errno)
        }
        return nil
}

//---------------------------------------------------------------------------------------------------------------------
// draw.Image implementation
//
//---------------------------------------------------------------------------------------------------------------------

func (fb *Framebuffer) ColorModel() color.Model {
	return color.RGBAModel
}


func (fb *Framebuffer) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(fb.Vinfo.Xres), Y: int(fb.Vinfo.Yres)},
	}
}


func (fb *Framebuffer) At(x, y int) color.Color {
	offset := gfx.GetPixelOffset(x, y, int(fb.Vinfo.Xres))
	
	c := color.RGBA{
		R: fb.MemOffscreen[offset + I_R],
		G: fb.MemOffscreen[offset + I_G],
		B: fb.MemOffscreen[offset + I_B],
	}
	
	/*if fb.Vinfo.Transp.Length != 0 {
		c.A = fb.Mem[offset + I_A]
	}*/
	return c
}


func (fb *Framebuffer) Set(x, y int, c color.Color) {
	offset := gfx.GetPixelOffset(x, y, int(fb.Vinfo.Xres))
	r, g, b, _ := c.RGBA()
	
	fb.MemOffscreen[offset + I_R] = uint8(r)
	fb.MemOffscreen[offset + I_G] = uint8(g)
	fb.MemOffscreen[offset + I_B] = uint8(b)
	
	/*if fb.Vinfo.Transp.Length != 0 {
		fb.Mem[offset + I_A] = uint8(a)
	}*/
}