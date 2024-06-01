package info

import (
	"encoding/json"
	"os/exec"
)

type Rasp5Temp struct {
	Cpu struct {
		Tctl struct {
			Temp float32 `json:"temp1_input"`
		} `json:"tctl"`
	} `json:"k10temp-pci-00c3"`
	Gpu struct {
		Edge struct {
			Temp float32 `json:"temp1_input"`
		} `json:"edge"`
	} `json:"amdgpu-pci-0400"`
}

func RaspberryTemps() *Rasp5Temp {
	out, _ := exec.Command("sensors", "-j").Output()

	r := Rasp5Temp{}
	json.Unmarshal(out, &r)

	return &r
}
