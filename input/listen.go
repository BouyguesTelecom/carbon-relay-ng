package input

import (
	"bytes"
	"net"
)

func listen(addr string, handler Handler) error {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	log.Notice("listening on %v/tcp", laddr)
	go acceptTcp(l, handler)

	udp_addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	udp_conn, err := net.ListenUDP("udp", udp_addr)
	if err != nil {
		return err
	}
	log.Notice("listening on %v/udp", udp_addr)
	go acceptUdp(udp_conn, handler)

	return nil
}

func acceptTcp(l *net.TCPListener, handler Handler) {
	for {
		// wait for a tcp connection
		c, err := l.AcceptTCP()
		if err != nil {
			log.Error(err.Error())
			break
		}
		log.Debug("listen.go: tcp connection from %v", c.RemoteAddr())
		// handle the connection
		go acceptTcpConn(c, handler)
	}
}

func acceptTcpConn(c net.Conn, handler Handler) {
	defer c.Close()
	handler.Handle(c)
}

func acceptUdp(l *net.UDPConn, handler Handler) {
	buffer := make([]byte, 65535)
	for {
		// read a packet into buffer
		b, addr, err := l.ReadFrom(buffer)
		if err != nil {
			log.Error(err.Error())
			break
		}
		log.Debug("listen.go: udp packet from %v (length: %d)", addr, b)
		// handle the packet
		handler.Handle(bytes.NewReader(buffer[:b]))
	}
}
