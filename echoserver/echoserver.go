package echoserver

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strconv"

// 	"golang.org/x/net/http2"
// 	"golang.org/x/net/http2/h2c"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// 	"google.golang.org/grpc/peer"
// 	"google.golang.org/grpc/reflection"
// )

import (
	"net"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"net/http"
	"strconv"

	"k8s-study/version"

	"github.com/spf13/cobra"
)

//go:generate protoc -I.  --go_out=plugins=grpc:. ./hello.proto

var CmdEchoServer = &cobra.Command{
	Use:   "echoserver",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

func StartHttpServer(port int) {
	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: http.HandlerFunc(handler),
	}
	addr := "0.0.0.0:" + strconv.Itoa(port)
	fmt.Printf("listening on %s, http1.0, http1.1\n", addr)
	server.ListenAndServe()
}

func main(cmd *cobra.Command, args []string) {
	StartHttpServer(8001)
	// go StartHttp2CleartextServer(8002)
	// go StartHttp2TLSServer(8003)
	// go StartGrpcServer(8004)
	// go StartGrpcWithTLSServer(8005)
	// go StartTCPServer(8006)
	// StartUDPServer(8007)
}

func getClientIP(req *http.Request) string {
	if req.Header.Get("X-Envoy-External-Address") != "" {
		return req.Header.Get("X-Envoy-External-Address")
	}

	ra, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ra
}

func printf(w io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, msg, args...)
	fmt.Printf(msg, args...)
}

func handler(w http.ResponseWriter, req *http.Request) {
	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	printf(w, "Version: %s\n", version.Version)
	printf(w, "HostName: %s\n", name)

	printf(w, "\nRequest Info:\n")
	printf(w, "    content-length: %d\n", req.ContentLength)
	printf(w, "    remote address: %s\n", req.RemoteAddr)
	printf(w, "    realIP: %s\n", getClientIP(req))
	printf(w, "    method: %s\n", req.Method)
	printf(w, "    path: %s\n", req.URL.Path)
	printf(w, "    query: %s\n", req.URL.RawQuery)
	printf(w, "    request_version: %s\n", req.Proto)
	printf(w, "    uri: %s\n", req.RequestURI)
	printf(w, "    tls: %t\n", req.TLS != nil)

	printf(w, "\nHeaders:\n")

	for name, headers := range req.Header {
		for _, h := range headers {
			printf(w, "    %v: %v\n", name, h)
		}
	}

	if req.Header.Get("Kalm-Sso-Userinfo") != "" {
		printf(w, "\nKalm SSO:\n")
		claims, err := base64.RawStdEncoding.DecodeString(req.Header.Get("Kalm-Sso-Userinfo"))

		if err != nil {
			printf(w, "Base64 decode error: %s\n", err.Error())
		} else {
			var out bytes.Buffer
			prefix := "  "
			if err := json.Indent(&out, claims, prefix, "  "); err != nil {
				printf(w, "json indent error: %s\n", err.Error())
			} else {
				printf(w, "%s%s\n", prefix, string(out.Bytes()))
			}
		}
	}

	printf(w, "\nBody:\n")
	if req.Body != nil && req.ContentLength > 0 {
		defer req.Body.Close()
		bs, err := ioutil.ReadAll(req.Body)

		if err != nil {
			printf(w, "Read body error: %s\n", err.Error())
		} else {
			printf(w, "%s", bs)
		}
	} else {
		printf(w, "No Body\n")
	}

	printf(w, "\nEnvironment Variables:\n")
	for _, env := range os.Environ() {
		printf(w, "%s\n", env)
	}

	printf(w, "\n")
}
