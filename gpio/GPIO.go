package gpio

// https://pip-assets.raspberrypi.com/categories/545-raspberry-pi-4-model-b/documents/RP-008248-DS-1-bcm2711-peripherals.pdf

import (
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

type PinMode int
type Direction int
type State int

const (
	Bcm PinMode = iota
	Board
)

const (
	NoDirection Direction = iota
	In
	Out
)

const (
	Low State = iota
	High
)

const numChannels = 40

// The following are indexes into the uint32 array, hence the official memory offsets are divided by 4.
const GPFSELStart = 0x00 >> 2
const GPSETStart = 0x1C >> 2
const GPCLRStart = 0x28 >> 2

type GPIO struct {
	pinMode    PinMode
	directions []Direction
	mmapFile   *os.File
	mmapBytes  []byte
	gpioMem    []uint32
}

func NewGPIO() *GPIO {
	g := &GPIO{
		pinMode:    Bcm,
		directions: make([]Direction, numChannels),
	}
	g.intiGPIOMmap()
	return g
}

func (g *GPIO) SetMode(mode PinMode) {
	g.assertGpioAvailable()
	if mode != Bcm {
		log.Fatal("Only BCM mode is supported in this implementation")
	}
	g.pinMode = mode
}

func (g *GPIO) Setup(pin int, direction Direction) {
	g.assertGpioAvailable()
	g.validatePin(pin)
	// There are 6 GPFSEL registers, each of which controls 10 pins. Each pin is configured using 3 bits.
	index := GPFSELStart + pin/10
	bitOffset := (pin % 10) * 3
	g.gpioMem[index] = (g.gpioMem[index] &^ (0b111 << bitOffset)) | (uint32(direction) << bitOffset)
	// TODO
	g.directions[pin] = direction
}

func (g *GPIO) Input(pin int) State {
	g.assertGpioAvailable()
	g.validatePin(pin)
	if g.directions[pin] != In {
		log.Fatal("Pin is not set to input")
	}
	// TODO
	return State(g.gpioMem[pin])
}

func (g *GPIO) Output(pin int, state State) {
	g.assertGpioAvailable()
	g.validatePin(pin)
	if g.directions[pin] != Out {
		log.Fatal("Pin is not set to output")
	}
	// There are 2 GPFSET and GPFCLR registers, each of which controls 32 pins. Each pin is configured using 1 bit.
	bitOffset := pin % 32
	var index int
	if state == Low {
		index = GPCLRStart + pin/32
	} else {
		index = GPSETStart + pin/32
	}
	g.gpioMem[index] = (g.gpioMem[index] &^ (0b1 << bitOffset)) | (uint32(state) << bitOffset)
}

func (g *GPIO) Cleanup() {
	g.closeGPIOMmap()
}

func (g *GPIO) intiGPIOMmap() {
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	g.mmapFile = f
	numElements := 1024
	fileSize := numElements * 4
	b, err := unix.Mmap(int(g.mmapFile.Fd()), 0, fileSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	g.mmapBytes = b
	g.gpioMem = unsafe.Slice((*uint32)(unsafe.Pointer(&b[0])), numElements)
}

func (g *GPIO) closeGPIOMmap() {
	if g.gpioMem == nil {
		return
	}
	if err := unix.Munmap(g.mmapBytes); err != nil {
		log.Fatal(err)
	}
	if err := g.mmapFile.Close(); err != nil {
		log.Fatal(err)
	}
	g.gpioMem = nil
}

func (g *GPIO) assertGpioAvailable() {
	if g.gpioMem == nil {
		log.Fatal("GPIO has been closed")
	}
}

func (g *GPIO) validatePin(pin int) {
	if pin < 0 || pin >= numChannels {
		log.Fatalf("Invalid pin number %d", pin)
	}
}
