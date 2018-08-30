/**
 *  author: lim
 *  data  : 18-8-30 下午10:57
 */

package proto

const (
	C uint8 = iota
	D
	Q
	C_Q
)

type Request struct {
	Cmds []uint8
}

type Response struct {
	Max    uint64
	Active map[uint64]uint8
}
