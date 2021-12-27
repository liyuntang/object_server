package es

import "fmt"

// create
func PutMetadata(name string, version int64, size int64, hash string) error{
	engine := getEngine()
	if engine == nil {
		return fmt.Errorf("connect to mysql is bad")
	}
	db, err:= engine.DB()
	if err != nil {
		return fmt.Errorf("connect to mysql is bad")
	}
	defer db.Close()

	metaData := &Object_server_metadata{}
	metaData.Name = name
	metaData.Version = version
	metaData.Size = size
	metaData.Hash = hash
	tx := engine.Create(metaData)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

