# acceptor

portquiz-like smol tcp-server

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

```bash
nmap -T4 127.0.0.1 -r -p1024-2048 -v
```