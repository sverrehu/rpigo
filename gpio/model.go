package gpio

import (
	"fmt"
	"os"
)

type Chip struct {
	f        *os.File
	Path     string
	Name     string
	Label    string
	NumLines int
}

type Line struct {
	Chip   *Chip
	Offset int
	Info   *LineInfo
	Req    *LineRequest
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/include/gpiod.h
type LineValue int

const (
	LineValueError LineValue = iota - 1
	LineValueInactive
	LineValueActive
)

func (lv LineValue) String() string {
	switch lv {
	case LineValueError:
		return "error"
	case LineValueInactive:
		return "inactive"
	case LineValueActive:
		return "active"
	default:
		return "unknown"
	}
}

type LineDirection int

const (
	LineDirectionAsIs LineDirection = iota + 1
	LineDirectionInput
	LineDirectionOutput
)

func (ld LineDirection) String() string {
	switch ld {
	case LineDirectionAsIs:
		return "as-is"
	case LineDirectionInput:
		return "input"
	case LineDirectionOutput:
		return "output"
	default:
		return "unknown"
	}
}

type LineEdge int

const (
	LineEdgeNone LineEdge = iota + 1
	LineEdgeRising
	LineEdgeFalling
	LineEdgeBoth
)

func (le LineEdge) String() string {
	switch le {
	case LineEdgeNone:
		return "none"
	case LineEdgeRising:
		return "rising"
	case LineEdgeFalling:
		return "falling"
	case LineEdgeBoth:
		return "both"
	default:
		return "unknown"
	}
}

type LineBias int

const (
	LineBiasAsIs LineBias = iota + 1
	LineBiasUnknown
	LineBiasDisabled
	LineBiasPullUp
	LineBiasPullDown
)

func (lb LineBias) String() string {
	switch lb {
	case LineBiasAsIs:
		return "as-is"
	case LineBiasUnknown:
		return "unknown"
	case LineBiasDisabled:
		return "disabled"
	case LineBiasPullUp:
		return "pull-up"
	case LineBiasPullDown:
		return "pull-down"
	default:
		return "unknown"
	}
}

type LineDrive int

const (
	LineDrivePushPull LineDrive = iota + 1
	LineDriveOpenDrain
	LineDriveOpenSource
)

func (ld LineDrive) String() string {
	switch ld {
	case LineDrivePushPull:
		return "push-pull"
	case LineDriveOpenDrain:
		return "open-drain"
	case LineDriveOpenSource:
		return "open-source"
	default:
		return "unknown"
	}
}

type LineClock int

const (
	LineClockMonotonic LineClock = iota + 1
	LineClockRealtime
	LineClockHTE
)

func (lc LineClock) String() string {
	switch lc {
	case LineClockMonotonic:
		return "monotonic"
	case LineClockRealtime:
		return "realtime"
	case LineClockHTE:
		return "hte"
	default:
		return "unknown"
	}
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/line-info.c
type LineInfo struct {
	Offset           int
	Name             string
	Used             bool
	Consumer         string
	Direction        LineDirection
	ActiveLow        bool
	Bias             LineBias
	Drive            LineDrive
	Edge             LineEdge
	EventClock       LineClock
	Debounced        bool
	DebouncePeriodUs uint64
}

func (li LineInfo) String() string {
	return fmt.Sprintf("offset: %d, name: %s, used: %t, consumer: %s, direction: %s, active-low: %t, bias: %s, drive: %s, edge: %s, event-clock: %s, debounced: %t, debounce-period-us: %d",
		li.Offset, li.Name, li.Used, li.Consumer, li.Direction, li.ActiveLow, li.Bias, li.Drive, li.Edge, li.EventClock, li.Debounced, li.DebouncePeriodUs)
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/line-settings.c
type LineSettings struct {
	Direction        LineDirection
	EdgeDetection    LineEdge
	Drive            LineDrive
	Bias             LineBias
	ActiveLow        bool
	EventClock       LineClock
	DebouncePeriodUs uint64
	OutputValue      LineValue
}

type perLineConfig struct {
	offset   int
	settings *LineSettings
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/line-config.c
type LineConfig struct {
	Offsets      []int
	Settings     LineSettings
	lineConfigs  []perLineConfig
	outputValues []LineValue
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/line-request.c
type LineRequest struct {
	ChipName string
	Offsets  [gpioV2LinesMax]uint32
	NumLines int
	Fd       int
}

// https://git.kernel.org/pub/scm/libs/libgpiod/libgpiod.git/tree/lib/request-config.c
type RequestConfig struct {
	Consumer        string
	EventBufferSize int
}

type EdgeEventType int

const (
	EdgeEventRisingEdge EdgeEventType = iota + 1
	EdgeEventFallingEdge
)

type EdgeEvent struct {
	EventType   EdgeEventType
	TimestampNs uint64
	LineOffset  int
	GlobalSeqno uint64
	LineSeqno   uint64
}

type EdgeEventBuffer struct {
	Capacity  int
	Events    []EdgeEvent
	eventData []gpioV2LineEvent
}

type InfoEventType int

const (
	InfoEventLineRequested InfoEventType = iota + 1
	InfoEventLineReleased
	InfoEventLineConfigChanged
)

type InfoEvent struct {
	EventType   InfoEventType
	TimestampNs uint64
	LineInfo    *LineInfo
}
