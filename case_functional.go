package main

import (
	"crypto/tls"
	"time"
)

/*********************************************** 1 */

///**
//要有侦听的 IP 地址 Addr 和端口号 Port ，这两个配置选项是必填的
//然后，还有协议 Protocol 、 Timeout 和MaxConns 字段，这几个字段是不能为空的，但是有默认值的，比如，协议是 TCP，超时30秒 和 最大链接数1024个
//还有一个 TLS ，这个是安全链接，需要配置相关的证书和私钥。这个是可以为空的
//*/
//
//type Server struct {
//	Addr     string
//	Port     int
//	Protocol string
//	Timeout  time.Duration
//	MaxConns int
//	TLS      *tls.Config
//}
//
///**
//Go 语言不支持重载函数，所以，你得用不同的函数名来应对不同的配置选项
//所以，针对这样的配置，我们需要有多种不同的创建不同配置 Server 的函数签名
//*/
//
//func NewDefaultServer(addr string, port int) (*Server, error) {
//	return &Server{addr, port, "tcp", 30 * time.Second, 100, nil}, nil
//}
//
//func NewTLSServer(addr string, port int, tls *tls.Config) (*Server, error) {
//	return &Server{addr, port, "tcp", 30 * time.Second, 100, tls}, nil
//}
//
//func NewServerWithTimeout(addr string, port int, timeout time.Duration) (*Server, error) {
//	return &Server{addr, port, "tcp", timeout, 100, nil}, nil
//}
//
//func NewTLSServerWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*Server, error) {
//	return &Server{addr, port, "tcp", 30 * time.Second, maxconns, tls}, nil
//}

/*********************************************** 2 */
///**
//最常见的方式是使用一个配置对象
//*/
//
//type Config struct {
//	Protocol string
//	Timeout  time.Duration
//	MaxConns int
//	TLS      *tls.Config
//}
//
//type Server struct {
//	Addr string
//	Port int
//	Conf *Config
//}
//
//func NewServer(addr string, port int, conf *Config) (*Server, error) {
//	return &Server{addr, port, conf}, nil
//}

/*********************************************** 3 */
/**
Builder 模式
*/

//type Server struct {
//	Addr     string
//	Port     int
//	Protocol string
//	Timeout  time.Duration
//	MaxConns int
//	TLS      *tls.Config
//}
//
//// ServerBuilder 使用一个builder类来做包装
//type ServerBuilder struct {
//	Server
//}
//
//func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
//	sb.Server.Addr = addr
//	sb.Server.Port = port
//	//其它代码设置其它成员的默认值
//	return sb
//}
//
//func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
//	sb.Server.Protocol = protocol
//	return sb
//}
//
//func (sb *ServerBuilder) WithMaxConn(maxconn int) *ServerBuilder {
//	sb.Server.MaxConns = maxconn
//	return sb
//}
//
//func (sb *ServerBuilder) WithTimeOut(timeout time.Duration) *ServerBuilder {
//	sb.Server.Timeout = timeout
//	return sb
//}
//
//func (sb *ServerBuilder) WithTLS(tls *tls.Config) *ServerBuilder {
//	sb.Server.TLS = tls
//	return sb
//}
//
//func (sb *ServerBuilder) Build() Server {
//	return sb.Server
//}
//
///**
//	sb := ServerBuilder{}
//	server := sb.Create("127.0.0.1", 8080).WithProtocol("udp").WithMaxConn(1024).WithTimeOut(30 * time.Second).Build()
// */

/*********************************************** 4 */

/**
Functional Options

这组代码传入一个参数，然后返回一个函数，返回的这个函数会设置自己的 Server 参数。
*/

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}
func MaxConns(maxconns int) Option {
	return func(s *Server) {
		s.MaxConns = maxconns
	}
}
func TLS(tls *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}

func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {

	srv := Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}
	for _, option := range options {
		option(&srv)
	}
	//...
	return &srv, nil
}

func main() {
	/**
	直觉式的编程；
	高度的可配置化；
	很容易维护和扩展；
	自文档；
	新来的人很容易上手；
	没有什么令人困惑的事（是 nil 还是空）。
	*/
	s1, _ := NewServer("localhost", 1024)
	s2, _ := NewServer("localhost", 2048, Protocol("udp"))
	s3, _ := NewServer("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))
	println(s1, s2, s3)
}
