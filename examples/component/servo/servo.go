package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

// use "pinctl | grep PWM" to find chip and channel. rpi5 typically has other numbers than previous models.
const pwmChip = 0
const pwmChannel = 3 // rpi5: 3, rpi<5: 1

func main() {
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()

	pwm, err := pwm.NewHardPWM(pwmChip, pwmChannel, 50)
	if err != nil {
		log.Panic(err)
	}
	servo, err := component.NewServo(pwm)
	if err != nil {
		log.Panic(err)
	}
	defer servo.Close()
	servo.SetAngle(0)
	time.Sleep(1 * time.Second)
	servo.SetAngle(45)
	time.Sleep(1 * time.Second)
	servo.SetAngle(90)
	time.Sleep(1 * time.Second)
	servo.SetAngle(135)
	time.Sleep(1 * time.Second)
	servo.SetAngle(180)
	time.Sleep(1 * time.Second)
}
