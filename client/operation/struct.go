package operation

type Object_server_markbench struct {
	Thread_id int	`gorm:"thread_id"`
	Host_name string 	`gorm:"host_name"`
	Code int	`gorm:"code"`
	Run_time int64	`gorm:"run_time"`
}

type objectInfo struct {
	name, path, apiServer, token string
	isBigObject bool
	size int64
}

type Object_server_metadata struct {
	ID        uint           `gorm:"primaryKey"`
	Name string
	Version int64
	Size int64
	Hash string
}