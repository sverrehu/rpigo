package component

import "github.com/sverrehu/rpigo/gpio"

type Component struct {
	Pin *Pin
}

type OutputComponent struct {
	Component
}

type BooleanOutputComponent struct {
	OutputComponent
}

type InputComponent struct {
	Component
}

type BooleanInputComponent struct {
	InputComponent
}

func NewBooleanOutputComponent(pin *Pin) (*BooleanOutputComponent, error) {
	c, err := NewOutputComponent(pin)
	if err != nil {
		return nil, err
	}
	return &BooleanOutputComponent{OutputComponent: *c}, nil
}

func NewOutputComponent(pin *Pin) (*OutputComponent, error) {
	c, err := NewComponent(pin)
	if err != nil {
		return nil, err
	}
	ls := &gpio.LineSettings{
		Direction: gpio.LineDirectionOutput,
	}
	err = pin.setup(ls)
	if err != nil {
		return nil, err
	}
	return &OutputComponent{Component: *c}, nil
}

func NewComponent(pin *Pin) (*Component, error) {
	return &Component{Pin: pin}, nil
}

func NewBooleanInputComponent(pin *Pin) (*BooleanInputComponent, error) {
	return NewBooleanInputComponentWithCallback(pin, nil)
}

func NewBooleanInputComponentWithCallback(pin *Pin, callback func(raising bool)) (*BooleanInputComponent, error) {
	c, err := NewInputComponent(pin)
	if err != nil {
		return nil, err
	}
	ls := &gpio.LineSettings{
		Direction: gpio.LineDirectionInput,
		Bias:      gpio.LineBiasPullUp,
		ActiveLow: true,
	}
	if callback != nil {
		ls.EdgeDetection = gpio.LineEdgeBoth
	}
	err = pin.setup(ls)
	if err != nil {
		return nil, err
	}
	if callback != nil {
		bridge := newCallbackBridge(pin, callback)
		go bridge.run()
	}
	return &BooleanInputComponent{InputComponent: *c}, nil
}

func NewInputComponent(pin *Pin) (*InputComponent, error) {
	c, err := NewComponent(pin)
	if err != nil {
		return nil, err
	}
	return &InputComponent{Component: *c}, nil
}

func (c *Component) Close() error {
	if c.Pin == nil {
		return nil
	}
	err := c.Pin.Close()
	c.Pin = nil
	return err
}

func (c *Component) IsClosed() bool {
	return c.Pin == nil
}
