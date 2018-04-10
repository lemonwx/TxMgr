/**
 *  author: lim
 *  data  : 18-4-10 下午9:48
 */

package main

import (
	"net"
	"net/rpc"
	"os"
	"fmt"
	"net/http"
	"sync/atomic"
	"sync"

	"github.com/lemonwx/log"
)

var addr string = "192.168.1.4:1235"
var NextV uint64
var vInuse map[uint64]string
var baseVersion uint64 = 1
var lock sync.RWMutex


type VSeq struct {
}

func setupLogger() {
	f, err := os.Create("xsql.log")
	if err != nil {
		fmt.Println("touch log file xsql.log failed: %v", err)
	}
	log.NewDefaultLogger(f)
	log.SetLevel(log.DEBUG)
	log.Debug("this is vseq's log")
}

func main() {
	vInuse = make(map[uint64]string, 1024)

	vSeq := new(VSeq)
	rpc.Register(vSeq)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(l)
	}

	http.Serve(l, nil)
}


func (v *VSeq) NextV(args uint8, reply *uint64) error {
	tmp := atomic.AddUint64(&baseVersion, 1)
	*reply = tmp
	lock.Lock()
	vInuse[tmp] = "test"
	lock.Unlock()
	return nil
}

func (v *VSeq) VInUser(args uint8, reply *[]uint64) error {
	lock.Lock()
	ret := make([]uint64, len(vInuse))
	idx := 0
	for k, _ := range vInuse {
		ret[idx] = k
		idx += 1
	}
	lock.Unlock()
	*reply = ret
	return nil
}

func (v *VSeq) Release(args uint64, reply *bool) error {
	lock.Lock()
	delete(vInuse, uint64(args))
	*reply = true
	lock.Unlock()
	return nil
}