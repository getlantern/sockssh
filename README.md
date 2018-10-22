sockssh
====

When you have a bunch of public facing servers to monitor, rather than opening several ports with complex authentication, a simpler way would be to rely on the handy tool you trust and use everyday - SSH. You probably heard or got used to `ssh -L`, i.e., [SSH tunneling](https://www.ssh.com/ssh/tunneling/example) for casual tasks, but how can your monitoring tool like Prometheus use it to fetch metrics from your thousands of servers?

This is how [sockssh](https://github.com/getlantern/sockssh) comes into play. It's a SOCKS5 server listening on the local port. When new proxy requests come in, it creates SSH connection to the destination server (if not already exists), extracts the port on the server to which you intend to connect as the SOCKS5 username, and establish a tunnel. The remote user and key file to authenticate is supplied as command line options.

# Usage

```sh
# Starts sockssh on the background
sockssh -socks5-port=8000 -ssh-user=ubuntu -ssh-key-file=/home/<user>/.ssh/id_rsa &
# Fetchs goroutine profile which serves on 127.0.0.1:4000 on the remote server
curl -x socks5://4000@localhost:8000 <remote-server>:22/debug/pprof/goroutine?debug=1

# If clients doesn't support SOCKS5 authentication, setting remote port as command line option
sockssh -socks5-port=8000 -ssh-user=ubuntu -ssh-key-file=/home/<user>/.ssh/id_rsa -remote-port=4000 &
# Note the remote port is omitted
curl -x socks5://localhost:8000 <remote-server>:22/debug/pprof/goroutine?debug=1
```

# License

Apache License 2.0
