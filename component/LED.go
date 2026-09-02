package component

type LED struct {
	BooleanOutputComponent
}

func NewLED(pin *Pin) (*LED, error) {
	c, err := NewBooleanOutputComponent(pin)
	if err != nil {
		return nil, err
	}
	return &LED{BooleanOutputComponent: *c}, nil
}

func (l *LED) Set(b bool) error {
	return l.Pin.Set(b)
}
