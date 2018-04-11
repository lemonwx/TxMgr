/**
 *  author: lim
 *  data  : 18-4-10 下午9:48
 */

package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

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

func (v *VSeq) NextV(args uint8, reply *[]byte) error {
	tmp := atomic.AddUint64(&baseVersion, 1)
	//*reply = make([]byte, 8)
	//binary.BigEndian.PutUint64(*reply, tmp)
	*reply = []byte(strconv.FormatUint(tmp, 10))
	lock.Lock()
	vInuse[tmp] = "test"
	lock.Unlock()
	return nil
}

func (v *VSeq) VInUser(args uint8, reply *[][]byte) error {
	lock.Lock()
	ret := make([][]byte, len(vInuse))
	idx := 0
	for k, _ := range vInuse {
		ret[idx] = []byte(strconv.FormatUint(k, 10))
		idx += 1
	}
	lock.Unlock()
	*reply = ret
	return nil
}

func (v *VSeq) Release(args []byte, reply *bool) error {

	tmp := string(args)
	fmt.Println(tmp)
	version, err := strconv.ParseUint(tmp, 10, 64)
	if err != nil {
		return err
	}

	lock.Lock()
	delete(vInuse, version)
	*reply = true
	lock.Unlock()
	return nil
}
