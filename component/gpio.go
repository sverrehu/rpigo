package component

import "github.com/sverrehu/rpigo/gpio"

type PinName string

const (
	ID_SDA       = "ID_SDA"
	ID_SCL       = "ID_SCL"
	GPIO2        = "GPIO2"
	GPIO3        = "GPIO3"
	GPIO4        = "GPIO4"
	GPIO5        = "GPIO5"
	GPIO6        = "GPIO6"
	GPIO7        = "GPIO7"
	GPIO8        = "GPIO8"
	GPIO9        = "GPIO9"
	GPIO10       = "GPIO10"
	GPIO11       = "GPIO11"
	GPIO12       = "GPIO12"
	GPIO13       = "GPIO13"
	GPIO14       = "GPIO14"
	GPIO15       = "GPIO15"
	GPIO16       = "GPIO16"
	GPIO17       = "GPIO17"
	GPIO18       = "GPIO18"
	GPIO19       = "GPIO19"
	GPIO20       = "GPIO20"
	GPIO21       = "GPIO21"
	GPIO22       = "GPIO22"
	GPIO23       = "GPIO23"
	GPIO24       = "GPIO24"
	GPIO25       = "GPIO25"
	GPIO26       = "GPIO26"
	GPIO27       = "GPIO27"
	HDMI_HPD_N   = "HDMI_HPD_N"
	STATUS_LED_G = "STATUS_LED_G"
	CTS0         = "CTS0"
	RTS0         = "RTS0"
	TXD0         = "TXD0"
	RXD0         = "RXD0"
	SD1_CLK      = "SD1_CLK"
	SD1_CMD      = "SD1_CMD"
	SD1_DATA0    = "SD1_DATA0"
	SD1_DATA1    = "SD1_DATA1"
	SD1_DATA2    = "SD1_DATA2"
	SD1_DATA3    = "SD1_DATA3"
	PWM0_OUT     = "PWM0_OUT"
	PWM1_OUT     = "PWM1_OUT"
	ETH_CLK      = "ETH_CLK"
	WIFI_CLK     = "WIFI_CLK"
	SDA0         = "SDA0"
	SCL0         = "SCL0"
	SMPS_SCL     = "SMPS_SCL"
	SMPS_SDA     = "SMPS_SDA"
	SD_CLK_R     = "SD_CLK_R"
	SD_CMD_R     = "SD_CMD_R"
	SD_DATA0_R   = "SD_DATA0_R"
	SD_DATA1_R   = "SD_DATA1_R"
	SD_DATA2_R   = "SD_DATA2_R"
	SD_DATA3_R   = "SD_DATA3_R"
)

type GPIO struct {
	Consumer string
	chip     *gpio.Chip
}

func NewGPIO() (*GPIO, error) {
	chip, err := gpio.NewChip()
	if err != nil {
		return nil, err
	}
	return &GPIO{Consumer: "rpigo-component", chip: chip}, nil
}

func (g *GPIO) Close() error {
	return g.chip.Close()
}

func (g *GPIO) Pin(pin PinName) (*Pin, error) {
	line, err := g.chip.FindLine(string(pin))
	if err != nil {
		return nil, err
	}
	return newPin(g, line), nil
}
