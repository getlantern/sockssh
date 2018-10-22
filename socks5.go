package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/armon/go-socks5"
)

const ctxKeyRemotePort = "RemotePort"

type socks struct {
	Dial       func(ctx context.Context, net, addr string) (net.Conn, error)
	remotePort string
}

func (s *socks) Serve(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Unable to listen: %q", err)
	}
	conf := &socks5.Config{
		Dial:     s.Dial,
		Rewriter: s,
	}
	if s.remotePort == "" {
		// Retrieve SOCKS5 user name as remote port on a per request basis
		conf.Credentials = alwaysValid{}
	}
	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("Unable to create SOCKS5 server: %v", err)
	}

	log.Printf("About to start SOCKS5 client proxy at %v", addr)
	return server.Serve(l)
}

func (s *socks) Rewrite(ctx context.Context, request *socks5.Request) (context.Context, *socks5.AddrSpec) {
	remotePort := s.remotePort
	if remotePort == "" {
		remotePort = request.AuthContext.Payload["Username"]
	}
	ctx = context.WithValue(ctx, ctxKeyRemotePort, remotePort)
	return ctx, request.DestAddr
}

type alwaysValid struct{}

func (s alwaysValid) Valid(user, password string) bool {
	return true
}
