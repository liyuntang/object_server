package es

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func SearchVersionStatus(maxVersion int64) ([]string, error){
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
	result := []*Result{}

	tx := engine.Model(metaData).Select("name, count(version) as total").Group("name").Find(&result)
	if tx.Error != nil && errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		// 说明object不存在
		return nil, tx.Error
	}
	// 说明object存在
	objectSlice := []string{}
	for _, data := range result {
		if data.Total > maxVersion {
			objectSlice = append(objectSlice, data.Name)
		}

	}
	return objectSlice, nil
}

func SearchIdOfObject(objectSlice []string, MIN_VERSION_COUNT int) (idSlice []uint, e error){
	engine := getEngine()
	if engine == nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return nil, fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()
	idList := []uint{}
	for _, objectName := range objectSlice {
		//metaData := &Object_server_metadata{}
		result := []Object_server_metadata{}

		tx := engine.Where("name = ?", objectName).Find(&result)
		if tx.Error != nil && errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// 说明object不存在
			return nil, tx.Error
		}
		length := len(result)
		dLength := length - MIN_VERSION_COUNT
		for i:=0; i< dLength; i++ {
			idList = append(idList, result[i].ID)
		}
	}
	return idList, nil
}

