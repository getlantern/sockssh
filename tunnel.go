package main

import (
	"io/ioutil"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	sshConfig   *ssh.ClientConfig
	idleTimeout time.Duration
}

type Server struct {
	config *ssh.ClientConfig
	addr   string
	client *ssh.Client
}

func NewRemoteServer(
	addr string,
	user string,
	keyFile string) (*Server, error) {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &Server{config, addr, nil}, nil
}

func (s *Server) ForwardTo(addr string) (net.Conn, error) {
	if s.client == nil {
		if client, err := ssh.Dial("tcp", s.addr, s.config); err != nil {
			return nil, err
		} else {
			s.client = client
		}
	}
	remoteConn, err := s.client.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return remoteConn, nil
}
