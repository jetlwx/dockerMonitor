package models

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jetlwx/comm"
)

type Res struct {
	ContainerName  string
	PodName        string
	CpuPercent     float64
	MemPercent     float64
	BlkReadToMB    float64
	BlkWriteToMB   float64
	TxToMb         float64
	RxToMb         float64
	TxDropPercent  float64
	TxErrorPercent float64
	RxDropPercent  float64
	RxErrorPercent float64
}

type Dlist struct {
	ContainerName string
	PodName       string
	ID            string
}

//get docker list
func DockerList(url string) (D []Dlist) {
	u := url + "/containers/json"
	res, code, err := comm.GetJsonFromUrlLongConn(u)
	if code != 200 || err != nil {
		log.Println("[E] http code=", code, "err=", err)
		return
	}

	li := []types.Container{}
	er := json.Unmarshal(res, &li)
	if er != nil {
		log.Println("[E] an error at unmarshar json:", er)
	}

	for _, v := range li {
		if v.Labels["io.kubernetes.container.name"] == "POD" {
			continue
		}
		d := Dlist{}
		d.ContainerName = v.Labels["io.kubernetes.container.name"]
		d.PodName = v.Labels["io.kubernetes.pod.name"]
		d.ID = v.ID
		D = append(D, d)
	}

	return D

}

//get docker information
///containers/{id}/stats
func GetDockinfo(url, id, containername, podname string, wg *sync.WaitGroup, showlog bool, cli client.Client, bp client.BatchPoints) {
	R := Res{}
	u := url + "/containers/" + id + "/stats?stream=false"
	if showlog {
		fmt.Println("url=", u)
	}
	res, code, err := comm.GetJsonFromUrlLongConn(u)
	if code != 200 || err != nil {
		log.Println("[E] http code=", code, "err=", err)
		return
	}

	r := types.StatsJSON{}
	json.Unmarshal([]byte(res), &r)

	//refence from https://github.com/docker/docker/blob/eb131c5383db8cac633919f82abad86c99bffbe5/cli/command/container/stats_helpers.go#L175-L188
	previousCPU := r.PreCPUStats.CPUUsage.TotalUsage
	previousSystem := r.PreCPUStats.SystemUsage
	cpuPercent := CalculateCPUPercentUnix(previousCPU, previousSystem, r)

	//memory useage percent
	memLimit := r.MemoryStats.Limit
	memUseage := r.MemoryStats.Usage
	memPercent := CaclaateMemoryPercent(memUseage, memLimit, r)

	// blk stats
	blkRead, blkWrite := CalculateBlockIO(r)
	blkReadToMB := float64(blkRead / 1024 / 1024)
	blkWriteToMB := float64(blkWrite / 1024 / 1024)

	//network IO
	rx, tx, txDropPercent, txErrorPercent, rxDropPercent, rxErrorPercent := CalculateNetwork(r.Networks)
	rxToMb := comm.MathRounding(float64(rx/1024/1024), 2)
	txToMb := comm.MathRounding(float64(tx/1024/1024), 2)

	R.ContainerName = containername
	R.PodName = podname
	R.BlkReadToMB = blkReadToMB
	R.BlkWriteToMB = blkWriteToMB
	R.CpuPercent = cpuPercent
	R.MemPercent = memPercent
	R.RxToMb = rxToMb
	R.TxToMb = txToMb
	R.TxDropPercent = txDropPercent
	R.TxErrorPercent = txErrorPercent
	R.RxDropPercent = rxDropPercent
	R.RxErrorPercent = rxErrorPercent
	if showlog {
		log.Printf("%+v", R)
		fmt.Println("")
	}

	if cli == nil && bp == nil { //=nil ,is mysql
		InsertContaier(containername)
		InsertRecorder(R)
	} else {
		t := []Tags{}
		t1 := Tags{}
		t1.Name = "ContainerName"
		t1.Value = containername
		t = append(t, t1)
		t1.Name = "PodName"
		t1.Value = podname
		t = append(t, t1)

		f := []Fileds{}
		f1 := Fileds{}
		f1.Name = "BlkReadToMB"
		f1.Value = blkReadToMB
		f = append(f, f1)

		f1.Name = "BlkWriteToMB"
		f1.Value = blkWriteToMB
		f = append(f, f1)

		f1.Name = "CpuPercent"
		f1.Value = cpuPercent
		f = append(f, f1)

		f1.Name = "MemPercent"
		f1.Value = memPercent
		f = append(f, f1)

		f1.Name = "RxToMb"
		f1.Value = rxToMb
		f = append(f, f1)

		f1.Name = "TxToMb"
		f1.Value = txToMb
		f = append(f, f1)

		f1.Name = "TxDropPercent"
		f1.Value = txDropPercent
		f = append(f, f1)

		f1.Name = "TxErrorPercent"
		f1.Value = txErrorPercent
		f = append(f, f1)

		f1.Name = "RxDropPercent"
		f1.Value = rxDropPercent
		f = append(f, f1)

		f1.Name = "RxErrorPercent"
		f1.Value = rxErrorPercent
		f = append(f, f1)
		AddInfluxDBPoint(cli, bp, "DockerData", t, f, time.Now())
	}

	wg.Done()
	//	fmt.Println("done action", wg)

}
