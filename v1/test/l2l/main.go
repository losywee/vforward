package main

import (
	"github.com/456vv/vforward/v1"
    "flag"
    "net"
    "time"
    "log"
    "fmt"
)
var fNetwork = flag.String("Network", "tcp", "网络地址类型")
var fALocal = flag.String("ALocal", "", "A本地监听网卡IP地址 (format \"12.13.14.15:123\")")
var fBLocal = flag.String("BLocal", "", "B本地监听网卡IP地址 (format \"22.23.24.25:234\")")

var fTimeout = flag.Duration("Timeout", time.Second*5, "转发连接时候，请求远程连接超时。单位：ns, us, ms, s, m, h")
var fMaxConn = flag.Int("MaxConn", 500, "限制连接最大的数量")
var fKeptIdeConn = flag.Int("KeptIdeConn", 2, "保持一方连接数量，以备快速互相连接。")
var fReadBufSize = flag.Int("ReadBufSize", 4096, "交换数据缓冲大小。单位：字节")


//commandline:l2l-main.exe -ALocal 127.0.0.1:1201 -BLocal 127.0.0.1:1202
func main(){
    flag.Parse()
    if flag.NFlag() == 0 {
        flag.PrintDefaults()
        return
    }
    var err error
    if *fALocal == "" || *fBLocal == "" {
        log.Printf("地址未填，A监听地址 %q, B监听地址 %q", *fALocal, *fBLocal)
        return
    }
    addr1, err := net.ResolveTCPAddr(*fNetwork, *fALocal)
    if err != nil {
        log.Println(err)
        return
    }
    addr2, err := net.ResolveTCPAddr(*fNetwork, *fBLocal)
    if err != nil {
        log.Println(err)
        return
    }
    addra := &vforward.Addr{
        Network:*fNetwork,
        Local: addr1,
    }
    addrb := &vforward.Addr{
        Network:*fNetwork,
        Local: addr2,
    }
    ll := &vforward.L2L{
        MaxConn: *fMaxConn,                     // 限制连接最大的数量
        KeptIdeConn: *fKeptIdeConn,             // 保持一方连接数量，以备快速互相连接。
        ReadBufSize: *fReadBufSize,             // 交换数据缓冲大小
    }
    lls, err := ll.Transport(addra, addrb)
    if err != nil {
        log.Println(err)
        return
    }
    exit := make(chan bool, 1)
    go func(){
        defer func(){
            lls.Close()
            exit <- true
            close(exit)
        }()
        log.Println("L2L启动了")

        var in0 string
        for err == nil  {
            log.Println("输入任何字符可以退出L2D!")
            fmt.Scan(&in0)
            if in0 != "" {
                return
            }
        }
    }()
    defer lls.Close()
    err = lls.Swap()
    if err != nil {
        log.Println("错误：%s", err)
    }
    <-exit
    log.Println("L2L退出了")
}