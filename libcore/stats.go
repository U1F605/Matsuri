package libcore

import (
	"net"
	"sync/atomic"
)

type AppStats struct {
	Uid          int32
	TcpConn      int32
	UdpConn      int32
	TcpConnTotal int32
	UdpConnTotal int32

	Uplink        int64
	Downlink      int64
	UplinkTotal   int64
	DownlinkTotal int64

	DeactivateAt int32

	NekoConnectionsJSON string
}

type appStats struct {
	tcpConn      int32
	udpConn      int32
	tcpConnTotal uint32
	udpConnTotal uint32

	uplink        uint64
	downlink      uint64
	uplinkTotal   uint64
	downlinkTotal uint64

	deactivateAt int64
}

type TrafficListener interface {
	UpdateStats(t *AppStats)
}

type statsConn struct {
	net.Conn
	uplink   *uint64
	downlink *uint64
}

func (c *statsConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	defer atomic.AddUint64(c.uplink, uint64(n))
	return
}

func (c *statsConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	defer atomic.AddUint64(c.downlink, uint64(n))
	return
}

type myStats struct {
	uplink   *uint64
	downlink *uint64
}
