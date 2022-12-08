# acceptor-http

A golang http server that listens on multiple ports.

Accepts connections on a set of specified ports - useful if you are behind a firewall and want to work out what is open without being rate limited by a service like http://portquiz.net/

#### Bind all ports from 1024-2048, stop if errors occur

```bash
go run acceptor.go
```

#### Bind all ports from 1-1024, stop if errors occur

```bash
go run acceptor.go -from=1 -until=1024
```

#### Bind all ports from 1-65535, continue even if there is an error

```bash
go run acceptor.go -from=1 -until=65535 -die=false
```

#### How to test with nmap

Note the T argument is the 'agressiveness'. Best to use T2 for portquiz.net

T4 will complete ports 1-65535 in about 3 minutes on my current connection (150up/down)

```bash
nmap -T4 127.0.0.1 -r -p1024-2048 -v
```