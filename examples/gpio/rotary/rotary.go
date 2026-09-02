package main

import (
	"log"

	"github.com/sverrehu/rpigo/gpio"
)

func main() {
	const clkPinOffset = 17
	const dtPinOffset = 27
	const swPinOffset = 22
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
		//DebouncePeriodUs: 20000,
	}
	requestConfig := gpio.RequestConfig{
		Consumer: "rotary-example",
	}
	lineConfig := gpio.LineConfig{
		Offsets:  []int{clkPinOffset, dtPinOffset, swPinOffset},
		Settings: lineSettings,
	}
	request, err := chip.RequestLines(&requestConfig, &lineConfig)
	if err != nil {
		log.Panic(err)
	}
	edgeEventBuffer := gpio.NewEdgeEventBuffer(-1)
	lineValue, _ := request.GetValue(clkPinOffset)
	lastClkState := lineValue == gpio.LineValueActive
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
			if edgeEvent.LineOffset == swPinOffset {
				rising := edgeEvent.EventType == gpio.EdgeEventRisingEdge
				if rising {
					println("Button pressed")
				} else {
					println("Button released")
				}
			} else {
				if edgeEvent.LineOffset == clkPinOffset {
					lineValue, _ := request.GetValue(dtPinOffset)
					dtState := lineValue == gpio.LineValueActive
					lineValue, _ = request.GetValue(clkPinOffset)
					clkState := lineValue == gpio.LineValueActive
					if clkState != lastClkState {
						if dtState != clkState {
							println("Rotated right")
						} else {
							println("Rotated left")
						}
						lastClkState = clkState
					}
				}
			}
		}
	}
}
