package gpio

import (
	"errors"

	"golang.org/x/sys/unix"
)

func fillChip(chip *Chip, gpioChip *gpioChipInfo) {
	chip.Name = cBytesToString(gpioChip.Name[:])
	chip.Label = cBytesToString(gpioChip.Label[:])
	chip.NumLines = int(gpioChip.NumLines)
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/line-info.c
func fillLineInfo(lineInfo *LineInfo, li *gpioV2LineInfo) {
	lineInfo.Offset = int(li.Offset)
	lineInfo.Name = cBytesToString(li.Name[:])
	lineInfo.Used = li.Flags&gpioV2LineFlagUsed != 0
	lineInfo.Consumer = cBytesToString(li.Consumer[:])
	if li.Flags&gpioV2LineFlagOutput != 0 {
		lineInfo.Direction = LineDirectionOutput
	} else {
		lineInfo.Direction = LineDirectionInput
	}
	lineInfo.ActiveLow = li.Flags&gpioV2LineFlagActiveLow != 0
	if li.Flags&gpioV2LineFlagBiasPullUp != 0 {
		lineInfo.Bias = LineBiasPullUp
	} else if li.Flags&gpioV2LineFlagBiasPullDown != 0 {
		lineInfo.Bias = LineBiasPullDown
	} else if li.Flags&gpioV2LineFlagBiasDisabled != 0 {
		lineInfo.Bias = LineBiasDisabled
	} else {
		lineInfo.Bias = LineBiasUnknown
	}
	if li.Flags&gpioV2LineFlagOpenDrain != 0 {
		lineInfo.Drive = LineDriveOpenDrain
	} else if li.Flags&gpioV2LineFlagOpenSource != 0 {
		lineInfo.Drive = LineDriveOpenSource
	} else {
		lineInfo.Drive = LineDrivePushPull
	}
	if (li.Flags&gpioV2LineFlagEdgeRising != 0) && (li.Flags&gpioV2LineFlagEdgeFalling != 0) {
		lineInfo.Edge = LineEdgeBoth
	} else if li.Flags&gpioV2LineFlagEdgeRising != 0 {
		lineInfo.Edge = LineEdgeRising
	} else if li.Flags&gpioV2LineFlagEdgeFalling != 0 {
		lineInfo.Edge = LineEdgeFalling
	} else {
		lineInfo.Edge = LineEdgeNone
	}
	if li.Flags&gpioV2LineFlagEventClockRealtime != 0 {
		lineInfo.EventClock = LineClockRealtime
	} else if li.Flags&gpioV2LineFlagEventClockHTE != 0 {
		lineInfo.EventClock = LineClockHTE
	} else {
		lineInfo.EventClock = LineClockMonotonic
	}
	lineInfo.Debounced = false
	lineInfo.DebouncePeriodUs = 0
	for i := 0; i < int(li.NumAttrs); i++ {
		attr := li.Attrs[i]
		if attr.Id == gpioV2LineAttrIdDebounce {
			lineInfo.Debounced = true
			lineInfo.DebouncePeriodUs = attr.Value
		}
	}
}

func toGpioV2LineRequest(requestConfig *RequestConfig, lineConfig *LineConfig) (*gpioV2LineRequest, error) {
	configs, err := normalizedLineConfigs(lineConfig)
	if err != nil {
		return nil, err
	}
	lr := &gpioV2LineRequest{}
	if requestConfig != nil {
		copy(lr.Consumer[:], stringToBytes(requestConfig.Consumer, gpioMaxNameSize))
		lr.EventBufferSize = uint32(requestConfig.EventBufferSize)
	}
	lr.NumLines = uint32(len(configs))
	for i, c := range configs {
		lr.Offsets[i] = uint32(c.offset)
	}
	if err := populateLineConfig(lineConfig, configs, &lr.Config); err != nil {
		return nil, err
	}
	return lr, nil
}

func normalizedLineConfigs(lineConfig *LineConfig) ([]perLineConfig, error) {
	if lineConfig == nil {
		return nil, errors.New("lineConfig cannot be nil")
	}
	var configs []perLineConfig
	if len(lineConfig.lineConfigs) > 0 {
		configs = append(configs, lineConfig.lineConfigs...)
	} else if len(lineConfig.Offsets) > 0 {
		for _, off := range lineConfig.Offsets {
			sCopy := lineConfig.Settings
			configs = append(configs, perLineConfig{
				offset:   off,
				settings: &sCopy,
			})
		}
	}
	seen := make(map[int]bool)
	for _, c := range configs {
		if seen[c.offset] {
			return nil, unix.EINVAL
		}
		seen[c.offset] = true
	}
	if len(configs) > gpioV2LinesMax {
		return nil, unix.E2BIG
	}
	return configs, nil
}

func populateLineConfig(lineConfig *LineConfig, configs []perLineConfig, uapiCfg *gpioV2LineConfig) error {
	attrIdx := 0

	// 1. Output values
	hasOutput := false
	for _, c := range configs {
		if c.settings != nil && c.settings.Direction == LineDirectionOutput {
			hasOutput = true
			break
		}
	}
	if hasOutput || len(lineConfig.outputValues) > 0 {
		if attrIdx >= gpioV2LineNumAttrsMax {
			return unix.E2BIG
		}
		var mask, vals uint64
		for i, c := range configs {
			if c.settings != nil && c.settings.Direction == LineDirectionOutput {
				lineMaskSetBit(&mask, i)
				if c.settings.OutputValue == LineValueActive {
					lineMaskSetBit(&vals, i)
				}
			}
		}
		for i, val := range lineConfig.outputValues {
			if i < len(configs) {
				lineMaskSetBit(&mask, i)
				if val == LineValueActive {
					lineMaskSetBit(&vals, i)
				} else {
					lineMaskAssignBit(&vals, i, false)
				}
			}
		}
		attr := &uapiCfg.Attrs[attrIdx]
		attrIdx++
		attr.Attr.Id = gpioV2LineAttrIdOutputValues
		attr.Attr.Value = vals
		attr.Mask = mask
	}

	// 2. Debounce periods
	var debounceDone uint64
	for i, c := range configs {
		if lineMaskTestBit(&debounceDone, i) {
			continue
		}
		lineMaskSetBit(&debounceDone, i)
		var periodI uint64
		if c.settings != nil {
			periodI = c.settings.DebouncePeriodUs
		}
		if periodI == 0 {
			continue
		}
		if attrIdx >= gpioV2LineNumAttrsMax {
			return unix.E2BIG
		}
		var mask uint64
		lineMaskSetBit(&mask, i)
		for j := i + 1; j < len(configs); j++ {
			var periodJ uint64
			if configs[j].settings != nil {
				periodJ = configs[j].settings.DebouncePeriodUs
			}
			if periodI == periodJ {
				lineMaskSetBit(&mask, j)
				lineMaskSetBit(&debounceDone, j)
			}
		}
		attr := &uapiCfg.Attrs[attrIdx]
		attrIdx++
		attr.Attr.Id = gpioV2LineAttrIdDebounce
		attr.Attr.Value = periodI
		attr.Mask = mask
	}

	// 3. Flags
	flags := make([]gpioV2LineFlag, len(configs))
	counts := make(map[gpioV2LineFlag]int)
	var maxCount int
	var defFlag gpioV2LineFlag
	for i, c := range configs {
		if c.settings != nil {
			flags[i] = lineSettingsToFlags(c.settings)
		}
		counts[flags[i]]++
		if counts[flags[i]] > maxCount {
			maxCount = counts[flags[i]]
			defFlag = flags[i]
		}
	}
	uapiCfg.Flags = defFlag

	var flagsDone uint64
	for i, fI := range flags {
		if fI == defFlag || lineMaskTestBit(&flagsDone, i) {
			continue
		}
		lineMaskSetBit(&flagsDone, i)
		if attrIdx >= gpioV2LineNumAttrsMax {
			return unix.E2BIG
		}
		var mask uint64
		lineMaskSetBit(&mask, i)
		for j := i + 1; j < len(configs); j++ {
			if flags[j] == fI {
				lineMaskSetBit(&mask, j)
				lineMaskSetBit(&flagsDone, j)
			}
		}
		attr := &uapiCfg.Attrs[attrIdx]
		attrIdx++
		attr.Attr.Id = gpioV2LineAttrIdFlags
		attr.Attr.Value = uint64(fI)
		attr.Mask = mask
	}
	uapiCfg.NumAttrs = uint32(attrIdx)
	return nil
}

func lineSettingsToFlags(s *LineSettings) gpioV2LineFlag {
	if s == nil {
		return 0
	}
	var flags gpioV2LineFlag
	switch s.Direction {
	case LineDirectionInput:
		flags |= gpioV2LineFlagInput
	case LineDirectionOutput:
		flags |= gpioV2LineFlagOutput
	}
	if s.ActiveLow {
		flags |= gpioV2LineFlagActiveLow
	}
	switch s.EdgeDetection {
	case LineEdgeRising:
		flags |= gpioV2LineFlagEdgeRising
	case LineEdgeFalling:
		flags |= gpioV2LineFlagEdgeFalling
	case LineEdgeBoth:
		flags |= gpioV2LineFlagEdgeRising | gpioV2LineFlagEdgeFalling
	}
	switch s.Drive {
	case LineDriveOpenDrain:
		flags |= gpioV2LineFlagOpenDrain
	case LineDriveOpenSource:
		flags |= gpioV2LineFlagOpenSource
	}
	switch s.Bias {
	case LineBiasPullUp:
		flags |= gpioV2LineFlagBiasPullUp
	case LineBiasPullDown:
		flags |= gpioV2LineFlagBiasPullDown
	case LineBiasDisabled:
		flags |= gpioV2LineFlagBiasDisabled
	}
	switch s.EventClock {
	case LineClockRealtime:
		flags |= gpioV2LineFlagEventClockRealtime
	case LineClockHTE:
		flags |= gpioV2LineFlagEventClockHTE
	}
	return flags
}

func toGpioV2LineConfig(lineConfig *LineConfig) (*gpioV2LineConfig, error) {
	configs, err := normalizedLineConfigs(lineConfig)
	if err != nil {
		return nil, err
	}
	uapiCfg := &gpioV2LineConfig{}
	if err := populateLineConfig(lineConfig, configs, uapiCfg); err != nil {
		return nil, err
	}
	return uapiCfg, nil
}

func stringToBytes(s string, maxBytes int) []byte {
	b := []byte(s)
	if len(b) > maxBytes {
		return b[:maxBytes]
	}
	return b
}
