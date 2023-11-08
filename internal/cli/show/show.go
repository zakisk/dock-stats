package show

import (
	"context"
	"encoding/json"

	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli/command/container"
	"github.com/docker/docker/api/types"
	sdkClient "github.com/docker/docker/client"
	"github.com/pterm/pterm"
	"github.com/zakisk/dock-stats/pkg/logger"
	"github.com/zakisk/dock-stats/pkg/utils"
	"github.com/zakisk/drawille-go"
)

type Attribute int

const escape = "\x1b"

// rootCmd represents the base command when called without any subcommands
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows containers stats",
	Long: `Command to render container real-time stats 
Usage: dock-stats show [CONTAINER]

Example: dock-stats show ubuntu
`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLog()
		if len(args) == 0 {
			log.Fatal("Please specify container name", log.Args())
		}

		confirm := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Before showing stats, It will clear your terminal, Do you agree?")
		result, _ := confirm.Show()
		if !result {
			return
		}

		client, err := sdkClient.NewClientWithOpts(sdkClient.FromEnv, sdkClient.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal("Error while creating docker client", log.Args("error", err))
		}
		statsResponse, err := client.ContainerStats(context.Background(), args[0], true)
		if err != nil {
			log.Fatal("Error while getting container stats", log.Args("error", err))
		}
		defer statsResponse.Body.Close()

		statChan, errChan := getStats(log, statsResponse)

		canvas := drawille.NewCanvas(103, 30)
		canvas.LineColors = []drawille.Color{
			drawille.DarkCyan, // cpu
			drawille.Yellow,   // memory
		}
		canvas.LabelColor = drawille.Crimson
		canvas.AxisColor = drawille.DarkOrchid
		canvas.NumDataPoints = 50
		canvas.HorizontalLabels = []string{}
		formatter := NewFormatter(statsResponse.OSType)

		for {
			select {
			case stats := <-statChan:
				{
					formatter.setStats(stats)
					formatter.Plot(canvas)
				}
			case err := <-errChan:
				{
					if err != nil {
						log.Error(err.Error(), log.Args())
						if err == io.EOF {
							break
						}
					}
				}
			}
		}
	},
}

func getStats(log *logger.Log, statsResponse types.ContainerStats) (chan *StatsItem, chan error) {
	statChan := make(chan *StatsItem)
	errChan := make(chan error)
	var (
		previousCPU    uint64
		previousSystem uint64
	)

	dec := json.NewDecoder(statsResponse.Body)
	go func() {
		for {
			var (
				stats                  *types.StatsJSON
				memPercent, cpuPercent float64
				blkRead, blkWrite      uint64
				mem, memLimit          float64
				pidsStatsCurrent       uint64
			)
			if err := dec.Decode(&stats); err != nil {
				dec = json.NewDecoder(io.MultiReader(dec.Buffered(), statsResponse.Body))
				errChan <- err
				if err == io.EOF {
					statsResponse.Body.Close()
					close(statChan)
					close(errChan)
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}

			daemonOSType := statsResponse.OSType

			// all these values are evaluated using copied code from docker's source
			// values showed in this app's output are exactly same as docker gives in `docker stats` command.
			if daemonOSType != "windows" {
				previousCPU = stats.PreCPUStats.CPUUsage.TotalUsage
				previousSystem = stats.PreCPUStats.SystemUsage
				cpuPercent = utils.CalculateCPUPercentUnix(previousCPU, previousSystem, stats)
				blkRead, blkWrite = utils.CalculateBlockIO(stats.BlkioStats)
				mem = utils.CalculateMemUsageUnixNoCache(stats.MemoryStats)
				memLimit = float64(stats.MemoryStats.Limit)
				memPercent = utils.CalculateMemPercentUnixNoCache(memLimit, mem)
				pidsStatsCurrent = stats.PidsStats.Current
			} else {
				cpuPercent = utils.CalculateCPUPercentWindows(stats)
				blkRead = stats.StorageStats.ReadSizeBytes
				blkWrite = stats.StorageStats.WriteSizeBytes
				mem = float64(stats.MemoryStats.PrivateWorkingSet)
			}
			netRx, netTx := utils.CalculateNetwork(stats.Networks)

			statsItem := &StatsItem{
				StatsEntry: container.StatsEntry{
					Name:             stats.Name,
					ID:               stats.ID,
					CPUPercentage:    cpuPercent,
					Memory:           mem,
					MemoryPercentage: memPercent,
					MemoryLimit:      memLimit,
					NetworkRx:        netRx,
					NetworkTx:        netTx,
					BlockRead:        float64(blkRead),
					BlockWrite:       float64(blkWrite),
					PidsCurrent:      pidsStatsCurrent,
				},
				daemonOSType: daemonOSType,
			}
			statChan <- statsItem
		}
	}()
	return statChan, errChan
}
