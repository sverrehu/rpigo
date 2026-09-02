package gpio

import "golang.org/x/sys/unix"

func (l *Line) Release() error {
	if l.Req == nil {
		return nil
	}
	err := l.Req.Close()
	l.Req = nil
	return err
}

func (l *Line) Setup(ls *LineSettings, consumer string) error {
	reqConfig := &RequestConfig{
		Consumer: consumer,
	}
	lineConfig := &LineConfig{
		Offsets:  []int{l.Offset},
		Settings: *ls,
	}
	req, err := l.Chip.RequestLines(reqConfig, lineConfig)
	if err != nil {
		return err
	}
	l.Req = req
	return nil
}

func (l *Line) SetValue(value LineValue) error {
	if l.Req == nil {
		return unix.EINVAL
	}
	return l.Req.SetValue(l.Offset, value)
}

func (l *Line) GetValue() (LineValue, error) {
	if l.Req == nil {
		return LineValueError, unix.EINVAL
	}
	return l.Req.GetValue(l.Offset)
}
