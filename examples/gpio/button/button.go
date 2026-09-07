package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

func main() {
	const pinOffset = 23
	chip, err := gpio.NewChip()
	if err != nil {
		log.Panic(err)
	}
	defer chip.Close()
	lineSettings := gpio.LineSettings{
		Direction: gpio.LineDirectionInput,
		Bias:      gpio.LineBiasPullUp,
		ActiveLow: true,
	}
	requestConfig := gpio.RequestConfig{
		Consumer: "button-example",
	}
	lineConfig := gpio.LineConfig{
		Offsets:  []int{pinOffset},
		Settings: lineSettings,
	}
	request, err := chip.RequestLines(&requestConfig, &lineConfig)
	if err != nil {
		log.Panic(err)
	}
	for {
		value, err := request.GetValue(pinOffset)
		if err != nil {
			log.Panic(err)
		}
		if value == gpio.LineValueActive {
			println("Button pressed")
		} else {
			println("Button released")
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
