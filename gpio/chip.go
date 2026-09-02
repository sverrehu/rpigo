package gpio

import (
	"os"
	"time"

	"golang.org/x/sys/unix"
)

func ChipOpen(chipDeviceFile string) (*Chip, error) {
	f, err := os.OpenFile(chipDeviceFile, os.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	ci, err := getChipInfo(f)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	chip := &Chip{f: f, Path: chipDeviceFile}
	fillChip(chip, ci)
	return chip, nil
}

func NewChip(paths ...string) (*Chip, error) {
	if len(paths) > 0 && paths[0] != "" {
		return ChipOpen(paths[0])
	}
	return ChipOpen("/dev/gpiochip0")
}

func (c *Chip) Close() error {
	if c.f == nil {
		return nil
	}
	err := c.f.Close()
	c.f = nil
	return err
}

func (c *Chip) WatchLineInfo(offs int) (*LineInfo, error) {
	if c.f == nil {
		return nil, unix.EINVAL
	}
	li, err := getLineInfoWatch(c.f, offs)
	if err != nil {
		return nil, err
	}
	lineInfo := &LineInfo{}
	fillLineInfo(lineInfo, li)
	return lineInfo, nil
}

func (c *Chip) UnwatchLineInfo(offs int) error {
	if c.f == nil {
		return unix.EINVAL
	}
	return unwatchLineInfo(c.f, offs)
}

func (c *Chip) WaitInfoEvent(timeout time.Duration) (bool, error) {
	if c.f == nil {
		return false, unix.EINVAL
	}
	return pollFd(int(c.f.Fd()), timeout)
}

func (c *Chip) ReadInfoEvent() (*InfoEvent, error) {
	if c.f == nil {
		return nil, unix.EINVAL
	}
	return readInfoEventFd(int(c.f.Fd()))
}

func (c *Chip) LineOffsetFromName(name string) (int, error) {
	if c.f == nil {
		return -1, unix.EINVAL
	}
	if name == "" {
		return -1, unix.EINVAL
	}
	for offset := 0; offset < int(c.NumLines); offset++ {
		li, err := getLineInfo(c.f, offset)
		if err != nil {
			return -1, err
		}
		if cBytesToString(li.Name[:]) == name {
			return offset, nil
		}
	}
	return -1, unix.ENOENT
}

func (c *Chip) FindLine(name string) (*Line, error) {
	offset, err := c.LineOffsetFromName(name)
	if err != nil {
		return nil, err
	}
	return c.GetLine(offset)
}

func (c *Chip) GetLine(offset int) (*Line, error) {
	info, err := c.LineInfo(offset)
	if err != nil {
		return nil, err
	}
	return &Line{
		Chip:   c,
		Offset: offset,
		Info:   info,
	}, nil
}

func (c *Chip) LineInfo(offs int) (*LineInfo, error) {
	if c.f == nil {
		return nil, unix.EINVAL
	}
	li, err := getLineInfo(c.f, offs)
	if err != nil {
		return nil, err
	}
	lineInfo := &LineInfo{}
	fillLineInfo(lineInfo, li)
	return lineInfo, nil
}

func (c *Chip) RequestLines(requestConfig *RequestConfig, lineConfig *LineConfig) (*LineRequest, error) {
	if c.f == nil {
		return nil, unix.EINVAL
	}
	if lineConfig == nil {
		return nil, unix.EINVAL
	}
	uapiReq, err := toGpioV2LineRequest(requestConfig, lineConfig)
	if err != nil {
		return nil, err
	}
	ci, err := getChipInfo(c.f)
	if err != nil {
		return nil, err
	}
	if err := getLine(c.f, uapiReq); err != nil {
		return nil, err
	}
	lr := &LineRequest{
		ChipName: cBytesToString(ci.Name[:]),
		Fd:       int(uapiReq.Fd),
		NumLines: int(uapiReq.NumLines),
	}
	copy(lr.Offsets[:], uapiReq.Offsets[:])
	return lr, nil
}
