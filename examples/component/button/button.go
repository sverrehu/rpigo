package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
)

func main() {
	const pinName = component.GPIO23
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	pin, err := gpio.Pin(pinName)
	if err != nil {
		log.Panic(err)
	}
	button, err := component.NewButton(pin)
	if err != nil {
		log.Panic(err)
	}
	for {
		value, err := button.Get()
		if err != nil {
			log.Panic(err)
		}
		if value {
			println("Button pressed")
		} else {
			println("Button released")
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
