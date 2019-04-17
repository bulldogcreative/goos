package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bulldogcreative/goos"
)

type logger struct{}

func (l *logger) Print(input string) {
	fmt.Println(input)
}

func main() {
	fmt.Println("Starting")

	log := &logger{}

	g := &goos.Goos{
		KeyID:    os.Getenv("aws_access_key_id"),
		Secret:   os.Getenv("aws_secret_access_key"),
		Endpoint: "https://nyc3.digitaloceanspaces.com",
		Region:   "us-east-2",
		Bucket:   os.Getenv("aws_bucket"),
		Logger:   log,
	}

	// logwriter, e := syslog.New(syslog.LOG_NOTICE, "goos")
	// if e == nil {
	// 	log.SetOutput(logwriter)
	// }

	handler := g.Handler()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
