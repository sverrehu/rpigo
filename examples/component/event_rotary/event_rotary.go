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
	rotary, err := component.NewRotaryWithCallback(clkPin, dtPin, swPin, rotationCallback, buttonCallback)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(time.Duration(10) * time.Hour)
	_ = rotary.Close()
}

func rotationCallback(direction component.RotaryDirection) {
	if direction == component.RotaryClockwise {
		fmt.Println("Rotated clockwise")
	} else {
		fmt.Println("Rotated counter-clockwise")
	}
}

func buttonCallback(pressed bool) {
	if pressed {
		fmt.Println("Button pressed")
	} else {
		fmt.Println("Button released")
	}
}
