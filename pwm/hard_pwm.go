package pwm

// https://www.kernel.org/doc/html/v7.2/driver-api/pwm.html#using-pwms-with-the-sysfs-interface
// pinctl | grep PWM
// Need to have 'dtoverlay=pwm-2chan' under '[all]' in /boot/firmware/config.txt
// otherwise you will get "open /sys/class/pwm/pwmchip0/export: no such file or directory"

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const apiRoot = "/sys/class/pwm"

type HardPWM struct {
	ChipNumber int
	PwmChannel int
	Frequency  float64
	Period     time.Duration // nanoseconds per Period
	timeOn     time.Duration
	dutyCycle  float64 // [0-100]
	chipRoot   string
	pwmRoot    string
}

func NewHardPWM(chipNumber, pwmChannel int, frequency float64) (*HardPWM, error) {
	pwm := &HardPWM{
		ChipNumber: chipNumber,
		PwmChannel: pwmChannel,
		Frequency:  frequency,
		Period:     time.Duration(1e9 / frequency),
		chipRoot:   fmt.Sprintf("%s/pwmchip%d", apiRoot, chipNumber),
		pwmRoot:    fmt.Sprintf("%s/pwmchip%d/pwm%d", apiRoot, chipNumber, pwmChannel),
	}
	err := pwm.export()
	if err != nil {
		return nil, err
	}
	err = pwm.setPeriod(pwm.Period.Nanoseconds())
	if err != nil {
		return nil, err
	}
	err = pwm.SetDutyCycle(0)
	if err != nil {
		return nil, err
	}
	return pwm, pwm.enable()
}

func (p *HardPWM) Close() error {
	err := p.disable()
	if err != nil {
		return err
	}
	return p.unexport()
}

func (p *HardPWM) SetDutyCycle(dutyCycle float64) error {
	if dutyCycle < 0 || dutyCycle > 100 {
		return fmt.Errorf("duty cycle must be between 0 and 100, keeping it at %d", p.dutyCycle)
	}
	p.dutyCycle = dutyCycle
	p.timeOn = time.Duration(float64(p.Period) * dutyCycle / 100)
	return p.setDutyCycle(p.timeOn.Nanoseconds())
}

func (p *HardPWM) GetPeriod() time.Duration {
	return p.Period
}

func (p *HardPWM) export() error {
	err := p.writeAPI(p.chipPath("export"), int64(p.PwmChannel))
	if errors.Is(err, unix.EBUSY) {
		log.Print("PWM already exported. Trying to go on.")
		return nil
	}
	if errors.Is(err, unix.ENOENT) {
		return fmt.Errorf("cannot initiate hardware PWM. make sure you have 'dtoverlay=pwm-2chan' under '[all]' in /boot/firmware/config.txt: %w", err)
	}
	if err == nil {
		time.Sleep(100 * time.Millisecond) // wait for subdir to settle
	}
	return err
}

func (p *HardPWM) unexport() error {
	return p.writeAPI(p.chipPath("unexport"), int64(p.PwmChannel))
}

func (p *HardPWM) enable() error {
	return p.writeAPI(p.pwmPath("enable"), 1)
}

func (p *HardPWM) disable() error {
	return p.writeAPI(p.pwmPath("enable"), 0)
}

func (p *HardPWM) setPeriod(periodNs int64) error {
	return p.writeAPI(p.pwmPath("period"), periodNs)
}

func (p *HardPWM) setDutyCycle(dutyCycleNs int64) error {
	return p.writeAPI(p.pwmPath("duty_cycle"), dutyCycleNs)
}

func (p *HardPWM) chipPath(filename string) string {
	return fmt.Sprintf("%s/%s", p.chipRoot, filename)
}

func (p *HardPWM) pwmPath(filename string) string {
	return fmt.Sprintf("%s/%s", p.pwmRoot, filename)
}

func (p *HardPWM) writeAPI(filename string, value int64) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, err = f.WriteString(fmt.Sprintf("%d\n", value))
	return err
}
