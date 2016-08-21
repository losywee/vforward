package vforward

import (
    "testing"
    "net"
    "time"
    "fmt"
    "bytes"
    "log"
    "os"
)
//如果测试出现错误提示，可能是你的电脑响应速度没跟上。由于关闭网络需要一段时间，所以我设置了延时。


func Test_D2D_0(t *testing.T){
    var exit bool

    addra := &Addr{
        Network:"tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 0,},
        Remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1201,},
    }
    addrb := &Addr{
        Network:"tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 0,},
        Remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1202,},
    }

    //服务器
    al, err := net.Listen(addra.Remote.Network(), addra.Remote.String())
    if err != nil {
        t.Fatal(err)
    }
    bl, err := net.Listen(addrb.Remote.Network(), addrb.Remote.String())
    if err != nil {
        al.Close()
        t.Fatal(err)
    }
    defer al.Close()
    defer bl.Close()
    funcgo := make(chan bool)
    go func(){
        funcgo<-true
        for  {
            al.Accept()
            if exit {
                return
            }
        }
    }()
    <-funcgo
    go func(){
        funcgo<-true
        for  {
            bl.Accept()
            if exit {
                return
            }
        }
    }()
    <-funcgo

    //客户端
    dd := &D2D{
        TryConnTime: time.Millisecond,
        MaxConn:5,
        KeptIdeConn:4,
        ReadBufSize:1024,
        ErrorLog:log.New(os.Stdout, "*", log.LstdFlags),
    }
    defer dd.Close()
    defer dd.Close()
    dds, err := dd.Transport(addra, addrb)
    if err != nil {
        t.Fatalf("发起请求失败：%s", err)
    }
    exitdd := make(chan bool, 1)
    go func(){
        defer func(){
            dd.Close()
            exitdd<-true
        }()
        time.Sleep(time.Second)
        if dds.ConnNum() != dd.MaxConn {
            t.Logf("连接数未达预计数量。返回为：%d，预计为：%d", dds.ConnNum(), dd.MaxConn)
        }
        dds.Close()
        dd.Close()
        time.Sleep(time.Second*3)
        //由于是异步关闭连接，所以需要等一会
        if dds.ConnNum() != 0 {
            t.Logf("还有连接没有被关闭。返回为：%d，预计为：0", dds.ConnNum())
        }
    }()
    defer dds.Close()
    defer dds.Close()
    dds.Swap()
    exit = <-exitdd
}


func Test_L2D_0(t *testing.T){
    var exit bool
    dial := &Addr{
        Network:"tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 0,},
        Remote: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1203,},
    }
    listen := &Addr{
        Network:"tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1204,},
    }

    //服务器
    al, err := net.Listen(dial.Network, dial.Remote.String())
    if err != nil {t.Fatal(err)}
    defer al.Close()
    serverRun := make(chan bool)
    go func(){
        serverRun<- true
        for  {
            al.Accept()
            if exit {
                return
            }
        }
    }()
    <-serverRun

    funexit := make(chan bool)
    dialRun := make(chan bool)
    go func(){
        dialRun<-true
        for {
            select {
            	case <-funexit:
            		return
            	default:
                    time.Sleep(time.Millisecond)
                	conn, err := net.Dial(listen.Network, listen.Local.String())
                    if err != nil {
                        continue
                    }
                    defer conn.Close()
            }
        }
    }()
    <-dialRun

    //客户端
    ld := &L2D{
        MaxConn:5,
        ReadBufSize:1024,
        ErrorLog:log.New(os.Stdout, "*", log.LstdFlags),
    }
    defer ld.Close()
    defer ld.Close()

    lds, err := ld.Transport(dial, listen)
    if err != nil {
        t.Fatalf("发起请求失败：%s", err)
    }
    _, err = ld.Transport(dial, listen)
    if err == nil {
        t.Fatalf("不应该可以重复调用")
    }

    exitld := make(chan bool)
    go func(){
        defer func(){
            lds.Close()
            exitld <- true
        }()
        time.Sleep(time.Second*2)
        if lds.ConnNum() != ld.MaxConn {
            t.Logf("连接数未达预计数量。返回为：%d，预计为：%d", lds.ConnNum(), ld.MaxConn)
        }
        lds.Close()
        //异步关闭需要等一会
        time.Sleep(time.Second*4)
        if lds.ConnNum() != 0 {
            t.Logf("还有连接没有被关闭。返回为：%d，预计为：0",  lds.ConnNum())
        }
    }()
    defer lds.Close()
    defer lds.Close()
    lds.Swap()
    funexit <- true
    exit = <-exitld
}


func Test_L2D_1(t *testing.T){
    dial := &Addr{
        Network:"udp",
        Local: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"),Port: 0,},
        Remote: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1205,},
    }
    listen := &Addr{
        Network:"udp",
        Local: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1206,},
    }

    //服务器
    server, err := net.ListenPacket(dial.Network, dial.Remote.String())
    if err != nil {t.Fatal(err)}
    defer server.Close()
    funcgo := make(chan bool)
    go func(){
        funcgo <- true
        for  {
            b := make([]byte, 1024)
        	n, laddr, err := server.ReadFrom(b)
            if err != nil {
                return
            }
        	server.WriteTo(b[:n], laddr)
        }
    }()
    <-funcgo

    //客户端
    ld := &L2D{
        MaxConn:5,
        ReadBufSize:1024,
        Timeout:time.Second*2,
        ErrorLog:log.New(os.Stdout, "*", log.LstdFlags),
    }
    defer ld.Close()
    defer ld.Close()

    lds, err := ld.Transport(dial, listen)
    if err != nil {
        t.Fatalf("发起请求失败：%s", err)
    }
    _, err = ld.Transport(dial, listen)
    if err == nil {
        t.Fatalf("不应该可以重复调用")
    }
    exitld := make(chan bool)
    go func(t *testing.T){
        defer func (){
            exitld<-true
        }()

        p := []byte("1234")
        for i:=0;i<6;i++ {
            client, err := net.DialUDP(listen.Network, nil, listen.Local.(*net.UDPAddr))
            if err != nil {
                lds.Close()
                t.Fatalf("可能是服务器还没启动：%s", err)
            }
            client.Write(p)
            b := make([]byte, 1024)
            client.SetReadDeadline(time.Now().Add(time.Second))
            n, err := client.Read(b)

            if err == nil && !bytes.Equal(p, b[:n]) {
                fmt.Println("UDP不正常，发送和收到数据不一致。发送的是：%s，返回的是：%s", p, b[:n])
            }
            client.Close()
        }
        if lds.ConnNum() != ld.MaxConn {
            t.Logf("连接数未达预计数量。返回为：%d，预计为：%d", lds.ConnNum(), ld.MaxConn)
        }

        ld.Close()
        lds.Close()
        time.Sleep(time.Second)
        if lds.ConnNum() != 0 {
            t.Logf("还有连接没有被关闭。返回为：%d，预计为：0", lds.ConnNum())
        }
    }(t)

    defer lds.Close()
    lds.Swap()
    <- exitld
}

func Test_L2L_0(t *testing.T){
    addra := &Addr{
        Network: "tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1207,},
    }
    addrb := &Addr{
        Network: "tcp",
        Local: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"),Port: 1208,},
    }
    ll := &L2L{
        MaxConn:5,
        KeptIdeConn:2,
        ErrorLog:log.New(os.Stdout, "*", log.LstdFlags),
    }
    defer ll.Close()
    lls, err := ll.Transport(addra, addrb)
    if err != nil {
        t.Fatalf("发起请求失败：%s", err)
    }

    exitgo := make(chan bool,1)
    go func(){
        defer func(){
            lls.Close()
            exitgo<-true
        }()
        for i:=0;i<ll.MaxConn*2;i++ {
        	conna, err := net.Dial(addra.Network, addra.Local.String())
            if err != nil {
                t.Fatalf("客户端连接失败：%s", err)
            }
            defer conna.Close()

        	connb, err := net.Dial(addrb.Network, addrb.Local.String())
            if err != nil {
                t.Fatalf("客户端连接失败：%s", err)
            }
            defer connb.Close()
            time.Sleep(time.Millisecond*100)
        }
        time.Sleep(time.Second)
        if lls.ConnNum() != ll.MaxConn {
            t.Logf("连接数未达预计数量。返回为：%d，预计为：%d", lls.ConnNum(), ll.MaxConn)
        }
        lls.Close()
        lls.Close()
    }()
    defer lls.Close()

    lls.Swap()
    <-exitgo

    time.Sleep(time.Second)
    if lls.ConnNum() != 0 {
        t.Logf("还有连接没有被关闭。返回为：%d，预计为：0", lls.ConnNum())
    }

    go func(){
        ll.Close()
        lls.Close()
        exitgo<- true
    }()
    lls.Swap()
    <-exitgo
}
