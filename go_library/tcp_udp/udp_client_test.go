package tcp_udp

import (
	"log"
	"net"
	"os"
	"testing"
)

func TestUdpClient(t *testing.T) {
	server, err := net.ResolveUDPAddr("udp", "10.0.10.11:61153")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		log.Fatalf("listen error %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("movie.wangzhaoju2.com"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}
}
