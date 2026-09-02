package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

func main() {
	const pinOffset = 18
	chip, err := gpio.NewChip()
	if err != nil {
		log.Panic(err)
	}
	defer chip.Close()
	lineSettings := gpio.LineSettings{
		Direction: gpio.LineDirectionOutput,
	}
	requestConfig := gpio.RequestConfig{
		Consumer: "led-example",
	}
	lineConfig := gpio.LineConfig{
		Offsets:  []int{pinOffset},
		Settings: lineSettings,
	}
	request, err := chip.RequestLines(&requestConfig, &lineConfig)
	if err != nil {
		log.Panic(err)
	}
	println("LED on")
	request.SetValue(pinOffset, gpio.LineValueActive)
	time.Sleep(time.Duration(1) * time.Second)
	println("LED off")
	request.SetValue(pinOffset, gpio.LineValueInactive)
}
