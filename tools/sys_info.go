package tools

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type SysInfo struct {
	MemPercent string `json:"mem-perc"`
	MemUsed    string `json:"mem-used"`
	MemTotal   string `json:"mem-total"`
	HomeUsed   string `json:"home-used"`
	HomeTotal  string `json:"home-total"`
	RootUsed   string `json:"root-used"`
	RootTotal  string `json:"root-total"`
}

func NewSysInfo() *SysInfo {
	m, err := mem.VirtualMemory()
	checkErr(err)

	parts, err := disk.Partitions(false)
	checkErr(err)

	var homePart *disk.UsageStat
	var rootPart *disk.UsageStat
	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		checkErr(err)

		if u.Path == "/" {
			rootPart = u
		} else if strings.HasPrefix(u.Path, "/home") {
			homePart = u
		}
	}

	return &SysInfo{
		MemPercent: fmt.Sprint(strconv.FormatFloat(m.UsedPercent, 'f', 2, 64), "%"),
		MemUsed:    fmt.Sprint(strconv.FormatUint(m.Used/1024/1024, 10), "mb"),
		MemTotal:   fmt.Sprint(strconv.FormatUint(m.Total/1024/1024, 10), "mb"),
		HomeUsed: fmt.Sprint(
			strconv.FormatUint(homePart.Used/1024/1024/1024, 10), "GB"),
		HomeTotal: fmt.Sprint(strconv.FormatUint(homePart.Total/1024/1024/1024, 10), "GB"),
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
