package gpio

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

const eventBufferMaxCapacity = gpioV2LinesMax * 16

func (e *EdgeEvent) Copy() *EdgeEvent {
	cpy := *e
	return &cpy
}

func NewEdgeEventBuffer(capacity int) *EdgeEventBuffer {
	if capacity <= 0 {
		capacity = 64
	}
	if capacity > eventBufferMaxCapacity {
		capacity = eventBufferMaxCapacity
	}
	return &EdgeEventBuffer{
		Capacity:  capacity,
		Events:    make([]EdgeEvent, 0, capacity),
		eventData: make([]gpioV2LineEvent, capacity),
	}
}

func (b *EdgeEventBuffer) ReadFd(fd int, maxEvents int) (int, error) {
	if fd < 0 {
		return 0, unix.EINVAL
	}
	if maxEvents <= 0 || maxEvents > b.Capacity {
		maxEvents = b.Capacity
	}
	if maxEvents == 0 {
		b.Events = b.Events[:0]
		return 0, nil
	}
	elemSize := int(unsafe.Sizeof(gpioV2LineEvent{}))
	buf := unsafe.Slice((*byte)(unsafe.Pointer(&b.eventData[0])), maxEvents*elemSize)
	n, err := unix.Read(fd, buf)
	if err != nil {
		return 0, err
	}
	if n < elemSize {
		return 0, unix.EIO
	}
	numRead := n / elemSize
	b.Events = b.Events[:0]
	for i := range numRead {
		raw := &b.eventData[i]
		evType := EdgeEventRisingEdge
		if raw.Id == uint32(gpioV2LineEventFallingEdge) {
			evType = EdgeEventFallingEdge
		}
		b.Events = append(b.Events, EdgeEvent{
			EventType:   evType,
			TimestampNs: raw.TimestampNs,
			LineOffset:  int(raw.Offset),
			GlobalSeqno: uint64(raw.SeqNo),
			LineSeqno:   uint64(raw.LineSeqNo),
		})
	}
	return numRead, nil
}
