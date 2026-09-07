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
	button, err := component.NewButtonWithCallback(pin, func(raising bool) {
		if raising {
			println("Button pressed")
		} else {
			println("Button released")
		}
	})
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(time.Duration(10) * time.Hour)
	_ = button.Close()
}
