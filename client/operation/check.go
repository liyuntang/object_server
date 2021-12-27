package operation

import "os"

func check() {
	if objectSize <= 0 || threads <= 0{
		logger.Println("sorry, object size is 0 or threads is 0")
		os.Exit(0)
	}
	_, err := os.Stat(tmpDir)
	if err != nil && os.IsNotExist(err) {
		logger.Println(tmpDir, "is not exist, create it")
		os.MkdirAll(tmpDir, 0755)
	}
}
