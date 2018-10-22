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
	for i := 0; i < 2; i++ {
		client := s.client.Load()
		if client == nil {
			if err = s.newSSH(); err != nil {
				return
			}
			// try again using the newly created SSH connection
			continue
		}
		conn, err = client.(*ssh.Client).Dial("tcp", addr)
		if err != nil {
			// Force creating a new SSH connection to retry, in case if the
			// current one was broken or closed by the server.
			s.client.Store(nil)
			continue
		}
		return conn, nil
	}
	return nil, err
}

func (s *remoteServer) newSSH() (err error) {
	log.Printf("Creating new SSH connection to %s", s.addr)
	client, err := ssh.Dial("tcp", s.addr, s.config)
	if err != nil {
		return err
	}
	s.client.Store(client)
	return nil
}
