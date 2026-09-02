package component

import (
	"log"

	"github.com/sverrehu/rpigo/gpio"
)

// TODO: test

type Rotary struct {
	InputComponent // used for swPin, which is a button
	Pressed        bool
	Steps          int
	clkPin         *Pin
	dtPin          *Pin
}

type RotaryDirection int

const (
	RotaryClockwise RotaryDirection = iota
	RotaryCounterClockwise
)

func NewRotary(clkPin, dtPin, swPin *Pin) (*Rotary, error) {
	return NewRotaryWithCallback(clkPin, dtPin, swPin, nil, nil)
}

func NewRotaryWithCallback(clkPin, dtPin, swPin *Pin, rotCallback func(dir RotaryDirection), buttonCallback func(raising bool)) (*Rotary, error) {
	c, err := NewInputComponent(swPin)
	if err != nil {
		return nil, err
	}
	lineSettings := gpio.LineSettings{
		Direction: gpio.LineDirectionInput,
		Bias:      gpio.LineBiasPullUp,
		ActiveLow: true,
	}
	err = dtPin.setup(&lineSettings)
	if err != nil {
		return nil, err
	}
	lineSettings.EdgeDetection = gpio.LineEdgeBoth
	err = clkPin.setup(&lineSettings)
	if err != nil {
		return nil, err
	}
	if swPin != nil {
		err = swPin.setup(&lineSettings)
		if err != nil {
			return nil, err
		}
	}
	r := &Rotary{InputComponent: *c, clkPin: clkPin, dtPin: dtPin}
	if swPin != nil {
		err = r.setupButtonCallback(buttonCallback)
		if err != nil {
			return nil, err
		}
	}
	err = r.setupRotationCallback(rotCallback)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rotary) setupButtonCallback(callback func(pressed bool)) error {
	bridge := newCallbackBridge(r.Component.Pin, func(raising bool) {
		r.Pressed = raising
		if callback != nil {
			callback(raising)
		}
	})
	go bridge.run()
	return nil
}

func (r *Rotary) setupRotationCallback(callback func(direction RotaryDirection)) error {
	lastClkState, err := r.clkPin.Get()
	if err != nil {
		return err
	}
	bridge := newCallbackBridge(r.clkPin, func(raising bool) {
		dtState, err := r.dtPin.Get()
		if err != nil {
			log.Printf("Error reading dtPin: %v", err)
			return
		}
		clkState, err := r.clkPin.Get()
		if err != nil {
			log.Printf("Error reading clkPin: %v", err)
			return
		}
		if clkState != lastClkState {
			var direction RotaryDirection
			if dtState != clkState {
				r.Steps++
				direction = RotaryClockwise
			} else {
				r.Steps--
				direction = RotaryCounterClockwise
			}
			lastClkState = clkState
			if callback != nil {
				callback(direction)
			}
		}
	})
	go bridge.run()
	return nil
}

func (r *Rotary) Get() (bool, error) {
	return r.Pin.Get()
}

func (r *Rotary) Close() error {
	totalError := r.InputComponent.Close()
	if r.clkPin != nil {
		err := r.clkPin.Close()
		r.clkPin = nil
		if err != nil {
			totalError = err
		}
	}
	if r.dtPin != nil {
		err := r.dtPin.Close()
		r.dtPin = nil
		if err != nil {
			totalError = err
		}
	}
	return totalError
}
