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
	MntUsed    string `json:"home-used"`
	MntTotal   string `json:"home-total"`
	RootUsed   string `json:"root-used"`
	RootTotal  string `json:"root-total"`
}

func NewSysInfo() *SysInfo {
	m, err := mem.VirtualMemory()
	checkErr(err)

	parts, err := disk.Partitions(false)
	checkErr(err)

	var rootPart *disk.UsageStat
	var mntParts *disk.UsageStat
	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)
		checkErr(err)

		fmt.Println(u.Path)

		if u.Path == "/" {
			rootPart = u
		} else if strings.HasPrefix(u.Path, "/mnt") {
			mntParts = u
		}
	}

	return &SysInfo{
		MemPercent: fmt.Sprint(strconv.FormatFloat(m.UsedPercent, 'f', 2, 64), "%"),
		MemUsed:    fmt.Sprint(strconv.FormatUint(m.Used/1024/1024, 10), "mb"),
		MemTotal:   fmt.Sprint(strconv.FormatUint(m.Total/1024/1024, 10), "mb"),
		MntUsed: fmt.Sprint(
			strconv.FormatUint(mntParts.Used/1024/1024/1024, 10), "GB"),
		MntTotal: fmt.Sprint(strconv.FormatUint(mntParts.Total/1024/1024/1024, 10), "GB"),
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
