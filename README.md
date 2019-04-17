# goos

_go object storage_ proxy.

## Example Usage

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bulldogcreative/goos/pkg/goos"
)

func main() {
	fmt.Println("Starting")

	g := &goos.Goos{
		KeyID:    os.Getenv("aws_access_key_id"),
		Secret:   os.Getenv("aws_secret_access_key"),
		Endpoint: "https://nyc3.digitaloceanspaces.com",
		Region:   "us-east-2",
		Bucket:   os.Getenv("aws_bucket"),
	}

	// logwriter, e := syslog.New(syslog.LOG_NOTICE, "goos")
	// if e == nil {
	// 	log.SetOutput(logwriter)
	// }

	handler := g.Handler()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

## nginx config

```nginx
server {
	server_name cdn.example.com;

	location / {
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_pass http://127.0.0.1:8080;
	}
}
```

## Response Headers

[Last-Modified](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified)
