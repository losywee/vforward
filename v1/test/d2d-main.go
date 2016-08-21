package main

import (
	"github.com/456vv/vforward/v1"
    "net"
    "log"
    "time"
    "flag"
    "fmt"
)

var fNetwork = flag.String("Network", "tcp", "网络地址类型")

var fALocal = flag.String("ALocal", "0.0.0.0", "A端本地发起连接地址")
var fALemote = flag.String("ALemote", "", "A端远程请求连接地址 (format \"12.13.14.15:123\")")

var fBLocal = flag.String("BLocal", "0.0.0.0", "B端本地发起连接地址")
var fBLemote = flag.String("BLemote", "", "B端远程请求连接地址 (format \"22.23.24.25:234\")")

var fTryConnTime = flag.Duration("TryConnTime", time.Millisecond*500, "尝试或发起连接时间，可能一方不在线，会一直尝试连接对方。单位：ns, us, ms, s, m, h")
var fTimeout = flag.Duration("Timeout", time.Second*5, "转发连接时候，请求远程连接超时。单位：ns, us, ms, s, m, h")
var fMaxConn = flag.Int("MaxConn", 500, "限制连接最大的数量")
var fKeptIdeConn = flag.Int("KeptIdeConn", 2, "保持一方连接数量，以备快速互相连接。")
var fReadBufSize = flag.Int("ReadBufSize", 4096, "交换数据缓冲大小。单位：字节")


//commandline:d2d-main.exe -ALemote 127.0.0.1:1201 -BLemote 127.0.0.1:1202
func main(){
    flag.Parse()
    if flag.NFlag() == 0 {
        flag.PrintDefaults()
        return
    }
    var err error
    if *fALemote == "" || *fBLemote == "" {
        log.Printf("地址未填，A端远程地址 %q, B端远程地址 %q", *fALemote, *fBLemote)
        return
    }
    rtcpaddr1, err := net.ResolveTCPAddr(*fNetwork, *fALemote)
    if err != nil {
        log.Println(err)
        return
    }
    rtcpaddr2, err := net.ResolveTCPAddr(*fNetwork, *fBLemote)
    if err != nil {
        log.Println(err)
        return
    }

    addra := &vforward.Addr{
        Network:*fNetwork,
        Local: &net.TCPAddr{IP: net.ParseIP(*fALocal),Port: 0,},
        Remote: rtcpaddr1,
    }
    addrb := &vforward.Addr{
        Network:*fNetwork,
        Local: &net.TCPAddr{IP: net.ParseIP(*fBLocal),Port: 0,},
        Remote: rtcpaddr2,
    }

	dd := &vforward.D2D{
        TryConnTime: *fTryConnTime,             // 尝试或发起连接时间，可能一方不在线，会一直尝试连接对方。
        MaxConn: *fMaxConn,                     // 限制连接最大的数量
        KeptIdeConn: *fKeptIdeConn,             // 保持一方连接数量，以备快速互相连接。
        Timeout: *fTimeout,                     // 发起连接超时
        ReadBufSize: *fReadBufSize,             // 交换数据缓冲大小
    }
    defer dd.Close()
    dds, err := dd.Transport(addra, addrb)
    if err != nil {
        log.Println(err)
        return
    }
    exit := make(chan bool, 1)
    go func(){
        defer func(){
            dds.Close()
            exit <- true
            close(exit)
        }()
        log.Println("D2D启动了")

        var in0 string
        for err == nil  {
            log.Println("输入任何字符可以退出D2D!")
            fmt.Scan(&in0)
            if in0 != "" {
                return
            }
        }
    }()
    defer dds.Close()
    err = dds.Swap()
    if err != nil {
        log.Println("错误：%s", err)
    }
    <-exit
    log.Println("D2D退出了")
}
