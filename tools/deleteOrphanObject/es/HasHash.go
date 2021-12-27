package es

import (
	"flag"
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var thread int

func init() {
	flag.IntVar(&thread, "t", 50, "threads")
}
func HasHash(hashSlice []string) ([]string, error) {
	engine := getEngine()
	if engine == nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	db, err := engine.DB()
	if err != nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()

	list := []string{}
	taskChannl := make(chan string)
	resultChannl := make(chan string)

	// 处理数据
	for i:=1;i<=thread;i++ {
		go func(id int) {
			for hash := range taskChannl {
				defer wg.Done()
				metaData := &Object_server_metadata{}
				//fmt.Println(id, ">>>>>>>>>>", hash)
				rows := engine.Select("id").Where("hash = ?", hash).Find(metaData).RowsAffected
				if rows == 0 {
					// 说明该hash不存在
					resultChannl <- hash
				}
			}
		}(i)
	}

	// 接收数据
	go func() {
		for hash := range resultChannl {
			//fmt.Println("出数据了", hash)
			list = append(list, hash)
		}
	}()

	// 分发数据
	for _, hash := range hashSlice {
		wg.Add(1)
		taskChannl <- hash
	}
	close(taskChannl)
	wg.Wait()
	//close(resultChannl)

	fmt.Println("over..................")
	return list, nil
}
