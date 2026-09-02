package component

import "github.com/sverrehu/rpigo/gpio"

type Pin struct {
	gpio *GPIO
	line *gpio.Line
}

func newPin(gpio *GPIO, line *gpio.Line) *Pin {
	return &Pin{gpio: gpio, line: line}
}

func (p *Pin) Close() error {
	err := p.line.Release()
	p.line = nil
	return err
}

func (p *Pin) IsClosed() bool {
	return p.line == nil
}

func (p *Pin) Set(b bool) error {
	if b {
		return p.line.SetValue(gpio.LineValueActive)
	}
	return p.line.SetValue(gpio.LineValueInactive)
}

func (p *Pin) Get() (bool, error) {
	v, err := p.line.GetValue()
	if err != nil {
		return false, err
	}
	return v == gpio.LineValueActive, nil
}

func (p *Pin) setup(ls *gpio.LineSettings) error {
	err := p.line.Setup(ls, p.gpio.Consumer)
	if err != nil {
		return err
	}
	return nil
}
