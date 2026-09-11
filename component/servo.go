package component

import (
	"fmt"
	"time"

	"github.com/sverrehu/rpigo/pwm"
)

type Servo struct {
	OutputComponent
	PWM             pwm.PWM
	PulseControlMin time.Duration
	PulseControlMax time.Duration
	MaxAngle        int
}

func NewServo(pwm pwm.PWM) (*Servo, error) {
	// Matches several servos, including SG92R (50 Hz)
	return NewServoWithSpec(pwm, 500*time.Microsecond, 2500*time.Microsecond, 180)
}

func NewServoWithSpec(pwm pwm.PWM, pulseControlMin, pulseControlMax time.Duration, maxAngle int) (*Servo, error) {
	c, err := NewOutputComponent(nil)
	if err != nil {
		return nil, err
	}
	s := &Servo{
		PWM:             pwm,
		OutputComponent: *c,
		PulseControlMin: pulseControlMin,
		PulseControlMax: pulseControlMax,
		MaxAngle:        maxAngle,
	}
	return s, nil
}

func (s *Servo) SetAngle(angle float64) error {
	if angle < 0 || angle > float64(s.MaxAngle) {
		return fmt.Errorf("angle must be between 0 and %d inclusive", s.MaxAngle)
	}
	pulseRange := s.PulseControlMax - s.PulseControlMin
	pulseWidth := s.PulseControlMin + time.Duration(int64(angle*float64(pulseRange)/float64(s.MaxAngle)))
	dutyCycle := float64(100*pulseWidth) / float64(s.PWM.GetPeriod())
	s.PWM.SetDutyCycle(dutyCycle)
	return nil
}

func (s *Servo) Close() error {
	s.PWM.Close()
	return nil
}
