# Proxy

| | |
|-|-|
| Language | Go |
| Version | 1 |

A simple proxy that interpones between client and server.
- Handles HTTP forwarding and HTTPS tunneling
- Methods allowed: GET, PUT, HEAD, CONNECT (for https)
- Logs each request/response as a JSON object in a log file

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
| Log_filename | name of the log file | "logging/logs.log" |

## Log file format
Each line in the log file is a JSON object with the following fields:
| Field         | Description                                 | Mandatory | Example                      |
|---------------|---------------------------------------------|-----------|------------------------------|
| time          | time of the request                         | true      | "2024-01-01T12:00:00Z"       |
| remote_addr   | IP address of the client                    | true      | "192.168.1.1"                |
| method        | HTTP method used                            | true      | "GET"                     |
| url           | requested URL                               | true      | "http://example.com/resource"|
| status        | HTTP response code from the server          | true      | 200                          |
| content_len   | length of the response content (in bytes)   | true      | 512                          |
| error         | error message if any occurred during request| false     | "timeout error"              |
| blocked       | whether the request was blocked by proxy (true/false)| true      | false                        |