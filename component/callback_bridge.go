package component

import (
	"errors"
	"log"
	"time"

	"github.com/sverrehu/rpigo/gpio"
	"golang.org/x/sys/unix"
)

const waitTimeout = time.Duration(500 * time.Millisecond)

type callbackBridge struct {
	pin      *Pin
	callback func(raising bool)
}

func newCallbackBridge(pin *Pin, callback func(raising bool)) *callbackBridge {
	return &callbackBridge{
		pin:      pin,
		callback: callback,
	}
}

func (cb *callbackBridge) run() {
	edgeEventBuffer := gpio.NewEdgeEventBuffer(-1)
	for {
		if cb.pin.IsClosed() {
			break
		}
		request := cb.pin.line.Req
		gotEvent, err := request.WaitEdgeEvents(waitTimeout)
		if err != nil {
			log.Printf("Got error during wait: %v", err)
			if errors.Is(err, unix.EINTR) {
				log.Printf("Got EINTR during wait")
			}
			log.Panic(err)
		}
		if cb.pin.IsClosed() {
			break
		}
		if !gotEvent {
			continue
		}
		numEvents, err := request.ReadEdgeEvents(edgeEventBuffer, 10)
		if err != nil {
			log.Printf("Got error during read: %v", err)
			if errors.Is(err, unix.EINTR) {
				log.Printf("Got EINTR during read")
			}
			log.Panic(err)
		}
		for i := range numEvents {
			edgeEvent := edgeEventBuffer.Events[i]
			value := false
			if edgeEvent.EventType == gpio.EdgeEventRisingEdge {
				value = true
			}
			cb.callback(value)
		}
	}
	log.Printf("Callback bridge stopped")
}
