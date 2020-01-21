package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const FMT_TIMESTAMP = "2006-01-02 15:04:05 MST"

type logWriter struct {}
func (w logWriter) Write(bytes []byte) (int, error) {
	ts := time.Now().Format(FMT_TIMESTAMP)
	return fmt.Printf("[%s] %s", ts, string(bytes))
}

type ConnectionHandler struct {}
func (c *ConnectionHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	var responseStr string
	var err error

	startTime := time.Now()
	contentType := req.Header.Get("Content-Type")

	responseStr = fmt.Sprintf("%s %s %s\n", req.Method, req.URL, req.Proto)
	responseStr += fmt.Sprintf("Host:           %s\n", req.Host)
	responseStr += fmt.Sprintf("Accept:         %s\n", req.Header.Get("Accept"))
	responseStr += fmt.Sprintf("Timestamp:      %s\n", time.Now().Format(FMT_TIMESTAMP))
	responseStr += fmt.Sprintf("User-Agent:     %s\n", req.UserAgent())
	if contentType != "" {
		responseStr += fmt.Sprintf("Content-Type: 	%s\n", contentType)
	}
	responseStr += fmt.Sprintf("Content-Length: %d\n", req.ContentLength)

	if req.Body != nil && req.ContentLength > 0 {
		bodyBuffer := make([]byte, req.ContentLength)
		_, _ = req.Body.Read(bodyBuffer)
		responseStr += fmt.Sprintf("%s", bodyBuffer)
	}

	resWriter.Header().Add("Content-Length", strconv.Itoa(len(responseStr)))
	resWriter.Header().Add("Cache-Control", "no-store")

	_, err = resWriter.Write([]byte(responseStr))

	resTime := time.Now().Sub(startTime).Milliseconds()

	log.Printf("%s %s - %s - %dms\n", req.Method, req.URL, req.RemoteAddr, resTime)

	if err != nil {
		log.Println(err)
	}

}

func main() {
	var port int
	var logger logWriter

	log.SetFlags(0)
	log.SetOutput(logger)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s [port]", os.Args[0])
	}

	port, err :=  strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Port must be an integer 0-65535, not '%s'", os.Args[1])
	}

	log.Printf("Starting server on port %d...", port)

	connHandler := ConnectionHandler{}
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           &connHandler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
