package es

import "fmt"

func GetMetaData(name string, version int64) (*Object_server_metadata, error) {
	engine := getEngine()
	if engine == nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()
	metaData := &Object_server_metadata{}
	//fmt.Println(">>>>>>>>>>>", name, version)
	tx := engine.Where("name = ? and version = ?", name, version).Find(metaData)

	if tx.Error != nil {
		return nil, err
	}

	return metaData, nil
}
