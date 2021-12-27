package operation

import (
	"client/heartbeat"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func post(objectName, objectPath string, size int64, isBigObject bool) (server, token string, e error) {
	apiServer := heartbeat.ChooseRandomDataServer()
	if apiServer == "" {
		return "", "", fmt.Errorf("get api server is bad, number of api server is 0")
	}
	url := fmt.Sprintf("http://%s/objects/%s", apiServer, objectName)
	logger.Println(url)
	// 计算sha256值
	objectDataReader, err := os.Open(objectPath)
	if err != nil   {
		return "", "", fmt.Errorf("open %s is bad, err is %v", objectPath, err)
	}
	defer objectDataReader.Close()
	hash, err := getBase64(isBigObject, objectDataReader)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest("POST", url, objectDataReader)
	if err != nil {
		logger.Println(err)
		return apiServer, "", err
	}

	// http.head添加配置
	req.Header.Set("size", fmt.Sprintf("%d", size))
	req.Header.Add("Digest", hash)
	req.Header.Add("isBigObjbect", fmt.Sprintf("%b", isBigObject))
	client := http.Client{}
	reps, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return apiServer, "", err
	}
	//logger.Println("post request run time is", time.Since(bTime))
	defer reps.Body.Close()
	if reps.StatusCode != 201 {
		if reps.StatusCode == http.StatusOK {
			// 说明object已经存在，不用head也不用上传了
			buf, _ := io.ReadAll(reps.Body)
			return apiServer, "", errors.New(string(buf))
		}
		// 说明post请求报错了
		return apiServer, "", errors.New(fmt.Sprintf("http code is not 201, code is %d", reps.StatusCode))
	}
	return apiServer, reps.Header.Get("location"), nil
}
func getBase64(isBigOjbect bool, reader io.Reader) (encode string, e error) {
	if isBigOjbect {
		logger.Println("计算大对象base编码")
		return getBase64OfBigOjbect(reader)
	} else {
		logger.Println("计算小对象base编码")
		return getBase64OfMinOjbect(reader)
	}
}

func getBase64OfMinOjbect(reader io.Reader) (encode string, e error)  {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.New(fmt.Sprintf("read data is bad, err is %v", err))
	}
	hash := sha256.New()
	_, err = hash.Write(buf)
	if err != nil {
		return "", errors.New("sha256 is bad")
	}
	// base64编码
	return url.PathEscape(fmt.Sprintf("SHA-256=%s",base64.StdEncoding.EncodeToString(hash.Sum(nil)))), nil
}
func getBase64OfBigOjbect(reader io.Reader) (encode string, e error) {
	hash := sha256.New()
	chunckSize := 1024 * 1024 * 5
	aaa := 0
	for {
		buf := make([]byte, chunckSize)
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return "", errors.New(fmt.Sprintf("read data is bad, err is %v", err))
			}
		}

		_, err = hash.Write(buf[:n])
		if err != nil {
			return "", errors.New("sha256 is bad")
		}
		aaa+=n

		//logger.Println("read data is objectDataReaderk, size is", aaa, "n is", n)
		if n == 0 {
			break
		}
	}
	// base64编码
	return url.PathEscape(fmt.Sprintf("SHA-256=%s", base64.StdEncoding.EncodeToString(hash.Sum(nil)))), nil
}