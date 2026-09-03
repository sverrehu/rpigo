package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
)

func main() {
	const pwmPinName = component.GPIO12
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	pwmPin, err := gpio.Pin(pwmPinName)
	if err != nil {
		log.Panic(err)
	}
	servo, err := component.NewServo(pwmPin)
	if err != nil {
		log.Panic(err)
	}
	servo.SetAngle(90)
	time.Sleep(2 * time.Second)
	servo.SetAngle(0)
	time.Sleep(2 * time.Second)
	servo.SetAngle(180)
	time.Sleep(2 * time.Second)
	servo.SetAngle(90)
	time.Sleep(2 * time.Second)
}
