package main

import (
	"io/ioutil"
	"log"
	"net"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
)

type remoteServer struct {
	config     *ssh.ClientConfig
	addr       string
	client     atomic.Value // *ssh.Client
	closeTimer *time.Timer
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
	return &remoteServer{config: config, addr: addr, closeTimer: time.NewTimer(0)}, nil
}

func (s *remoteServer) ForwardTo(addr string, sshIdleClose time.Duration) (conn net.Conn, err error) {
	if !s.closeTimer.Stop() {
		// We do not drain the channel here to avoid race condition with the
		// goroutine created below. This branch is very unlikely to be hit in
		// reality though.
	}
	s.closeTimer.Reset(sshIdleClose)
	var redialSSH bool
	for i := 0; i < 2; i++ {
		client := s.client.Load()
		if client == nil || redialSSH {
			log.Printf("Creating new SSH connection to %s", s.addr)
			client, err = ssh.Dial("tcp", s.addr, s.config)
			if err != nil {
				return
			}
			go func() {
				_ = <-s.closeTimer.C
				log.Printf("Closing SSH connection to %s", s.addr)
				client.(*ssh.Client).Close()
			}()
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
