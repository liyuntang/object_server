package main

import (
	"client/heartbeat"
	"client/operation"
	"errors"
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var model string
var object string

func init() {
	flag.StringVar(&model, "model", "auto", "运行模式，包含auto、custom，auto: 自动生成测试数据， custom:上传指定文件")
	flag.StringVar(&object, "object", "", "需要上传的对象，请使用绝对路径")
}

func main() {
	flag.Parse()

	//test.Test()
	//return

	GOTIME := "2006-01-02 15:04:05"
	go heartbeat.ListenHeartbeat()
	// 查看是佛获取到了api server地址
	fmt.Println(time.Now().Format(GOTIME), "获取接口服务地址")
	checkApi()
	fmt.Println(time.Now().Format(GOTIME), "获取接口服务地址成功")

	db := getDb()
	sqlDb, _ := db.DB()
	defer sqlDb.Close()
	fmt.Println("connect to mysql is ok")

	// 判断运行模式
	err := checkParameter()
	if err != nil {
		fmt.Println(err)
		return
	}
	if model == "auto" {
		operation.Auto(db)
	} else {
		operation.Custom(object)
	}
}

func getDb() *gorm.DB {
	dsn := "object_server:123.com@tcp(10.10.10.197:9000)/db_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connect to mysql is bad, err is", err)
		return nil
	}
	return db
}

func checkParameter() error {
	if model == "auto" {
		// 自动模式，client自己生成数据,该模式下需要指定thread，object size
	} else if model == "custom" {
		//表示上传指定的数据,此时需要指定对象名称，判断name是否为空
		if len(object) == 0 {
			return errors.New("没有指定object参数")
		}
	} else {
		return errors.New("model设置错误，目前仅支持auto、custom")
	}
	return nil
}
func checkApi() {
	num := 0
	for num <= 0 {
		num = len(heartbeat.ChooseRandomDataServer())
		//fmt.Println("api server num is", num)
	}
}

/*
断点续传流程介绍：

下载：
	客户端在get对象请求时通过设置range头部来告诉接口服务需要从什么位置开始输出对象的数据


数据上传：

	1、接口服务会对数据进行散列值进行校验，当发生网络故障时，如果上传的数据跟期望的不一致，那么整个上传的数据都会被丢弃，所以断点上传在一开始就需要客户端和接口服务
做好约定，使用特定的接口进行上传

	2、客户端在知道自己上传大对象时就主动该用post接口，提供对象的散列值和大小

	3、客户端post对象后会得到一个token，对token进行put可以上传数据，在上传时客户端需要指定range头部来告诉接口服务上传数据的范围，接口服务对token进行解码，
获取6个分片所在的数据服务地址及uuid，分别调用patch将数据写入各个临时对象

	4、客户端每次用put方法访问token之前都需要先用head方法获取当前已上传了多少数据，接口服务对token进行解码，获取6个分片所在的数据服务地址及uuid，仅对第一个分片
调用head获取该分片的当前长度，将这个数据乘以4，作为Content-Length响应头部返回给客户端


*/
