package fbdev

/*
#include <linux/fb.h>
#include <sys/mman.h>
#include <linux/kd.h>
*/
import "C"

/*
 * linux/fb.h
 */
const (
	FBIOGET_VSCREENINFO uintptr = C.FBIOGET_VSCREENINFO
	FBIOPUT_VSCREENINFO uintptr = C.FBIOPUT_VSCREENINFO
	FBIOGET_FSCREENINFO uintptr = C.FBIOGET_FSCREENINFO
)

type fixScreenInfo C.struct_fb_fix_screeninfo
type varScreenInfo C.struct_fb_var_screeninfo
type bitField C.struct_fb_bitfield


/*
 * sys/mman.h
 */
const (
	PROT_READ int = C.PROT_READ	
	PROT_WRITE int = C.PROT_WRITE
	PROT_EXEC int = C.PROT_EXEC
)

const (
	MAP_SHARED int = C.MAP_SHARED
	MAP_PRIVATE int = C.MAP_PRIVATE
)

/*
 * linux/kd.h
 */
const (
	KDSETMODE uintptr = C.KDSETMODE
	KD_TEXT int = C.KD_TEXT
	KD_GRAPHICS int = C.KD_GRAPHICS
	KDGETMODE uintptr = C.KDGETMODE
)