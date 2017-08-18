package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/jetlwx/dockerMonitor/models"
)

var (
	dockerAPIurl string
	interval     int
	dbdriver     string
	mysqluser    string
	mysqlpass    string
	mysqlhost    string
	mysqlport    int
	mysqldb      string
	showlog      bool
	threads      int

	influxurl  string
	influxuser string
	influxpass string
	influxdb   string
	cli        client.Client
	bp         client.BatchPoints
)

func Init() {
	flag.StringVar(&dbdriver, "db", "influxdb", "backend database,mysql or influxdb,default mysql")
	flag.StringVar(&mysqluser, "myuser", "docker", "mysql db user,default docker")
	flag.StringVar(&mysqlpass, "mypass", "docker", "mysql db user's pass ,default docker")
	flag.StringVar(&mysqlhost, "myhost", "172.16.6.156", "mysql db host,default 127.0.0.1")
	flag.IntVar(&mysqlport, "myport", 3306, "mysql db host port ,default 3306")
	flag.StringVar(&mysqldb, "mydb", "docker", "mysqldb name ,default docker")
	flag.StringVar(&dockerAPIurl, "dockerAPIurl", "http://172.16.16.2:8888", "docker remote API ,default http://127.0.0.1:8888, eg: /usr/bin/dockerd   -H 127.0.0.1:8888 -H unix:///var/run/docker.sock")
	flag.IntVar(&interval, "interval", 120, "the frequency(seconds) of collect data ,default 120")
	flag.BoolVar(&showlog, "log", true, "showlog ?,default false")
	flag.IntVar(&threads, "threads", 2, "how many thread to collect,default 2")

	flag.StringVar(&influxurl, "influxurl", "http://172.16.18.2:8086", "influxdb API address")
	flag.StringVar(&influxdb, "influxdb", "docker", "infulx database name ,default docker")
	flag.StringVar(&influxuser, "influxuser", "docker", "influx login username ,default docker")
	flag.StringVar(&influxpass, "influxpassword", "docker", "influxdb login password,default docker")
	fmt.Println("Version: 20170811")
}

func main() {
	Init()
	flag.Parse()

	switch dbdriver {
	case "mysql":
		models.DBinit(mysqluser, mysqlpass, mysqlhost, mysqldb, mysqlport)
	case "influxdb":
		cli, bp = models.InitInfluxDBClient(influxurl, influxdb, influxuser, influxpass, 5)
	}

	interval64 := int64(interval)
	runtime.GOMAXPROCS(threads)
	wg := sync.WaitGroup{}
	count := int64(0)

	for {
		if count%interval64 == 0 {
			l := models.DockerList(dockerAPIurl)
			wg.Add(len(l))
			//	fmt.Println("l----->", len(l))

			for _, v := range l {
				go models.GetDockinfo(dockerAPIurl, v.ID, v.ContainerName, v.PodName, &wg, showlog, cli, bp)
			}
			wg.Wait()
			models.WriteDB(cli, bp)
		}

		time.Sleep(1 * time.Second)
		count++
	}
}
