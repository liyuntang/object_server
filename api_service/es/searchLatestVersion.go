package es

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func SearchLatestVersion(name string) (int64, error){
	engine := getEngine()
	if engine == nil {
		return 0, fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return 0, fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()

	metaData := &Object_server_metadata{}
	metaData.Name = name
	tx := engine.Where("name = ?", name).Last(metaData)
	//fmt.Println(metaData, tx.Error)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.Error == nil{
		return metaData.Version, nil
	}

	return 0, tx.Error
}

func CheckObjectISExsit(name string) (bool, int64, error) {
	engine := getEngine()
	if engine == nil {
		return false, -1, fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return false, -1, fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()

	metaData := &Object_server_metadata{}
	metaData.Name = name
	tx := engine.Where("name = ?", name).Last(metaData)
	//fmt.Println(metaData, tx.Error)
	if tx.Error != nil && errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		// 说明object不存在
		return false, -1, tx.Error
	}
	// 说明object存在
	return true, metaData.Version, nil

}


