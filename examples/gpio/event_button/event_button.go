package main

import (
	"log"

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
		Direction:     gpio.LineDirectionInput,
		Bias:          gpio.LineBiasPullUp,
		ActiveLow:     true,
		EdgeDetection: gpio.LineEdgeBoth,
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
	edgeEventBuffer := gpio.NewEdgeEventBuffer(-1)
	for {
		gotEvent, err := request.WaitEdgeEvents(-1)
		if err != nil {
			log.Panic(err)
		}
		if !gotEvent {
			continue
		}
		numEvents, err := request.ReadEdgeEvents(edgeEventBuffer, 10)
		if err != nil {
			log.Panic(err)
		}
		for i := range numEvents {
			edgeEvent := edgeEventBuffer.Events[i]
			if edgeEvent.EventType == gpio.EdgeEventRisingEdge {
				println("Button pressed")
			} else {
				println("Button released")
			}
		}
	}
}
