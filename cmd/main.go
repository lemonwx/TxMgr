/**
 *  author: lim
 *  data  : 18-4-10 下午9:48
 */

package main

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"github.com/lemonwx/VSequence/base"
	"github.com/lemonwx/log"
)

var addr string = "192.168.1.2:1235"
var NextV uint64
var vInuse map[uint64]uint8
var baseVersion uint64 = 1000
var lock sync.RWMutex

type VSeq struct {
}

func setupLogger() {
	/*
		f, err := os.Create("v.log")
		if err != nil {
		}
		log.NewDefaultLogger(f)*/

	log.NewDefaultLogger(os.Stdout)
	log.SetLevel(log.DEBUG)
	log.Debug("this is vseq's log")
}

func main() {

	setupLogger()
	vInuse = make(map[uint64]uint8, 1024)

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

	lock.Lock()
	baseVersion += 1
	vInuse[baseVersion] = 1
	*reply = baseVersion
	lock.Unlock()
	return nil
}

func (v *VSeq) InUseAndNext(args uint8, reply *base.UseAndNext) error {
	lock.Lock()
	baseVersion += 1
	ret := make(map[uint64]uint8, len(vInuse))
	for k, v := range vInuse {
		ret[k] = v
	}
	log.Debug(reply, *reply)
	reply.Next = baseVersion
	reply.InUse = ret
	vInuse[baseVersion] = 1
	lock.Unlock()
	return nil
}

func (v *VSeq) VInUser(args uint8, reply *map[uint64]uint8) error {
	lock.RLock()
	ret := make(map[uint64]uint8, len(vInuse))
	for k, v := range vInuse {
		ret[k] = v
	}
	*reply = ret
	lock.RUnlock()
	return nil
}

func (v *VSeq) Release(args uint64, reply *bool) error {

	lock.Lock()
	delete(vInuse, args)
	*reply = true
	lock.Unlock()
	return nil
}
