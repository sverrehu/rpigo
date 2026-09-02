package component

import (
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

// TODO: very unfinished

type Servo struct {
	OutputComponent
	PulseControlMin time.Duration
	PulseControlMax time.Duration
	MaxAngle        int
	PWM             *gpio.SoftPWM
}

func NewServo(pin *Pin, frequency float64) (*Servo, error) {
	c, err := NewOutputComponent(pin)
	if err != nil {
		return nil, err
	}
	s := &Servo{
		OutputComponent: *c,
		PulseControlMin: time.Duration(500 * time.Microsecond),
		PulseControlMax: time.Duration(2500 * time.Microsecond),
		MaxAngle:        180,
		PWM:             gpio.NewSoftPWM(pin.line, frequency),
	}
	return s, nil
}

func (s *Servo) Set(b bool) error {
	return s.Pin.Set(b)
}

func (s *Servo) Close() error {
	s.PWM.Close()
	return nil
}
