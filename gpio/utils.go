package gpio

import (
	"errors"
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func pollFd(fd int, timeout time.Duration) (bool, error) {
	var timeoutMs int
	if timeout < 0 {
		timeoutMs = -1
	} else {
		timeoutMs = int(timeout.Milliseconds())
	}
	pfds := []unix.PollFd{
		{
			Fd:     int32(fd),
			Events: unix.POLLIN | unix.POLLPRI,
		},
	}
	for {
		n, err := unix.Poll(pfds, timeoutMs)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// NOTE: the timeoutMs logic is wrong when this happens.
				// Should account for the time elapsed since the last poll.
				log.Print("Got EINTR. Retrying Poll...")
				continue
			}
			return false, err
		}
		return n > 0, nil
	}
}

func readInfoEventFd(fd int) (*InfoEvent, error) {
	var chg gpioV2LineInfoChanged
	buf := (*[unsafe.Sizeof(chg)]byte)(unsafe.Pointer(&chg))[:]
	n, err := unix.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if uintptr(n) < unsafe.Sizeof(chg) {
		return nil, unix.EIO
	}
	var lineInfo LineInfo
	fillLineInfo(&lineInfo, &chg.Info)
	return &InfoEvent{
		EventType:   InfoEventType(chg.EventType),
		TimestampNs: chg.TimestampNs,
		LineInfo:    &lineInfo,
	}, nil
}

func lineMaskTestBit(mask *uint64, nr int) bool {
	return (*mask & (1 << nr)) != 0
}

func lineMaskSetBit(mask *uint64, nr int) {
	*mask |= 1 << nr
}

func lineMaskAssignBit(mask *uint64, nr int, val bool) {
	if val {
		lineMaskSetBit(mask, nr)
	} else {
		*mask &= ^(1 << nr)
	}
}
