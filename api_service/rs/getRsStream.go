package rs

import (
	"api_service/objectstream"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
)

type RsGetStream struct {
	*decoder
}

func NewRsGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RsGetStream, error) {
	if len(locateInfo) + len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	//logger.Println(">>>>>>>>>>>>>", len(locateInfo), len(dataServers))
	readers := make([]io.Reader, ALL_SHARDS)
	for i:=0;i< ALL_SHARDS;i++ {
		server := locateInfo[i]
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil{
			readers[i] = reader
		}
		//logger.Println("i is bad, err is", e)
	}

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil {
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}

	dec := NewDecoder(readers, writers, size)
	return &RsGetStream{dec}, nil
}
type decoder struct {
	readers []io.Reader
	writers []io.Writer
	enc reedsolomon.Encoder
	size int64
	cache []byte
	cacheSize int
	total int64
}

func NewDecoder(readersq []io.Reader, writersq []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{readersq, writersq, enc, size, nil, 0, 0}
}

func (d *decoder)Read(p []byte)(n int, err error) {
	if d.cacheSize == 0 {
		e := d.getData()
		if e != nil {
			return 0, e
		}
	}
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

func (d *decoder)getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			n, e := io.ReadFull(d.readers[i], shards[i])
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	e := d.enc.Reconstruct(shards)
	if e != nil {
		return e
	}
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}

	for i:=0;i<DATA_SHARDS;i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}

func (s *RsGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

func (s *RsGetStream)Seek(offset int64, whence int) (int64, error){
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}