package vforward

import (
	"net"
    "errors"
    "io"
    "strings"
    "time"
)
const DefaultReadBufSize int = 4096         // 默认交换数据缓冲大小

//Addr 地址包含本地，远程
type Addr struct {
    Network string
	Local, Remote   net.Addr        // 本地，远程
}

func temporaryError(err error, delay time.Duration, maxDelay time.Duration)(time.Duration, bool) {
    if ne, ok := err.(net.Error); ok && ne.Temporary() {
    	if delay == 0 {
    		delay = 5 * time.Millisecond
    	} else {
    		delay *= 2
    	}
    	if max := 1 * maxDelay; delay > max {
    		delay = max
    	}
    	time.Sleep(delay)
    	return delay, true
    }
    return delay, false
}

func copyData(dst io.Writer, src io.ReadCloser, bufferSize int)(written int64, err error){
    defer src.Close()
    buf := make([]byte, bufferSize)
    written, err = io.CopyBuffer(dst, src, buf)
    return
}

func connectListen(addr *Addr) (interface{}, error) {
    switch addr.Network {
    	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
            return net.Listen(addr.Network, addr.Local.String())
    	case "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unixgram":
            return net.ListenPacket(addr.Network, addr.Local.String())
        default:
            if strings.HasPrefix(addr.Network, "ip:") && len(addr.Network) > 3 {
                return net.ListenPacket(addr.Network, addr.Local.String())
            }
    }
    return nil, errors.New("netForward: 监听地址类型是未知的")
}

func connectUDP(addr *Addr) (net.Conn, error) {
	switch addr.Network {
    	case "udp", "udp4", "udp6":
    		return net.DialUDP(addr.Network, nil, addr.Remote.(*net.UDPAddr))
    	case "ip", "ip4", "ip6":
    		return net.DialIP(addr.Network, nil, addr.Remote.(*net.IPAddr))
    	case "unix", "unixpacket", "unixgram":
    		return net.DialUnix(addr.Network, nil, addr.Remote.(*net.UnixAddr))
        default:
            if strings.HasPrefix(addr.Network, "ip:") && len(addr.Network) > 3 {
    	    	return net.DialIP(addr.Network, nil, addr.Remote.(*net.IPAddr))
            }
	}
    return nil, errors.New("netForward: 远程地址类型是未知的")
}
