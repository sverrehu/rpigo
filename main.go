package main

import (
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

func main() {
	pwm := gpio.NewSoftPWM(nil, 50)
	pwm.SetDutyCycle(10)
	time.Sleep(5 * time.Hour)
}
