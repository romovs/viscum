package mouse

/*
#include <linux/input.h>
*/
import "C"

/*
 * linux/input.h
 */

const (
	EV_SYN uint16		= C.EV_SYN 
	EV_KEY uint16		= C.EV_KEY 
	EV_REL uint16		= C.EV_REL
	EV_ABS uint16		= C.EV_ABS
	EV_MSC uint16		= C.EV_MSC
	EV_SW uint16		= C.EV_SW
	EV_LED uint16		= C.EV_LED
	EV_SND uint16		= C.EV_SND
	EV_REP uint16		= C.EV_REP
	EV_FF uint16		= C.EV_FF
	EV_PWR uint16		= C.EV_PWR
	EV_FF_STATUS uint16	= C.EV_FF_STATUS
	EV_MAX uint16		= C.EV_MAX
	EV_CNT uint16		= C.EV_CNT
)

const(
	REL_X uint16		= C.REL_X
	REL_Y uint16		= C.REL_Y
	REL_Z uint16		= C.REL_Z
	REL_RX uint16		= C.REL_RX
	REL_RY uint16		= C.REL_RY
	REL_RZ uint16		= C.REL_RZ
	REL_HWHEEL uint16	= C.REL_HWHEEL
	REL_DIAL uint16		= C.REL_DIAL
	REL_WHEEL uint16	= C.REL_WHEEL
	REL_MISC uint16		= C.REL_MISC
	REL_MAX uint16		= C.REL_MAX
	REL_CNT uint16		= C.REL_CNT
)

type inputEvent C.struct_input_event
type timeVal C.struct_timeval
