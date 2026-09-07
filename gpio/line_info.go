package gpio

func (li *LineInfo) Copy() *LineInfo {
	cpy := *li
	return &cpy
}
