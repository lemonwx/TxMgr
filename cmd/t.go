/**
 *  author: lim
 *  data  : 18-4-10 下午10:08
 */

package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	//a := []byte{1,2}
	str := string(123)
	fmt.Println(str)


	cli, err := rpc.DialHTTP("tcp", "192.168.1.4:1235")
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	var reply []byte
	err = cli.Call("VSeq.NextV", uint8(0), &reply)
	if err != nil {
		panic(err)
	}

	fmt.Println(27, reply)

	var re [][]byte
	err = cli.Call("VSeq.VInUser", uint8(0), &re)
	if err != nil {
		panic(err)
	}
	fmt.Println(re)

	/*
		var rel bool
		err = cli.Call("VSeq.Release", reply, &rel)
		if err != nil {
			panic(err)
		}
		fmt.Println(rel)
	*/
}
