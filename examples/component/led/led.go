package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
)

func main() {
	const pinName = component.GPIO18
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	pin, err := gpio.Pin(pinName)
	if err != nil {
		log.Panic(err)
	}
	led, err := component.NewLED(pin)
	if err != nil {
		log.Panic(err)
	}
	println("LED on")
	_ = led.Set(true)
	time.Sleep(time.Duration(1) * time.Second)
	println("LED off")
	_ = led.Set(false)
}
