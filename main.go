package main

import (
	"context"
	"flag"
	"net"
	"sync"

	"github.com/vharitonsky/iniflags"
)

var (
	socks5Port = flag.String("socks5-port", "8080", "The port on which the local SOCKS5 server is listen.")
	remotePort = flag.String("remote-port", "", "The port to access on the remote servers. Leave it empty to allow the remote port being passed as SOCKS5 user name on a per-request basis.")
	sshUser    = flag.String("ssh-user", "", "User name on the remote servers.")
	sshKeyFile = flag.String("ssh-key-file", "", "The path of the private key file to authenticate the user on the remote servers.")
)

func main() {
	var mx sync.Mutex
	remoteServers := make(map[string]*remoteServer)

	iniflags.Parse()
	s := socks{
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			mx.Lock()
			s, exists := remoteServers[addr]
			if !exists {
				s, err = NewRemoteServer(addr, *sshUser, *sshKeyFile)
				if err != nil {
					mx.Unlock()
					return
				}
				remoteServers[addr] = s
			}
			mx.Unlock()
			// Passed as username in SOCKS5 request
			remotePort := ctx.Value(ctxKeyRemotePort).(string)
			return s.ForwardTo(net.JoinHostPort("127.0.0.1", remotePort))
		},
		remotePort: *remotePort,
	}
	s.Serve(net.JoinHostPort("127.0.0.1", *socks5Port))
}
