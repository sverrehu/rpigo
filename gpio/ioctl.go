package gpio

import (
	"bytes"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const ioctlTypeGPIO = 0xb4

// From https://github.com/torvalds/linux/blob/master/include/uapi/asm-generic/ioctl.h
const iocNRBits = 8
const iocTypeBits = 8
const iocSizeBits = 14
const ioDirBits = 2

const iocNRMask = (1 << iocNRBits) - 1
const iocTypeMask = (1 << iocTypeBits) - 1
const iocSizeMask = (1 << iocSizeBits) - 1
const iocDirMask = (1 << ioDirBits) - 1

const iocNRShift = 0
const iocTypeShift = iocNRShift + iocNRBits
const iocSizeShift = iocTypeShift + iocTypeBits
const iocDirShift = iocSizeShift + iocSizeBits

const iocNone = 0b00
const iocWrite = 0b01
const iocRead = 0b10
const iocWriteRead = 0b11

// From https://github.com/torvalds/linux/blob/master/include/uapi/linux/gpio.h
const gpioMaxNameSize = 32
const gpioV2LinesMax = 64
const gpioV2LineNumAttrsMax = 10

var (
	gpioGetChipInfoIoctl        = gpioIor(0x01, unsafe.Sizeof(gpioChipInfo{}))
	gpioGetLineInfoUnwatchIoctl = gpioIowr(0x0c, unsafe.Sizeof(uint32(0)))
	gpioV2GetLineInfoIoctl      = gpioIowr(0x05, unsafe.Sizeof(gpioV2LineInfo{}))
	gpioV2GetLineInfoWatchIoctl = gpioIowr(0x06, unsafe.Sizeof(gpioV2LineInfo{}))
	gpioV2GetLineIoctl          = gpioIowr(0x07, unsafe.Sizeof(gpioV2LineRequest{}))
	gpioV2LineSetConfigIoctl    = gpioIowr(0x0d, unsafe.Sizeof(gpioV2LineConfig{}))
	gpioV2LineGetValuesIoctl    = gpioIowr(0x0e, unsafe.Sizeof(gpioV2LineValues{}))
	gpioV2LineSetValuesIoctl    = gpioIowr(0x0f, unsafe.Sizeof(gpioV2LineValues{}))
)

type gpioChipInfo struct {
	Name     [gpioMaxNameSize]byte
	Label    [gpioMaxNameSize]byte
	NumLines uint32
}

type gpioV2LineFlag uint64

const (
	gpioV2LineFlagUsed gpioV2LineFlag = 1 << iota
	gpioV2LineFlagActiveLow
	gpioV2LineFlagInput
	gpioV2LineFlagOutput
	gpioV2LineFlagEdgeRising
	gpioV2LineFlagEdgeFalling
	gpioV2LineFlagOpenDrain
	gpioV2LineFlagOpenSource
	gpioV2LineFlagBiasPullUp
	gpioV2LineFlagBiasPullDown
	gpioV2LineFlagBiasDisabled
	gpioV2LineFlagEventClockRealtime
	gpioV2LineFlagEventClockHTE
)

type gpioV2LineValues struct {
	Bits uint64
	Mask uint64
}

type gpioV2LineAttrId uint32

const (
	gpioV2LineAttrIdFlags = iota + 1
	gpioV2LineAttrIdOutputValues
	gpioV2LineAttrIdDebounce
)

type gpioV2LineAttribute struct {
	Id      gpioV2LineAttrId
	Padding uint32
	Value   uint64
}

type gpioV2LineConfigAttribute struct {
	Attr gpioV2LineAttribute
	Mask uint64
}

type gpioV2LineConfig struct {
	Flags    gpioV2LineFlag
	NumAttrs uint32
	Padding  [5]uint32
	Attrs    [gpioV2LineNumAttrsMax]gpioV2LineConfigAttribute
}

type gpioV2LineRequest struct {
	Offsets         [gpioV2LinesMax]uint32
	Consumer        [gpioMaxNameSize]byte
	Config          gpioV2LineConfig
	NumLines        uint32
	EventBufferSize uint32
	Padding         [5]uint32
	Fd              int32
}

type gpioV2LineInfo struct {
	Name     [gpioMaxNameSize]byte
	Consumer [gpioMaxNameSize]byte
	Offset   uint32
	NumAttrs uint32
	Flags    gpioV2LineFlag
	Attrs    [gpioV2LineNumAttrsMax]gpioV2LineAttribute
	Padding  [4]uint32
}

type gpioV2LineChangedType int

const (
	gpioV2LineChangedRequested gpioV2LineChangedType = iota + 1
	gpioV2LineChangedReleased
	gpioV2LineChangedConfig
)

type gpioV2LineInfoChanged struct {
	Info        gpioV2LineInfo
	TimestampNs uint64
	EventType   uint32
	Padding     [5]uint32
}

type gpioV2LineEventId int

const (
	gpioV2LineEventRisingEdge gpioV2LineEventId = iota + 1
	gpioV2LineEventFallingEdge
)

type gpioV2LineEvent struct {
	TimestampNs uint64
	Id          uint32
	Offset      uint32
	SeqNo       uint32
	LineSeqNo   uint32
	Padding     [6]uint32
}

func getChipInfo(f *os.File) (*gpioChipInfo, error) {
	var ci gpioChipInfo
	err := ioctl(f, gpioGetChipInfoIoctl, uintptr(unsafe.Pointer(&ci)))
	if err != nil {
		return nil, err
	}
	return &ci, nil
}

func getLineInfo(f *os.File, offs int) (*gpioV2LineInfo, error) {
	ci := gpioV2LineInfo{Offset: uint32(offs)}
	err := ioctl(f, gpioV2GetLineInfoIoctl, uintptr(unsafe.Pointer(&ci)))
	if err != nil {
		return nil, err
	}
	return &ci, nil
}

func getLineInfoWatch(f *os.File, offs int) (*gpioV2LineInfo, error) {
	ci := gpioV2LineInfo{Offset: uint32(offs)}
	err := ioctl(f, gpioV2GetLineInfoWatchIoctl, uintptr(unsafe.Pointer(&ci)))
	if err != nil {
		return nil, err
	}
	return &ci, nil
}

func unwatchLineInfo(f *os.File, offs int) error {
	offset := uint32(offs)
	return ioctl(f, gpioGetLineInfoUnwatchIoctl, uintptr(unsafe.Pointer(&offset)))
}

func getLine(f *os.File, request *gpioV2LineRequest) error {
	return ioctl(f, gpioV2GetLineIoctl, uintptr(unsafe.Pointer(request)))
}

func lineSetConfig(fd int, cfg *gpioV2LineConfig) error {
	return ioctlFd(fd, gpioV2LineSetConfigIoctl, uintptr(unsafe.Pointer(cfg)))
}

func lineGetValues(fd int, vals *gpioV2LineValues) error {
	return ioctlFd(fd, gpioV2LineGetValuesIoctl, uintptr(unsafe.Pointer(vals)))
}

func lineSetValues(fd int, vals *gpioV2LineValues) error {
	return ioctlFd(fd, gpioV2LineSetValuesIoctl, uintptr(unsafe.Pointer(vals)))
}

func ioc(dir, tp, nr, size uintptr) uintptr {
	return ((dir & iocDirMask) << iocDirShift) | ((tp & iocTypeMask) << iocTypeShift) | ((nr & iocNRMask) << iocNRShift) | ((size & iocSizeMask) << iocSizeShift)
}

func io(tp, nr, size uintptr) uintptr {
	return ioc(iocNone, tp, nr, size)
}

func ior(tp, nr, size uintptr) uintptr {
	return ioc(iocRead, tp, nr, size)
}

func iow(tp, nr, size uintptr) uintptr {
	return ioc(iocWrite, tp, nr, size)
}

func iowr(tp, nr, size uintptr) uintptr {
	return ioc(iocWriteRead, tp, nr, size)
}

func gpioIor(nr, size uintptr) uintptr {
	return ior(ioctlTypeGPIO, nr, size)
}

func gpioIow(nr, size uintptr) uintptr {
	return iow(ioctlTypeGPIO, nr, size)
}

func gpioIowr(nr, size uintptr) uintptr {
	return iowr(ioctlTypeGPIO, nr, size)
}

func ioctl(f *os.File, op uintptr, ptr uintptr) error {
	return ioctlFd(int(f.Fd()), op, ptr)
}

func ioctlFd(fd int, op uintptr, ptr uintptr) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), op, ptr)
	if errno != 0 {
		return errno
	}
	return nil
}

func cBytesToString(a []byte) string {
	before, _, ok := bytes.Cut(a, []byte{0})
	if !ok {
		return string(a)
	}
	return string(before)
}
