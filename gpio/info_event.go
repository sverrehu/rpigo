package gpio

func (e *InfoEvent) Copy() *InfoEvent {
	cpy := *e
	if e.LineInfo != nil {
		cpy.LineInfo = e.LineInfo.Copy()
	}
	return &cpy
}
