package main

import (
	"deleteOldMetadata/es"
	"fmt"
)

const MIN_VERSION_COUNT = 5

func main() {
	// 查询所有版本超过MIN_VERSION_COUNT数量的object
	objectNameSlice, err := es.SearchVersionStatus(MIN_VERSION_COUNT)
	if err != nil {
		fmt.Println("get meta data is bad, err is", err)
		return
	}
	// 判断是否有需要删除的数据
	if len(objectNameSlice) == 0 {
		// 说明没有需要删除的object
		fmt.Println("there no object need to delete")
		return
	}

	// 根据返回的objectNameSlice查询该objectName对应的版本号
	idSlice, err := es.SearchIdOfObject(objectNameSlice, MIN_VERSION_COUNT)
	if err != nil {
		fmt.Println("search id of object is bad, err is", err)
		return
	}
	// 清理数据
	if err = es.DeleteDataById(idSlice); err != nil {
		fmt.Println("delete meta data is bad, err is", err)
		return
	}
	fmt.Println("delete meta data is ok")



}
