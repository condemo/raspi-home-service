package tools

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type SysInfo struct {
	MemPercent   string `json:"mem-perc"`
	MemUsed      string `json:"mem-used"`
	MemTotal     string `json:"mem-total"`
	ToshibaUsed  string `json:"toshiba-used"`
	ToshibaTotal string `json:"toshiba-total"`
	SeagateUsed  string `json:"seagate-used"`
	SeagateTotal string `json:"seagate-total"`
	RootUsed     string `json:"root-used"`
	RootTotal    string `json:"root-total"`
}

func NewSysInfo() *SysInfo {
	m, err := mem.VirtualMemory()
	checkErr(err)

	parts, err := disk.Partitions(false)
	checkErr(err)

	cpuPercent, err := cpu.Percent(5, true)
	checkErr(err)

	fmt.Println(cpuPercent)

	var rootPart *disk.UsageStat
	var mntToshiba *disk.UsageStat
	var mntSeagate *disk.UsageStat
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

	return &SysInfo{
		MemPercent: fmt.Sprint(strconv.FormatFloat(m.UsedPercent, 'f', 2, 64), "%"),
		MemUsed:    fmt.Sprint(strconv.FormatUint(m.Used/1024/1024, 10), "mb"),
		MemTotal:   fmt.Sprint(strconv.FormatUint(m.Total/1024/1024, 10), "mb"),
		ToshibaUsed: fmt.Sprint(
			strconv.FormatUint(mntToshiba.Used/1024/1024/1024, 10), "GB"),
		ToshibaTotal: fmt.Sprint(strconv.FormatUint(mntToshiba.Total/1024/1024/1024, 10), "GB"),
		SeagateUsed: fmt.Sprint(
			strconv.FormatUint(mntSeagate.Used/1024/1024/1024, 10), "GB"),
		SeagateTotal: fmt.Sprint(strconv.FormatUint(mntSeagate.Total/1024/1024/1024, 10), "GB"),
		RootUsed: fmt.Sprint(
			strconv.FormatUint(rootPart.Used/1024/1024/1024, 10), "GB"),
		RootTotal: fmt.Sprint(strconv.FormatUint(rootPart.Total/1024/1024/1024, 10), "GB"),
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
