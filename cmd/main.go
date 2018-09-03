package main

import (
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/lemonwx/TxMgr/proto"
	"github.com/lemonwx/log"
)

var (
	addr     string = "192.168.1.2:1235"
	max      uint64 = 1000
	active          = map[uint64]bool{}
	reqQueue        = make(chan *proto.Request, 1024)
)

type VSeq struct {
}

func setupLogger() {
	/*
		f, err := os.Create("v.log")
		if err != nil {
		}
		log.NewDefaultLogger(f)*/

	log.NewDefaultLogger(os.Stdout)
	log.SetLevel(log.ERROR)
	log.Debug("this is vseq's log")
}

func main() {

	setupLogger()

	vSeq := new(VSeq)
	rpc.Register(vSeq)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(l)
	}

	go handleReq()
	http.Serve(l, nil)
}

func handleReq() {
	for {
		req := <-reqQueue
		resp := &proto.Response{
			Maxs:   make([]uint64, 0),
			Active: make(map[uint64]bool),
		}
		hasQ := false
		for _, cmd := range req.Cmds {
			switch cmd {
			case proto.Q:
				hasQ = true
			case proto.C:
				max += 1
				active[max] = false
				resp.Maxs = append(resp.Maxs, max)
			case proto.C_Q:
				hasQ = true
				max += 1
				active[max] = false
				resp.Maxs = append(resp.Maxs, max)
			case proto.D:
				gtid := req.ToDels[0]
				req.ToDels = req.ToDels[1:]
				delete(active, gtid)
			default:
				log.Errorf("receive unexpected cmd: %v", cmd)
			}
		}

		if hasQ {
			resp.Active = active
		}

		log.Debugf("resp to client")
		req.Resp <- resp
	}
}

func (v *VSeq) PushReq(req *proto.Request, resp **proto.Response) error {
	req.Resp = make(chan *proto.Response, 1)
	log.Debugf("receive %d requests merge to response", len(req.Cmds))
	reqQueue <- req
	*resp = <-req.Resp
	log.Debugf("response %d merge requests", len(req.Cmds))
	return nil
}
