package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
)

func main() {
	const clkPinName = component.GPIO17
	const dtPinName = component.GPIO27
	const swPinName = component.GPIO22
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	clkPin, err := gpio.Pin(clkPinName)
	if err != nil {
		log.Panic(err)
	}
	dtPin, err := gpio.Pin(dtPinName)
	if err != nil {
		log.Panic(err)
	}
	swPin, err := gpio.Pin(swPinName)
	if err != nil {
		log.Panic(err)
	}
	rotary, err := component.NewRotary(clkPin, dtPin, swPin)
	if err != nil {
		log.Panic(err)
	}
	for {
		buttonStatus := "released"
		if rotary.Pressed {
			buttonStatus = "pressed"
		}
		fmt.Printf("Button %s, steps: %d\n", buttonStatus, rotary.Steps)
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
