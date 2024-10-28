package gid

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"time"

	"ascale/pkg/conf/env"
	"ascale/pkg/net/ip"

	"github.com/bwmarrin/snowflake"
)

var workerID = int64(1)
var generator *snowflake.Node

func Init() (err error) {
	ipStr := env.IP
	if ipStr == "" {
		ipStr = ip.InternalIP()
	}

	if ipStr != "" {
		ipInt32 := ip.InetAtoN(ipStr)
		workerID = int64(ipInt32 % 1024)
	} else {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		workerID = int64(rnd.Intn(1023))
	}

	if generator, err = snowflake.NewNode(workerID); err != nil {
		return
	}
	return
}

func NewID() (ts int64) {
	return generator.Generate().Int64()
}

func NewIDInt32() (ts int32) {
	return int32(generator.Generate().Int64() % 100000000)
}

func EncodeInt64ToString(id int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(id))
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeStringToInt64(str string) (id int64, err error) {
	bytes, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return
	}
	id = int64(binary.LittleEndian.Uint64(bytes))
	return
}
