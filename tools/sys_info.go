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

type DiskInfo struct {
	ToshibaUsed  string
	ToshibaTotal string
	SeagateUsed  string
	SeagateTotal string
	RootUsed     string
	RootTotal    string
}

func newDiskInfo() *DiskInfo {
	var rootPart *disk.UsageStat
	var mntToshiba *disk.UsageStat
	var mntSeagate *disk.UsageStat

	parts, err := disk.Partitions(false)

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		checkErr(err)

		if u.Path == "/" {
			rootPart = u
		} else if strings.HasPrefix(u.Path, "/mnt/toshiba") {
			mntToshiba = u
		} else if strings.HasPrefix(u.Path, "/mnt/seagate") {
			mntSeagate = u
		}
	}
	checkErr(err)

	return &DiskInfo{
		ToshibaUsed: fmt.Sprint(
			strconv.FormatUint(mntToshiba.Used/1024/1024/1024, 10), "GB"),
		ToshibaTotal: fmt.Sprint(
			strconv.FormatUint(mntToshiba.Total/1024/1024/1024, 10), "GB"),
		SeagateUsed: fmt.Sprint(
			strconv.FormatUint(mntSeagate.Used/1024/1024/1024, 10), "GB"),
		SeagateTotal: fmt.Sprint(
			strconv.FormatUint(mntSeagate.Total/1024/1024/1024, 10), "GB"),
		RootUsed: fmt.Sprint(
			strconv.FormatUint(rootPart.Used/1024/1024/1024, 10), "GB"),
		RootTotal: fmt.Sprint(strconv.FormatUint(rootPart.Total/1024/1024/1024, 10), "GB"),
	}
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

	speedF, err := os.Open("/sys/devices/platform/cooling_fan/hwmon/hwmon1/fan1_input")
	checkErr(err)
	defer stateF.Close()

	fs, err := io.ReadAll(speedF)
	checkErr(err)

	fanSpeed := strings.TrimSuffix(string(fs), "\n")
	checkErr(err)

	f.FanSpeed = fanSpeed

	return f
}

type SysInfo struct {
	*DiskInfo
	*MemInfo
	*CpuInfo
	*FanInfo
}

func NewSysInfo() *SysInfo {
	return &SysInfo{
		newDiskInfo(),
		newMemInfo(),
		newCpuInfo(),
		newFanInfo(),
	}
}

func (s *SysInfo) Update() {
	s.CpuInfo = newCpuInfo()
	s.MemInfo = newMemInfo()
	s.DiskInfo = newDiskInfo()
	s.FanInfo = newFanInfo()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
