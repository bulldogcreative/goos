package goos

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Goos holds the information to connect to an S3 compatible Object Storage Service
type Goos struct {
	KeyID    string
	Secret   string
	Endpoint string
	Region   string
	Bucket   string
}

func (g *Goos) session() *session.Session {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(g.KeyID, g.Secret, ""),
		Endpoint:    aws.String(g.Endpoint),
		Region:      aws.String(g.Region),
	}
	sess := session.New(s3Config)

	return sess
}

// Handler will return the handler we need
func (g *Goos) Handler() http.HandlerFunc {
	svc := s3.New(g.session())
	fn := func(w http.ResponseWriter, r *http.Request) {

		ip := r.Header.Get("X-Real-Ip")

		if r.URL.String() == "/" {
			logMessage(ip, r.URL.String(), "/")
			notFound(w)
			return
		}

		url, err := url.QueryUnescape(r.URL.String())
		if err != nil {
			logMessage(ip, r.URL.String(), "404")
			notFound(w)
			return
		}

		input := &s3.GetObjectInput{
			Bucket: aws.String(g.Bucket),
			Key:    aws.String(url),
		}
		result, err := svc.GetObject(input)
		if err != nil {
			logMessage(ip, url, "404")
			notFound(w)
			return
		}
		defer result.Body.Close()

		w.Header().Set("Content-Length", strconv.FormatInt(*result.ContentLength, 10))
		w.Header().Set("Last-Modified", result.LastModified.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		w.Header().Set("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))
		w.Header().Set("Cache-Control", "max-age:290304000")
		w.Header().Set("Etag", *result.ETag)

		_, err = io.Copy(w, result.Body)
		if err != nil {
			logMessage(ip, url, "500")
			notFound(w)
			return
		}

		logMessage(ip, url, "200")
	}

	return http.HandlerFunc(fn)
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "max-age:0, private")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found."))
}

func logMessage(remote string, url string, status string) {
	fmt.Println("[" + remote + "] " + "[" + url + "] " + "[" + status + "]")
	//log.Print("[" + remote + "] [" + url + "] [" + status + "]")
}
