// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs _defines.go

package mouse

const (
	EV_SYN		uint16	= 0x0
	EV_KEY		uint16	= 0x1
	EV_REL		uint16	= 0x2
	EV_ABS		uint16	= 0x3
	EV_MSC		uint16	= 0x4
	EV_SW		uint16	= 0x5
	EV_LED		uint16	= 0x11
	EV_SND		uint16	= 0x12
	EV_REP		uint16	= 0x14
	EV_FF		uint16	= 0x15
	EV_PWR		uint16	= 0x16
	EV_FF_STATUS	uint16	= 0x17
	EV_MAX		uint16	= 0x1f
	EV_CNT		uint16	= 0x20
)

const (
	REL_X		uint16	= 0x0
	REL_Y		uint16	= 0x1
	REL_Z		uint16	= 0x2
	REL_RX		uint16	= 0x3
	REL_RY		uint16	= 0x4
	REL_RZ		uint16	= 0x5
	REL_HWHEEL	uint16	= 0x6
	REL_DIAL	uint16	= 0x7
	REL_WHEEL	uint16	= 0x8
	REL_MISC	uint16	= 0x9
	REL_MAX		uint16	= 0xf
	REL_CNT		uint16	= 0x10
)

type inputEvent struct {
	Time	timeVal
	Type	uint16
	Code	uint16
	Value	int32
}
type timeVal struct {
	Sec	int64
	Usec	int64
}
