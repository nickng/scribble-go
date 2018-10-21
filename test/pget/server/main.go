package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rhu1/scribble-go-runtime/runtime/session2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2"
	"github.com/rhu1/scribble-go-runtime/runtime/transport2/tcp"
	"github.com/rhu1/scribble-go-runtime/test/pget/PGet/Basic"
	S "github.com/rhu1/scribble-go-runtime/test/pget/PGet/Basic/family_1/S_1to1"
)

func main() {
	K := flag.Int("K", 2, "Specify parameter K")
	Port := flag.Int("port", 8080, "Specify listening port")
	flag.Parse()

	protocol := Basic.New()
	S := protocol.New_family_1_S_1to1(*K, 1)

	ln, err := tcp.Listen(*Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()

	for i := 1; i <= *K; i++ {
		if i == 1 {
			if err := S.F_1to1and1toK_Accept(i, ln, new(HTTPFormatter)); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		} else {
			if err := S.F_1toK_not_1to1_Accept(i, ln, new(HTTPFormatter)); err != nil {
				log.Fatalf("cannot accept: %v", err)
			}
		}
	}

	S.Run(serverBody)
}

func serverBody(s *S.Init) S.End {
	s0 := s.F_1_Gather_Head()
	s1 := s0.F_1_Scatter_Res()
	s2 := s1.F_1toK_Gather_Get()
	sEnd := s2.F_1toK_Scatter_Res()
	return *sEnd
}

// HTTPFormatter is a server-side HTTP formatter.
type HTTPFormatter struct {
	c transport2.BinChannel
}

// Wrap wraps a server-side TCP connection.
func (f *HTTPFormatter) Wrap(c transport2.BinChannel) { f.c = c }

// Serialize emulates sending of a file requested.
func (f *HTTPFormatter) Serialize(m session2.ScribMessage) error {
	file := `Content of HTTP file`
	res := &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          ioutil.NopCloser(strings.NewReader(file)),
		ContentLength: int64(len(file)),
	}
	return res.Write(f.c)
}

// Deserialize emulates server reading an HTTP Request.
func (f *HTTPFormatter) Deserialize(m *session2.ScribMessage) error {
	_, err := http.ReadRequest(bufio.NewReader(f.c))
	if err != nil {
		return err
	}
	return nil
}
