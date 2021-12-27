package heartbeat

import (
	"data_server/common"
	"data_server/rabbitmq"
	"flag"
	"github.com/shirou/gopsutil/disk"
	"log"
	"os"
	"time"
)
var RABBITMQ_SERVER string
var LISTEN_ADDRESS string
var device string
var logger *log.Logger
func init() {
	logger = common.WriteLog("logs/log.file")
	flag.StringVar(&RABBITMQ_SERVER, "r", "amqp://admin:123.com@10.10.30.207:5672", "rabbit server")
	flag.StringVar(&LISTEN_ADDRESS, "listen", "", "api server")
	flag.StringVar(&device, "device", "/data", "数据盘")
}

type DataServerInfo struct {
	HostName string
	Address string
	UsageStat *disk.UsageStat
}
func StartHeartbeat() {
	q := rabbitmq.New(RABBITMQ_SERVER)
	defer q.Close()
	for {
		status := getDiskStat()
		q.Publish("apiServers", status)
		time.Sleep(5*time.Second)
	}
}

func getDiskStat() *DataServerInfo {
	s, err := disk.Usage(device)
	if err != nil {
		panic(err)
	}

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	info := &DataServerInfo{}
	info.UsageStat = s
	info.HostName = name
	info.Address = LISTEN_ADDRESS
	return info
}
