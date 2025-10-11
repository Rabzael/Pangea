# Proxy

| | |
|-|-|
| Language | Go |
| Version | 1 |

A simple proxy that interpones between client and server.
- Handles HTTP forwarding and HTTPS tunneling
- Methods allowed: GET, PUT, HEAD, CONNECT (for https)

## How to use
#### Run as a **Go script**:
```
$ go run main.go <path/to/config/file.json>
```

#### Or build as a standalone executable:
```
$ go build -o bin/proxy
$ bin/proxy <path/to/config/file.json>
```

## Configuration file options
*Please find a sample in the 'resources' folder.*

| Option | Description | Example |
|--------|-------------|---------|
| Port | the port to listen | 5173 |