package rs

import (
	"api_service/common"
	"api_service/objectstream"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"log"
)

const (
	DATA_SHARDS = 2
	PARITY_SHARDS = 1
	ALL_SHARDS = DATA_SHARDS + PARITY_SHARDS
	BLOCK_PER_SHARD = 1024*1024*3
	BLOCK_SIZE = BLOCK_PER_SHARD * DATA_SHARDS
)

var logger *log.Logger
func init() {
	logger = common.WriteLog("logs/log.file")
}

type RSputStream struct {
	*encoder
}

func (s *RSputStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
func NewRSputStream(dataServers []string, hash string, size int64) (*RSputStream, error) {
	//logger.Println("data servers is", dataServers)
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismath")
	}
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		/*
		url中的转译问题：
			1、包含/或者%2F可以正常请求
			2、包含//或者%2F%2F不可以正常请求
		 */
		name := fmt.Sprintf("%s.%d", hash, i)
		server := dataServers[i]
		logger.Println("i is", i, "name is", name, "server is", server)
		writers[i], e = objectstream.NewTempPutStream(server, name, perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	return &RSputStream{enc}, nil
}
type encoder struct {
	writers []io.Writer
	enc reedsolomon.Encoder
	cache []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

func (e *encoder)Write(p []byte) (n int, err error){
	length := len(p)

	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		//logger.Println("length is", length, "len of cache is", len(e.cache), "block_size is", BLOCK_SIZE)
		// length is 33 len of cache is 0 block_size is 32000
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i:= range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}

