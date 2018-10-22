package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/vharitonsky/iniflags"
)

var (
	addr       = flag.String("addr", ":8080", "Address to listen with SOCKS5")
	serverAddr = flag.String("server", "", "The address of the remote server")
	sshUser    = flag.String("ssh-user", "lantern", "SSH user on the remote server")
	sshKeyFile = flag.String("ssh-key-file", "", "SSH key file to authenticate on the remote server")
)

func main() {
	iniflags.Parse()
	s, err := NewRemoteServer(*serverAddr, *sshUser, *sshKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	tunnel(*addr, s)
}

func tunnel(listenAdddr string, remoteServer *Server) error {
	s := socks{
		Dial: func(ctx context.Context, net, addr string) (net.Conn, error) {
			return remoteServer.ForwardTo(addr)
		},
	}
	return s.Serve(listenAdddr)
}
