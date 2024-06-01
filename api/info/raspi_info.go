package info

import (
	"encoding/json"
	"os/exec"
)

type Rasp5Temp struct {
	Cpu struct {
		Tctl struct {
			Temp float32 `json:"temp1_input"`
		} `json:"temp1"`
	} `json:"cpu_thermal-virtual-0"`
	Fan struct {
		Fan1 struct {
			Speed float32 `json:"fan1_input"`
		} `json:"fan1"`
	} `json:"pwmfan-isa-0000"`
}

func RaspberryTemps() *Rasp5Temp {
	out, _ := exec.Command("sensors", "-j").Output()

	r := Rasp5Temp{}
	json.Unmarshal(out, &r)

	return &r
}
