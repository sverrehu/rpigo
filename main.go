package main

import (
	"time"

	"github.com/sverrehu/rpigo/pwm"
)

func main() {
	pwm := pwm.NewSoftPWM(nil, 50)
	pwm.SetDutyCycle(10)
	time.Sleep(5 * time.Hour)
}
