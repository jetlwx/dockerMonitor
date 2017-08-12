package models

import (
	"strings"

	dockerType "github.com/docker/docker/api/types"
	"github.com/jetlwx/comm"
)

//calculate  cpu percent
func CalculateCPUPercentUnix(previousCPU, previousSystem uint64, v dockerType.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return comm.MathRounding(cpuPercent, 2)
}

//calculate the memory useage
func CaclaateMemoryPercent(useage, limit uint64, v dockerType.StatsJSON) float64 {
	memPercent := comm.MathRounding(float64(useage)/float64(limit)*100, 2)
	return memPercent

}

//calculate BlockIO
func CalculateBlockIO(blkio dockerType.StatsJSON) (blkRead uint64, blkWrite uint64) {
	for _, bioEntry := range blkio.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return
}

//calculate Network IO
func CalculateNetwork(network map[string]dockerType.NetworkStats) (rx, tx, txDropPercent, txErrorPercent, rxDropPercent, rxErrorPercent float64) {
	var txDrop, rxDrop, txError, rxError float64
	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
		txDrop += float64(v.TxDropped)
		rxDrop += float64(v.RxDropped)
		txError += float64(v.TxErrors)
		rxError += float64(v.RxErrors)
	}
	if tx <= 0 {
		txDropPercent = 0.0
		txErrorPercent = 0.0
	} else {
		txDropPercent = comm.MathRounding(txDrop/tx, 2)
		txErrorPercent = comm.MathRounding(txError/rx, 2)
	}
	if rx <= 0 {
		rxDropPercent = 0.0
		rxErrorPercent = 0.0
	} else {
		rxDropPercent = comm.MathRounding(rxDrop/rx, 2)
		rxErrorPercent = comm.MathRounding(rxDrop/rx, 2)
	}

	return rx, tx, txDropPercent, txErrorPercent, rxDropPercent, rxErrorPercent
}
