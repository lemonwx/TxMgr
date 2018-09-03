/**
 *  author: lim
 *  data  : 18-9-3 下午8:39
 */

package proto

import "time"

const (
	Q uint8 = iota
	C
	D
	C_Q
)

type Request struct {
	Cmds   []uint8
	ToDels []uint64
	Resp   chan *Response
	Ts     time.Time
}

type Response struct {
	Maxs   []uint64
	Active map[uint64]bool
}
