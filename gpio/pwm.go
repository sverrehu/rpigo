package gpio

import (
	"log"
	"time"
)

type SoftPWM struct {
	line      *Line
	Frequency float64
	Period    time.Duration // nanoseconds per Period
	dutyCycle int           // [0-100]
	timeOn    time.Duration
	timeOff   time.Duration
	stop      bool
	startTime time.Time
	periods   uint64
}

func NewSoftPWM(line *Line, frequency float64) *SoftPWM {
	pwm := &SoftPWM{
		line:      line,
		Frequency: frequency,
		Period:    time.Duration(1e9 / frequency),
		stop:      false,
		periods:   0,
	}
	pwm.SetDutyCycle(0)
	go pwm.run()
	return pwm
}

func (p *SoftPWM) Close() {
	p.stop = true
}

func (p *SoftPWM) SetDutyCycle(dutyCycle int) {
	if dutyCycle < 0 || dutyCycle > 100 {
		log.Printf("duty cycle must be between 0 and 100, keeping it at %d", p.dutyCycle)
		return
	}
	p.dutyCycle = dutyCycle
	p.timeOn = time.Duration(float64(p.Period) * float64(dutyCycle) / 100)
	p.timeOff = p.Period - p.timeOn
}

func (p *SoftPWM) setValue(value LineValue) error {
	if p.line == nil {
		return nil
	}
	return p.line.SetValue(value)
}

func (p *SoftPWM) run() {
	go p.dumpStatsForever()
	err := error(nil)
	delay := time.Duration(0)
	p.startTime = time.Now()
	for !p.stop {
		startTime := time.Now()
		timeAdjust := delay
		if p.dutyCycle > 0 && p.dutyCycle < 100 {
			timeAdjust = timeAdjust / 2
		}
		if p.dutyCycle > 0 {
			err = p.setValue(LineValueActive)
			time.Sleep(p.timeOn + timeAdjust)
		}
		if err == nil && p.dutyCycle < 100 {
			err = p.setValue(LineValueInactive)
			time.Sleep(p.timeOff + timeAdjust)
		}
		if err != nil {
			log.Print("PWM error: ", err)
			p.stop = true
		}
		p.periods++
		delay = p.Period - time.Since(startTime)
	}
	_ = p.setValue(LineValueInactive)
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
