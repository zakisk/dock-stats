package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/fatih/color"
	"golang.org/x/term"
)

func IsDockerInstalled() error {
	return exec.Command("docker", "version").Run()
}

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// code copied from docker's source code.
func CalculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
		onlineCPUs  = float64(v.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func CalculateCPUPercentWindows(v *types.StatsJSON) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                         // Convert to number of 100ns intervals
	possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors

	// Intervals used
	intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals) * 100.0
	}
	return 0.00
}

func CalculateBlockIO(blkio types.BlkioStats) (uint64, uint64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return blkRead, blkWrite
}

func CalculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}

// calculateMemUsageUnixNoCache calculate memory usage of the container.
// Cache is intentionally excluded to avoid misinterpretation of the output.
//
// On cgroup v1 host, the result is `mem.Usage - mem.Stats["total_inactive_file"]` .
// On cgroup v2 host, the result is `mem.Usage - mem.Stats["inactive_file"] `.
//
// This definition is consistent with cadvisor and containerd/CRI.
// * https://github.com/google/cadvisor/commit/307d1b1cb320fef66fab02db749f07a459245451
// * https://github.com/containerd/cri/commit/6b8846cdf8b8c98c1d965313d66bc8489166059a
//
// On Docker 19.03 and older, the result was `mem.Usage - mem.Stats["cache"]`.
// See https://github.com/moby/moby/issues/40727 for the background.
func CalculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	// cgroup v1
	if v, isCgroup1 := mem.Stats["total_inactive_file"]; isCgroup1 && v < mem.Usage {
		return float64(mem.Usage - v)
	}
	// cgroup v2
	if v := mem.Stats["inactive_file"]; v < mem.Usage {
		return float64(mem.Usage - v)
	}
	return float64(mem.Usage)
}

func CalculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

func Box() {
	// Function to print colored text
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// Define the size of the box
	width, height := 10, 5

	// Print the top of the box
	fmt.Print(red("┌"))
	for i := 0; i < width-2; i++ {
		fmt.Print(red("─"))
	}
	fmt.Println(red("┐"))

	// Print the middle of the box
	for i := 0; i < height-2; i++ {
		fmt.Print(red("│"))
		for j := 0; j < width-2; j++ {
			fmt.Print(blue("█")) // Use any color to fill the inside of the box
		}
		fmt.Println(red("│"))
	}

	// Print the bottom of the box
	fmt.Print(red("└"))
	for i := 0; i < width-2; i++ {
		fmt.Print(red("─"))
	}
	fmt.Println(red("┘"))
}

func GetLine() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return ""
	}

	dashes := ""
	for i := 0; i < width; i++ {
		dashes += "-"
	}

	return dashes
}
