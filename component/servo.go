package component

import (
	"fmt"
	"log"
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
	pulseRange := s.PulseControlMax - s.PulseControlMin
	pulseWidth := s.PulseControlMin + time.Duration(int64(angle*float64(pulseRange)/float64(s.MaxAngle)))
	dutyCycle := float64(100*pulseWidth) / float64(s.PWM.Period)
	s.PWM.SetDutyCycle(dutyCycle)
	log.Printf("angle: %f -> dutyCycle: %e, on: %s, off: %s", angle, dutyCycle, s.PWM.TimeOn, s.PWM.TimeOff)
	return nil
}

func (s *Servo) Close() error {
	s.PWM.Close()
	return nil
}
