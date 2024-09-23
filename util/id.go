package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
)

var (
	prefix string
	reqid  uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func RequestID() string {
	myid := atomic.AddUint64(&reqid, 1)
	return fmt.Sprintf("%s-%06d", prefix, myid)
}

func ClientID(addr string) string {
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	return fmt.Sprintf("%s/%s", addr, b64[0:10])
}
