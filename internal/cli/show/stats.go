package show

import (
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/go-units"
)

type StatsItem struct {
	container.StatsEntry
	daemonOSType string
}

func (s *StatsItem) NetIO() string {
	if s.IsInvalid {
		return "--"
	}
	return units.HumanSizeWithPrecision(s.NetworkRx, 3) + " / " + units.HumanSizeWithPrecision(s.NetworkTx, 3)
}

func (s *StatsItem) BlockIO() string {
	if s.IsInvalid {
		return "--"
	}
	return units.HumanSizeWithPrecision(s.BlockRead, 3) + " / " + units.HumanSizeWithPrecision(s.BlockWrite, 3)
}

func (s *StatsItem) MemUsage(os string) string {
	if s.IsInvalid {
		return "-- / --"
	}
	if os == "windows" {
		return units.BytesSize(s.Memory)
	}
	return units.BytesSize(s.Memory) + " / " + units.BytesSize(s.MemoryLimit)
}
