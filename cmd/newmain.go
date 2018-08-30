/**
 *  author: lim
 *  data  : 18-8-30 下午10:54
 */

package main

import (
	"net"

	"encoding/binary"

	"github.com/lemonwx/TxMgr/proto"
	"github.com/lemonwx/log"
)

var (
	addr  = "192.168.1.102:1235"
	queue = make(chan *proto.Request, 1024)

	max    uint64 = 1000
	active        = map[uint64]uint8{}
)

func main() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go handleCmd()
	handleConn(l)
}

func handleCmd() {
	for {
		req := <-queue
		resps := map[int]*proto.Response{}
		hasQ := false
		size := len(req.Cmds)

		for idx, cmd := range req.Cmds {
			switch cmd {
			case proto.C:
				max += 1
				resps[idx] = &proto.Response{Max: max}
			case proto.Q:
				hasQ = true
			case proto.D:
				delete(active)
			case proto.C_Q:
				hasQ = true
			}
		}
	}
}

func handleConn(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go serve(conn)
	}
}

func serve(conn net.Conn) {
	for {
		header := []byte{0, 0, 0, 0}
		if _, err := conn.Read(header); err != nil {
			log.Errorf("read proto header failed: %v", err)
			return
		}

		size := binary.LittleEndian.Uint32(header)
		content := make([]byte, size)
		if _, err := conn.Read(content); err != nil {
			log.Errorf("read proto content failed: %v", err)
			return
		}

		for content

		mergeReq := &proto.Request{
			Cmds: content,
		}

		queue <- mergeReq
	}
}
