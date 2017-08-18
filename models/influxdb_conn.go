package models

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type Tags struct {
	Name  string
	Value string
}
type Fileds struct {
	Name  string
	Value interface{}
}

/*
config set:
     	InfluxConfig.Addr = influxurl
		InfluxConfig.Username = influxuser
		InfluxConfig.Password = influxpass
		InfluxConfig.Timeout = 5 * time.Second
*/
// func InitInfluxClient(conf client.HTTPConfig) (cc client.Client) {
// 	c, err := client.NewHTTPClient(conf)
// 	if err != nil {
// 		log.Fatalln("Error creating InfluxDB Client: ", err.Error())
// 	}
// 	return c
// }

func InitInfluxDBClient(httpAPI, dbname, username, password string, httptimeout int) (client.Client, client.BatchPoints) {
	c1 := client.HTTPConfig{}
	b1 := client.BatchPointsConfig{}
	c1.Addr = httpAPI
	c1.Password = password
	//c1.Timeout = time.Duration(httptimeout)
	c1.Username = username
	b1.Database = dbname

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(c1)
	if err != nil {
		log.Fatal("Create http client error:", err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(b1)
	if err != nil {
		log.Fatal("Create New Batch Points Error:", err)
	}

	return c, bp
}

func gentTag(t []Tags) (n map[string]string) {
	m := make(map[string]string)
	for _, v := range t {
		fmt.Println("v.name=", v.Name, "v.vale=", v.Value)
		m[v.Name] = v.Value
	}
	return m
}

func gentField(f []Fileds) (n map[string]interface{}) {
	m := make(map[string]interface{})
	for _, v := range f {
		m[v.Name] = v.Value
	}
	return m
}

func AddInfluxDBPoint(c client.Client, bp client.BatchPoints, measureName string, tag []Tags, field []Fileds, t time.Time) {
	// Create a point and add to batch
	fmt.Println("tag=", tag)
	fmt.Println("field=", field)
	tags := gentTag(tag)
	fileds := gentField(field)
	pt, err := client.NewPoint(measureName, tags, fileds, t)
	if err != nil {
		log.Println("Create New Point error:", err)
	}

	bp.AddPoint(pt)
	return
}

func WriteDB(c client.Client, bp client.BatchPoints) {
	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Println("Write the batch error:", err)
	}

	err := c.Close()
	if err != nil {
		log.Println("[E] error at close conn", err)
	} else {
		log.Println("close ok")
	}

	return
}
