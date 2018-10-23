package main

import (
	"io/ioutil"
	"log"
	"net"
	"sync/atomic"

	"golang.org/x/crypto/ssh"
)

type remoteServer struct {
	config *ssh.ClientConfig
	addr   string
	client atomic.Value // *ssh.Client
}

func NewRemoteServer(
	addr string,
	user string,
	keyFile string) (*remoteServer, error) {
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
	return &remoteServer{config: config, addr: addr}, nil
}

func (s *remoteServer) ForwardTo(addr string) (conn net.Conn, err error) {
	var redialSSH bool
	for i := 0; i < 2; i++ {
		client := s.client.Load()
		if client == nil || redialSSH {
			log.Printf("Creating new SSH connection to %s", s.addr)
			client, err = ssh.Dial("tcp", s.addr, s.config)
			if err != nil {
				return
			}
			s.client.Store(client)
		}
		conn, err = client.(*ssh.Client).Dial("tcp", addr)
		if err != nil {
			redialSSH = true
			continue
		}
		return conn, nil
	}
	return nil, err
}
