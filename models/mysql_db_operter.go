package models

import (
	"log"
	"time"
)

func InsertContaier(c string) {
	con := new(Containers)
	total, er := engine.Where("container_name=?", c).Count(con)
	if er != nil {
		log.Println("an error at check the recorder,", er)
		return
	}
	if total == 0 {
		cc := new(Containers)
		cc.ContainerName = c
		_, err := engine.Insert(cc)
		if err != nil {
			log.Println("an error at insert container names:", err)
			return
		}
	}

	return
}

func InsertRecorder(r Res) {
	r1 := new(Dockermonitor)
	r1.TimeStamp = time.Now()
	r1.BlkReadToMB = r.BlkReadToMB
	r1.BlkWriteToMB = r.BlkWriteToMB
	r1.ContainerName = r.ContainerName
	r1.CpuPercent = r.CpuPercent
	r1.MemPercent = r.MemPercent
	r1.PodName = r.PodName
	r1.RxDropPercent = r.RxDropPercent
	r1.RxErrorPercent = r.RxErrorPercent
	r1.RxToMb = r.RxToMb
	r1.TxDropPercent = r.TxDropPercent
	r1.TxErrorPercent = r.TxErrorPercent
	r1.TxToMb = r.TxToMb
	_, err := engine.Insert(r1)
	if err != nil {
		log.Println("an error at insert recorders:", err)
		return
	}

	return
}
