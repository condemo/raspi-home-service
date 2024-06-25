package tools

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// The file with the fan speed changes folder with each reboot
// so you need to see where it is when loading the module
var fanFilePath string

type MemInfo struct {
	MemPercent string
	MemUsed    string
	MemTotal   string
}

func newMemInfo() *MemInfo {
	m, err := mem.VirtualMemory()
	checkErr(err)

	return &MemInfo{
		MemPercent: fmt.Sprint(strconv.FormatFloat(m.UsedPercent, 'f', 2, 64), "%"),
		MemUsed:    fmt.Sprint(strconv.FormatUint(m.Used/1024/1024, 10), "mb"),
		MemTotal:   fmt.Sprint(strconv.FormatUint(m.Total/1024/1024, 10), "mb"),
	}
}

type USBDrive struct {
	Name  string
	Used  string
	Total string
}

type DiskInfo struct {
	RootUsed  string
	RootTotal string
	USBDrives []*USBDrive
}

func newDiskInfo() *DiskInfo {
	var dl []*USBDrive

	di := &DiskInfo{}

	parts, err := disk.Partitions(false)

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		checkErr(err)

		// TODO: Mucho codigo repetido, mejorable
		if u.Path == "/" {
			di.RootUsed = fmt.Sprint(
				strconv.FormatUint(u.Used/1024/1024/1024, 10), "GB")
			di.RootTotal = fmt.Sprint(
				strconv.FormatUint(u.Total/1024/1024/1024, 10), "GB")
		} else if strings.HasPrefix(u.Path, "/mnt/") {
			n := strings.ToTitle(strings.TrimPrefix(u.Path, "/mnt/"))
			uDrive := &USBDrive{
				Name:  n,
				Used:  fmt.Sprint(strconv.FormatUint(u.Used/1024/1024/1024, 10), "GB"),
				Total: fmt.Sprint(strconv.FormatUint(u.Total/1024/1024/1024, 10), "GB"),
			}
			dl = append(dl, uDrive)
		} else if strings.HasPrefix(u.Path, "/media/condemopi/") {
			n := fmt.Sprint(
				"USB ", strings.ToTitle(strings.TrimPrefix(u.Path, "/media/condemopi/")))
			uDrive := &USBDrive{
				Name:  n,
				Used:  fmt.Sprint(strconv.FormatUint(u.Used/1024/1024/1024, 10), "GB"),
				Total: fmt.Sprint(strconv.FormatUint(u.Total/1024/1024/1024, 10), "GB"),
			}
			dl = append(dl, uDrive)
		}
	}

	di.USBDrives = dl

	checkErr(err)

	return di
}

type CpuInfo struct {
	CpuTemp      string
	CoreInfoList []string
}

func newCpuInfo() *CpuInfo {
	cpuPer, err := cpu.Percent(0, true)
	checkErr(err)

	c := &CpuInfo{}

	for _, cpu := range cpuPer {
		c.CoreInfoList = append(
			c.CoreInfoList,
			fmt.Sprintf("%s%%", strconv.FormatFloat(cpu, 'f', 2, 64)))
	}

	f, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
	checkErr(err)
	defer f.Close()

	cb, err := io.ReadAll(f)
	checkErr(err)

	cs := strings.TrimSuffix(string(cb), "\n")

	cTemp, err := strconv.ParseInt(cs, 10, 64)
	checkErr(err)

	c.CpuTemp = fmt.Sprintf("%d°C", cTemp/1000)

	return c
}

type FanInfo struct {
	FanSpeed  string
	FanStatus bool
}

func newFanInfo() *FanInfo {
	f := &FanInfo{}

	stateF, err := os.Open("/sys/devices/virtual/thermal/cooling_device0/cur_state")
	checkErr(err)
	defer stateF.Close()

	s, err := io.ReadAll(stateF)
	checkErr(err)

	if state, _ := strconv.ParseInt(
		strings.TrimSuffix(string(s), "\n"), 10, 8); state >= 1 {
		f.FanStatus = true
	}

	ff, err := os.Open(fanFilePath)
	checkErr(err)
	defer ff.Close()

	fb, err := io.ReadAll(ff)
	checkErr(err)

	fanSpeed := strings.TrimSuffix(string(fb), "\n")
	checkErr(err)

	f.FanSpeed = fanSpeed

	return f
}

type NetInfo struct {
	NetUp   string
	NetDown string
}

func newNetInfo() *NetInfo {
	nf, err := os.Open("/home/condemopi/scripts/custom_output/net_speed")
	checkErr(err)

	nb, err := io.ReadAll(nf)
	checkErr(err)

	netSpeed := strings.TrimSuffix(string(nb), "\n")
	speedSlice := strings.Split(netSpeed, "|")

	return &NetInfo{
		NetUp:   speedSlice[0],
		NetDown: speedSlice[1],
	}
}

type SysInfo struct {
	*DiskInfo
	*MemInfo
	*CpuInfo
	*FanInfo
	*NetInfo
}

func NewSysInfo() *SysInfo {
	return &SysInfo{
		newDiskInfo(),
		newMemInfo(),
		newCpuInfo(),
		newFanInfo(),
		newNetInfo(),
	}
}

func (s *SysInfo) Update() {
	s.CpuInfo = newCpuInfo()
	s.MemInfo = newMemInfo()
	s.DiskInfo = newDiskInfo()
	s.FanInfo = newFanInfo()
	s.NetInfo = newNetInfo()
}

// this init check where the fan1_input file is located
func init() {
	p := "/sys/devices/platform/cooling_fan/hwmon/"
	items, err := os.ReadDir(p)
	if err != nil {
		log.Fatalf("error leyendo Items %s", err)
	}

	for _, item := range items {
		if item.IsDir() {
			subItems, err := os.ReadDir(p + item.Name())
			if err != nil {
				log.Fatalf("error leyendo subItems %s", err)
			}
			for _, subItem := range subItems {
				if !subItem.IsDir() {
					if subItem.Name() == "fan1_input" {
						fanFilePath = p + item.Name() + "/" + subItem.Name()
					}
				}
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("sys_info error: %v", err)
	}
}
