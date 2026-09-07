package pwm

import "time"

type PWM interface {
	Close() error
	SetDutyCycle(dutyCycle float64) error
	GetPeriod() time.Duration
}
