package main

import (
	"flag"
	"fmt"
	"sync/atomic"

	"github.com/ahmadhabibi14/gnet-starter/protocol"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

type Server struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	network      string
	addr         string
	multicore    bool
	connected    int32
	disconnected int32
}

func NewServer(network string, port int, multicore bool) *Server {
	return &Server{
		network: network,
		addr: fmt.Sprintf(":%d", port),
		multicore: multicore,
	}
}

func (s *Server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof(
		"running server on %s with multicore=%t",
		fmt.Sprintf("%s://%s", s.network, s.addr), s.multicore,
	)
	s.eng = eng
	return
}

func (s *Server) OnShutdown(eng gnet.Engine) {
	logging.Infof("server close")
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	disconnected := atomic.AddInt32(&s.disconnected, 1)
	connected := atomic.AddInt32(&s.connected, -1)
	if connected == 0 {
		logging.Infof("all %d connections are closed", disconnected)
		action = gnet.None
	}
	return
}


func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(new(protocol.SimpleCodec))
	atomic.AddInt32(&s.connected, 1)
	out = []byte("sweetness\r\n")
	return
}

func (s *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	logging.Error(err)
	fmt.Printf("Received: %s\n", buf)

	action = gnet.None
	return
}

func main() {
	var port int
	var multicore bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore=true")
	flag.Parse()
	
	server := NewServer("tcp", port, multicore)
	protoAddr := server.network+"://"+server.addr

	err := gnet.Run(server, protoAddr, gnet.WithMulticore(multicore))
	logging.Infof("server exits with error: %v", err)
}
