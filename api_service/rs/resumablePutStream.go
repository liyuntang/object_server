package rs

import (
	"api_service/objectstream"
	"encoding/base64"
	"encoding/json"
)

type resumableToken struct {
	Name string
	Size int64
	Hash string
	Servers []string
	Uuids []string
}

type RSResumablePutStream struct {
	*RSputStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, err := NewRSputStream(dataServers, hash, size)
	if err != nil {
		//logger.Println("err................")
		return nil, err
	}
	//logger.Println("NewRSputStream is ok...........")
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	//logger.Println("api make token", uuids)
	//logger.Println("object name is", name)
	//logger.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	token := &resumableToken{Name: name, Size: size,Hash: hash, Servers: dataServers, Uuids: uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

func (s *RSResumablePutStream) ToToken() string{
	buf, _ := json.Marshal(s)
	//logger.Println(">>>>>>>>>", string(buf))
	token := base64.StdEncoding.EncodeToString(buf)
	//logger.Println("?????????????", token)
	return token
}









