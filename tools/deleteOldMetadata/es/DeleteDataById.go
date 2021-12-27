package es

import "fmt"

func DeleteDataById(ids []uint) error{
	engine := getEngine()
	if engine == nil {
		return fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return  fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()
	meta := &Object_server_metadata{}

	tx := engine.Delete(&meta, ids)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
