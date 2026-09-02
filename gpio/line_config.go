package gpio

import (
	"slices"

	"golang.org/x/sys/unix"
)

func NewLineConfig() *LineConfig {
	return &LineConfig{}
}

func (c *LineConfig) AddLineSettings(offsets []int, settings *LineSettings) error {
	if settings == nil {
		return unix.EINVAL
	}
	if len(c.lineConfigs)+len(offsets) > gpioV2LinesMax {
		return unix.E2BIG
	}
	seen := make(map[int]bool)
	for _, lc := range c.lineConfigs {
		seen[lc.offset] = true
	}
	for _, off := range offsets {
		if seen[off] {
			return unix.EINVAL
		}
		seen[off] = true
	}
	for _, off := range offsets {
		c.lineConfigs = append(c.lineConfigs, perLineConfig{
			offset:   off,
			settings: settings.Copy(),
		})
		c.Offsets = append(c.Offsets, off)
	}
	return nil
}

func (c *LineConfig) GetLineSettings(offset int) (*LineSettings, error) {
	for _, lc := range c.lineConfigs {
		if lc.offset == offset {
			return lc.settings.Copy(), nil
		}
	}
	if len(c.lineConfigs) == 0 && len(c.Offsets) > 0 {
		if slices.Contains(c.Offsets, offset) {
			return c.Settings.Copy(), nil
		}
	}
	return nil, unix.ENOENT
}

func (c *LineConfig) SetOutputValues(values []LineValue) error {
	if len(values) > gpioV2LinesMax {
		return unix.E2BIG
	}
	for _, val := range values {
		if val != LineValueInactive && val != LineValueActive {
			return unix.EINVAL
		}
	}
	c.outputValues = make([]LineValue, len(values))
	copy(c.outputValues, values)
	return nil
}

func (c *LineConfig) SetOutputValueOffset(index int, value LineValue) error {
	if index < 0 || index >= gpioV2LinesMax {
		return unix.EINVAL
	}
	if value != LineValueInactive && value != LineValueActive {
		return unix.EINVAL
	}
	if index >= len(c.outputValues) {
		newVals := make([]LineValue, index+1)
		copy(newVals, c.outputValues)
		c.outputValues = newVals
	}
	c.outputValues[index] = value
	return nil
}

func (c *LineConfig) ConfiguredOffsets() []int {
	if len(c.lineConfigs) > 0 {
		offs := make([]int, len(c.lineConfigs))
		for i, lc := range c.lineConfigs {
			offs[i] = lc.offset
		}
		return offs
	}
	return append([]int(nil), c.Offsets...)
}

func (c *LineConfig) NumConfiguredOffsets() int {
	return len(c.ConfiguredOffsets())
}
