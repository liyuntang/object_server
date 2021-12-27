package operation

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func put(url string, offset int64, reader io.Reader) error {
	/* 上传数据，这里用的是put方法，有个地方需要注意，put请求头部要加上range信息，格式是range:bytes=XXX-
	 */

	//logger.Println("url is", url)
	//logger.Println("offset is", offset)
	request, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		logger.Println("new request of put is bad, err is", err)
		return err
	}
	request.Header.Set("range", fmt.Sprintf("bytes=%d-", offset))
	client := http.Client{}
	resp, err := client.Do(request)
	if resp.StatusCode == http.StatusOK {
		//logger.Println("status code is", resp.StatusCode)
		resp.Body.Close()
		logger.Println("数据上传成功")
		return nil
	}
	//logger.Println("数据上传失败, err is", err)
	return errors.New(fmt.Sprintf("数据上传失败, err is %v, code is %d", err, resp.StatusCode))
}

func head(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		logger.Println("head operation is bad, err is", err)
		return "", errors.New("head operation is bad")
	}
	return resp.Header.Get("content-length"), nil
}
