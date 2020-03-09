package main

import (
	"net"
)

func main(){

	list, _ := net.Listen("tcp4", "127.0.0.1:8888")

	defer list.Close()

	conn, _ := list.Accept()

	defer conn.Close()

	b := []byte("hello world...\n")

	_, _ = conn.Write(b)

	//_, _ = io.WriteString(os.Stdout, "\n===============\n| Hello World |\n===============\n\n")
}