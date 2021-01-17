# TalkTasy

A Golang library to securely talk with services over Mutual TLS (mTLS).

## Features

- [x] Support basic communication over HTTP Protocol.
- [x] Support basic communication over WebSocket Protocol.

## Example

### Server
Creating a server with `talktasy` is as easy as this:

```go
package main

import (
    "log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/tasylab/talktasy"
)

func main() {
	srv := talktasy.NewServer(
		"./ca.pem",
		"./server.crt",
		"./server.key",
	)
	log.Fatal(srv.Listen("localhost:8443", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Request verified! \nPath: %s\n", r.URL.Path)
    }))
}
```

### Client
Creating a client with `talktasy`:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tasylab/talktasy"
)

func main() {
	client := talktasy.NewClient(
		"./ca.pem",
		"./client.crt",
		"./client.key",
	)
	req, _ := http.NewRequest("GET", "https://localhost:8443", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	fmt.Printf("%s\n", body)
}
```