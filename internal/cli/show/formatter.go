package show

import (
	"fmt"
	"strings"

	"github.com/zakisk/dock-stats/pkg/utils"
	"github.com/zakisk/drawille-go"
)

var (
	winOSType = "windows"
)

type formatter struct {
	id                string
	name              string
	data              [][]float64
	daemonOSType      string
	netIO             string
	blockIO           string
	memoryUsage       string
	pids              uint64
	privateWorkingSet string
}

func NewFormatter(daemonOSType string) *formatter {
	var data [][]float64
	if daemonOSType != winOSType {
		data = make([][]float64, 2)
	} else {
		data = make([][]float64, 1)
	}
	return &formatter{data: data, daemonOSType: daemonOSType}
}

func (f *formatter) Plot(canvas drawille.Canvas) {
	utils.ClearScreen()
	containerId := drawille.ColorString(fmt.Sprintf("CONTAINER ID: %s", f.id), drawille.Aqua)
	name := drawille.ColorString(fmt.Sprintf("NAME: %s\n\n", f.name), drawille.DarkSlateBlue)
	space := strings.Repeat(" ", 6)
	net := drawille.ColorString(fmt.Sprintf("NETWORK I/O (Received/Transferred): %s%s", f.netIO, space), drawille.Tomato)
	block := drawille.ColorString(fmt.Sprintf("BLOCK I/O (Read/Write): %s\n\n", f.blockIO), drawille.Chartreuse)
	diff := len(net) - len(containerId)
	space = strings.Repeat(" ", diff - 1)
	fmt.Printf("%s%s", containerId, space)
	fmt.Print(name, net, block)
	if f.daemonOSType != winOSType {
		processes := fmt.Sprintf("PROCESSES (PIDs): %d\n\n", f.pids)
		fmt.Print(drawille.ColorString(processes, drawille.DeepPink))
	} else {
		privateWorkingSet := fmt.Sprintf("PRIVATE WORKING SET: %s\n\n", f.privateWorkingSet)
		fmt.Print(drawille.ColorString(privateWorkingSet, drawille.DeepPink))
	}
	// print line
	fmt.Printf("%s\n\n", drawille.ColorString(utils.GetLine(), drawille.Goldenrod))

	// CPU and Memory (OS != Windows) Stats
	fmt.Printf("%s\n\n", canvas.PlotWithMinAndMax(f.data, 0, 100, true))
	cpuUsage := fmt.Sprintf("%.2f", f.data[0][len(f.data[0])-1]) + "%\t\t"
	fmt.Print(drawille.ColorString("██ CPU USAGE(%): "+cpuUsage, drawille.DarkCyan))
	if f.daemonOSType != winOSType {
		memUsage := fmt.Sprintf("%.2f", f.data[1][len(f.data[1])-1]) + "% " + fmt.Sprintf("(%s)\n\n", f.memoryUsage)
		memString := drawille.ColorString("██ Memory USAGE(%): "+memUsage, drawille.Yellow)
		fmt.Print(memString)
	}
}

func (f *formatter) setStats(item *StatsItem) {
	f.id = item.ID[:12]
	f.name = item.Name
	f.data[0] = append(f.data[0], item.CPUPercentage)
	if f.daemonOSType != winOSType {
		f.data[1] = append(f.data[1], item.MemoryPercentage)
	}
	f.netIO = item.NetIO()
	f.blockIO = item.BlockIO()
	f.memoryUsage = item.MemUsage(f.daemonOSType)
	f.pids = item.PidsCurrent
	f.privateWorkingSet = item.MemUsage(f.daemonOSType)
}
