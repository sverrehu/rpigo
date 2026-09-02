package gpio

func NewLineSettings() *LineSettings {
	return &LineSettings{
		Direction:     LineDirectionAsIs,
		EdgeDetection: LineEdgeNone,
		Drive:         LineDrivePushPull,
		Bias:          LineBiasAsIs,
		ActiveLow:     false,
		EventClock:    LineClockMonotonic,
		OutputValue:   LineValueInactive,
	}
}

func (s *LineSettings) Copy() *LineSettings {
	cpy := *s
	return &cpy
}
