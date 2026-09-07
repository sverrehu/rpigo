package gpio

import (
	"time"

	"golang.org/x/sys/unix"
)

func NewLineRequest() *LineRequest {
	return &LineRequest{}
}

func (r *LineRequest) RequestedOffsets() []int {
	offs := make([]int, r.NumLines)
	for i := 0; i < r.NumLines; i++ {
		offs[i] = int(r.Offsets[i])
	}
	return offs
}

func (r *LineRequest) offsetToBit(offset int) int {
	for i := 0; i < r.NumLines; i++ {
		if r.Offsets[i] == uint32(offset) {
			return i
		}
	}
	return -1
}

func (r *LineRequest) GetValue(offset int) (LineValue, error) {
	if r.Fd < 0 {
		return LineValueError, unix.EINVAL
	}
	bit := r.offsetToBit(offset)
	if bit < 0 {
		return LineValueError, unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	lineMaskSetBit(&uapiVals.Mask, bit)
	if err := lineGetValues(r.Fd, &uapiVals); err != nil {
		return LineValueError, err
	}
	if lineMaskTestBit(&uapiVals.Bits, bit) {
		return LineValueActive, nil
	}
	return LineValueInactive, nil
}

func (r *LineRequest) GetValuesSubset(offsets []int) (map[int]LineValue, error) {
	if r.Fd < 0 {
		return nil, unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	for _, off := range offsets {
		bit := r.offsetToBit(off)
		if bit < 0 {
			return nil, unix.EINVAL
		}
		lineMaskSetBit(&uapiVals.Mask, bit)
	}
	if err := lineGetValues(r.Fd, &uapiVals); err != nil {
		return nil, err
	}
	res := make(map[int]LineValue, len(offsets))
	for _, off := range offsets {
		bit := r.offsetToBit(off)
		if lineMaskTestBit(&uapiVals.Bits, bit) {
			res[off] = LineValueActive
		} else {
			res[off] = LineValueInactive
		}
	}
	return res, nil
}

func (r *LineRequest) GetValues() ([]LineValue, error) {
	if r.Fd < 0 {
		return nil, unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	for i := 0; i < r.NumLines; i++ {
		lineMaskSetBit(&uapiVals.Mask, i)
	}
	if err := lineGetValues(r.Fd, &uapiVals); err != nil {
		return nil, err
	}
	vals := make([]LineValue, r.NumLines)
	for i := 0; i < r.NumLines; i++ {
		if lineMaskTestBit(&uapiVals.Bits, i) {
			vals[i] = LineValueActive
		} else {
			vals[i] = LineValueInactive
		}
	}
	return vals, nil
}

func (r *LineRequest) SetValue(offset int, value LineValue) error {
	if r.Fd < 0 {
		return unix.EINVAL
	}
	bit := r.offsetToBit(offset)
	if bit < 0 {
		return unix.EINVAL
	}
	if value != LineValueInactive && value != LineValueActive {
		return unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	lineMaskSetBit(&uapiVals.Mask, bit)
	lineMaskAssignBit(&uapiVals.Bits, bit, value == LineValueActive)
	return lineSetValues(r.Fd, &uapiVals)
}

func (r *LineRequest) SetValuesSubset(values map[int]LineValue) error {
	if r.Fd < 0 {
		return unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	for off, val := range values {
		bit := r.offsetToBit(off)
		if bit < 0 {
			return unix.EINVAL
		}
		if val != LineValueInactive && val != LineValueActive {
			return unix.EINVAL
		}
		lineMaskSetBit(&uapiVals.Mask, bit)
		lineMaskAssignBit(&uapiVals.Bits, bit, val == LineValueActive)
	}
	return lineSetValues(r.Fd, &uapiVals)
}

func (r *LineRequest) SetValues(values []LineValue) error {
	if r.Fd < 0 || len(values) < r.NumLines {
		return unix.EINVAL
	}
	var uapiVals gpioV2LineValues
	for i := 0; i < r.NumLines; i++ {
		if values[i] != LineValueInactive && values[i] != LineValueActive {
			return unix.EINVAL
		}
		lineMaskSetBit(&uapiVals.Mask, i)
		lineMaskAssignBit(&uapiVals.Bits, i, values[i] == LineValueActive)
	}
	return lineSetValues(r.Fd, &uapiVals)
}

func (r *LineRequest) ReconfigureLines(config *LineConfig) error {
	if r.Fd < 0 || config == nil {
		return unix.EINVAL
	}
	uapiReq, err := toGpioV2LineRequest(nil, config)
	if err != nil {
		return err
	}
	if int(uapiReq.NumLines) != r.NumLines {
		return unix.EINVAL
	}
	for i := 0; i < r.NumLines; i++ {
		if uapiReq.Offsets[i] != r.Offsets[i] {
			return unix.EINVAL
		}
	}
	return lineSetConfig(r.Fd, &uapiReq.Config)
}

func (r *LineRequest) WaitEdgeEvents(timeout time.Duration) (bool, error) {
	if r.Fd < 0 {
		return false, unix.EINVAL
	}
	return pollFd(r.Fd, timeout)
}

func (r *LineRequest) ReadEdgeEvents(buffer *EdgeEventBuffer, maxEvents int) (int, error) {
	if r.Fd < 0 || buffer == nil {
		return 0, unix.EINVAL
	}
	return buffer.ReadFd(r.Fd, maxEvents)
}

func (r *LineRequest) Close() error {
	if r.Fd < 0 {
		return nil
	}
	err := unix.Close(r.Fd)
	r.Fd = -1
	return err
}
