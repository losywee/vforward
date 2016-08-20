# vforward [![Build Status](https://travis-ci.org/456vv/vforward.svg?branch=master)](https://travis-ci.org/456vv/vforward)
go/golang TCP/UDP port forwarding，端口转发，主动连接，被动连接，大多用于内网端口反弹。
<br/>
最近更新20160820：<a href="/v1/update.txt">update.txt</a>
<br/>
列表：
====================
    const DefaultReadBufSize int = 4096                                             // 默认交换数据缓冲大小
    type Addr struct {                                                      // 地址
        Network       string                                                        // 网络类型
        Local, Remote net.Addr                                                      // 本地，远程
    }
    type D2D struct {                                                       // D2D（内网to内网）
        TryConnTime time.Duration                                                   // 尝试或发起连接时间，可能一方不在线，会一直尝试连接对方。
        MaxConn     int                                                             // 限制连接最大的数量
        KeptIdeConn int                                                             // 保持一方连接数量，以备快速互相连接。
        ReadBufSize int                                                             // 交换数据缓冲大小
        ErrorLog    *log.Logger                                                     // 日志
    }
        func (dd *D2D) Close() error                                                // 关闭
        func (dd *D2D) Transport(a, b *Addr) (*D2DSwap, error)                      // 建立连接
    type D2DSwap struct {}                                                   // D2D交换数据
        func (dds *D2DSwap) Close() error                                           // 关闭
        func (dds *D2DSwap) ConnNum() int                                           // 当前连接数
        func (dds *D2DSwap) Swap() error                                            // 开始交换
    type L2D struct {                                                        // L2D（端口转发）
        MaxConn     int                                                             // 限制连接最大的数量
        ReadBufSize int                                                             // 交换数据缓冲大小
        Timeout     time.Duration                                                   // 发起连接超时
        ErrorLog    *log.Logger                                                     // 日志
    }
        func (ld *L2D) Close() error                                                // 关闭
        func (ld *L2D) Transport(raddr, laddr *Addr) (*L2DSwap, error)              // 建立连接
    type L2DSwap struct {}                                                    // L2D交换数据
        func (lds *L2DSwap) Close() error                                           // 关闭
        func (lds *L2DSwap) ConnNum() int                                           // 当前连接数
        func (lds *L2DSwap) Swap() error                                            // 开始交换
    type L2L struct {                                                         // L2L（内网to内网）
        MaxConn     int                                                             // 限制连接最大的数量
        KeptIdeConn int                                                             // 保持一方连接数量，以备快速互相连接。
        ReadBufSize int                                                             // 交换数据缓冲大小
        ErrorLog    *log.Logger                                                     // 日志
    }
        func (ll *L2L) Close() error                                                // 关闭
        func (ll *L2L) Transport(aaddr, baddr *Addr) (*L2LSwap, error)              // 建立连接
    type L2LSwap struct {}                                                    // L2L交换数据
        func (lls *L2LSwap) Close() error                                           // 关闭
        func (lls *L2LSwap) ConnNum() int                                           // 当前连接数
        func (lls *L2LSwap) Swap() error                                            // 开始交换
<br/>