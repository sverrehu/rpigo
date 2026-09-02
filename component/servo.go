package component

import (
	"fmt"
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

func NewServo(pin *Pin) (*Servo, error) {
	// Matches several servos, including SG92R
	return NewServoWithSpec(pin, 500*time.Microsecond, 2500*time.Microsecond, 180, 50)
}

func NewServoWithSpec(pin *Pin, pulseControlMin, pulseControlMax time.Duration, maxAngle int, frequency float64) (*Servo, error) {
	c, err := NewOutputComponent(pin)
	if err != nil {
		return nil, err
	}
	s := &Servo{
		OutputComponent: *c,
		PulseControlMin: pulseControlMin,
		PulseControlMax: pulseControlMax,
		MaxAngle:        maxAngle,
		PWM:             gpio.NewSoftPWM(pin.line, frequency),
	}
	return s, nil
}

func (s *Servo) SetAngle(angle float64) error {
	if angle < 0 || angle > float64(s.MaxAngle) {
		return fmt.Errorf("angle must be between 0 and %d inclusive", s.MaxAngle)
	}
	// TODO: calculate duty cycle
	dutyCycle := 0
	s.PWM.SetDutyCycle(dutyCycle)
	return nil
}

func (s *Servo) Close() error {
	s.PWM.Close()
	return nil
}
