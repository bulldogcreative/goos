package goos

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Handler will return the handler we need
func S3Handler(s *session.Session, bucket string) http.HandlerFunc {
	svc := s3.New(s)
	fn := func(w http.ResponseWriter, r *http.Request) {

		url, err := url.QueryUnescape(r.URL.String())
		if err != nil {
			log(r.RemoteAddr, r.URL.String(), "404")
			http.NotFound(w, r)
		}

		input := &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(url),
		}
		result, err := svc.GetObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchKey:
					log(r.RemoteAddr, url, "404")
					http.NotFound(w, r)
				default:
					w.Write([]byte(aerr.Error()))
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				w.Write([]byte(err.Error()))
			}

			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(*result.ContentLength, 10))
		w.Header().Set("Last-Modified", result.LastModified.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		w.Header().Set("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))
		w.Header().Set("Cache-Control", "max-age:290304000, public")
		w.Header().Set("Etag", *result.ETag)

		// fmt.Println(result)

		bf := new(bytes.Buffer)
		bf.ReadFrom(result.Body)

		w.Write([]byte(bf.String()))

		// Print Request Details
		log(r.RemoteAddr, url, "200")
	}

	return http.HandlerFunc(fn)
}

func log(remote string, url string, status string) {
	fmt.Println("[" + time.Now().Format(time.RFC3339) + "] [" + remote + "] " + "[" + url + "] " + "[" + status + "]")
}
