package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bulldogcreative/goos/pkg/goos"
)

func main() {
	fmt.Println("Starting")
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("aws_access_key_id"), os.Getenv("aws_secret_access_key"), ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-2"),
	}
	sess := session.New(s3Config)

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "goos")
	if e == nil {
		log.SetOutput(logwriter)
	}

	handler := goos.S3Handler(sess, os.Getenv("aws_bucket"))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
