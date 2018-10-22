package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/armon/go-socks5"
)

type socks struct {
	Dial func(ctx context.Context, net, addr string) (net.Conn, error)
}

func (s *socks) Serve(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Unable to listen: %q", err)
	}
	conf := &socks5.Config{
		Dial: s.Dial,
	}
	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("Unable to create SOCKS5 server: %v", err)
	}

	log.Printf("About to start SOCKS5 client proxy at %v", addr)
	return server.Serve(l)
}
