package component

type Button struct {
	BooleanInputComponent
}

func NewButton(pin *Pin) (*Button, error) {
	c, err := NewBooleanInputComponent(pin)
	if err != nil {
		return nil, err
	}
	return &Button{BooleanInputComponent: *c}, nil
}

func NewButtonWithCallback(pin *Pin, callback func(raising bool)) (*Button, error) {
	c, err := NewBooleanInputComponentWithCallback(pin, callback)
	if err != nil {
		return nil, err
	}
	return &Button{BooleanInputComponent: *c}, nil
}

func (b *Button) Get() (bool, error) {
	return b.Pin.Get()
}
