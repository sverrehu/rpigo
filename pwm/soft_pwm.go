package pwm

import (
	"fmt"
	"log"
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

type SoftPWM struct {
	line      *gpio.Line
	Frequency float64
	Period    time.Duration // nanoseconds per Period
	dutyCycle float64       // [0-100]
	TimeOn    time.Duration
	TimeOff   time.Duration
	stop      bool
	startTime time.Time
	periods   uint64
}

func NewSoftPWM(line *gpio.Line, frequency float64) (*SoftPWM, error) {
	pwm := &SoftPWM{
		line:      line,
		Frequency: frequency,
		Period:    time.Duration(1e9 / frequency),
		stop:      false,
		periods:   0,
	}
	ls := &gpio.LineSettings{
		Direction: gpio.LineDirectionOutput,
	}
	err := line.Setup(ls, "SoftPWM")
	if err != nil {
		return nil, err
	}
	err = pwm.SetDutyCycle(0)
	if err != nil {
		return nil, err
	}
	go pwm.run()
	return pwm, nil
}

func (p *SoftPWM) Close() error {
	p.stop = true
	return p.line.Release()
}

func (p *SoftPWM) SetDutyCycle(dutyCycle float64) error {
	if dutyCycle < 0 || dutyCycle > 100 {
		return fmt.Errorf("duty cycle must be between 0 and 100, keeping it at %d", p.dutyCycle)
	}
	p.dutyCycle = dutyCycle
	p.TimeOn = time.Duration(float64(p.Period) * dutyCycle / 100)
	p.TimeOff = p.Period - p.TimeOn
	return nil
}

func (p *SoftPWM) GetPeriod() time.Duration {
	return p.Period
}

func (p *SoftPWM) setValue(value gpio.LineValue) error {
	if p.line == nil {
		return nil
	}
	return p.line.SetValue(value)
}

func (p *SoftPWM) run() {
	go p.dumpStatsForever()
	err := error(nil)
	p.startTime = time.Now()
	for !p.stop {
		if p.dutyCycle > 0 {
			startTime := time.Now()
			err = p.setValue(gpio.LineValueActive)
			d := time.Since(startTime)
			time.Sleep(p.TimeOn - d)
		}
		if err == nil && p.dutyCycle < 100 {
			startTime := time.Now()
			err = p.setValue(gpio.LineValueInactive)
			d := time.Since(startTime)
			time.Sleep(p.TimeOff - d)
		}
		if err != nil {
			log.Print("PWM error: ", err)
			p.stop = true
		}
		p.periods++
	}
	_ = p.setValue(gpio.LineValueInactive)
}

func (p *SoftPWM) dumpStatsForever() {
	for !p.stop {
		time.Sleep(5 * time.Second)
		p.dumpStats()
	}
}

func (p *SoftPWM) dumpStats() {
	frequency := float64(p.periods) / time.Since(p.startTime).Seconds()
	log.Printf("period: %s, periods: %d, time since start: %s, frequency: %f Hz (wanted %f Hz)", p.Period, p.periods, time.Since(p.startTime), frequency, p.Frequency)
}
