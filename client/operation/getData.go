package operation

import (
	"fmt"
	"os"
	"time"
)

func createFile(fileName string, data *string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s%s%s", fileName, time.Now().UnixNano(), *data))
	if err != nil {
		return err
	}
	return nil
}

func getData(lenth int64) *string {
	text := "aaa"
	neirong := time.Now().UnixNano()*2
	for int64(len(text)) < lenth {
		text = text + fmt.Sprintf("%d", neirong)
	}
	return &text
}