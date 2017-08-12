package models

import (
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"log"
	"strconv"

	"github.com/jetlwx/comm"
)

type Containers struct {
	ID            int64  `xorm:"notnull autoincr pk"`
	ContainerName string `xorm:"varchar(100) "`
}
type Dockermonitor struct {
	ID             int64     `xorm:"notnull autoincr pk"`
	TimeStamp      time.Time `xorm:"notnull"`
	ContainerName  string    `xorm:"varchar(100) "`
	PodName        string    `xorm:"varchar(100)"`
	CpuPercent     float64   `xorm:"default 0"`
	MemPercent     float64   `xorm:"default 0"`
	BlkReadToMB    float64   `xorm:"default 0"`
	BlkWriteToMB   float64   `xorm:"default 0"`
	TxToMb         float64   `xorm:"default 0"`
	RxToMb         float64   `xorm:"default 0"`
	TxDropPercent  float64   `xorm:"default 0"`
	TxErrorPercent float64   `xorm:"default 0"`
	RxDropPercent  float64   `xorm:"default 0"`
	RxErrorPercent float64   `xorm:"default 0"`
}

//-----------------------------------
var engine *xorm.Engine

//初始化数据库
func DBinit(dbuser, dbpass, dbhost, dbname string, dbport int) {
	var err error

	dbsource := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + strconv.Itoa(dbport) + ")/" + dbname + "?charset=utf8"

	engine, err = xorm.NewEngine("mysql", dbsource)
	if err != nil {

		//log.Fatalf("[ E ] 创建数据库引擎失败 %v", err)
		log.Fatalln("E", "创建数据库引擎失败 %v", comm.CustomerErr(err))

	}

	//是否显示打印SQL信息
	b, _ := strconv.ParseBool(beego.AppConfig.String("showsqllog"))
	if b {
		engine.ShowSQL(b)
	}
	//数据库PING测试
	if engine.Ping() == nil {
		log.Println("[ D ]数据库" + dbhost + ":" + strconv.Itoa(dbport) + "连接成功！！！")

	}

	//同步至数据库
	//注意：增加表后，要增加new(表名)
	err = engine.Sync2(new(Containers),
		new(Dockermonitor))
	if err != nil {

		log.Fatalln("[ E ]同步数据结构至数据库失败，可能原因为：%v ", err)

	}

}
