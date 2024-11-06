package socks4

import (
	"io"
	"net"
	"net/url"
	"strconv"

	errors "github.com/bdandy/go-errors"
	"golang.org/x/net/proxy"
)

type socks4 struct {
	url    *url.URL
	dialer proxy.Dialer
}

type Error = errors.Error

func init() {
	proxy.RegisterDialerType("socks4", func(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
		return socks4{url: u, dialer: d}, nil
	})

	proxy.RegisterDialerType("socks4a", func(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
		return socks4{url: u, dialer: d}, nil
	})
}

func (s socks4) Dial(network, addr string) (c net.Conn, err error) {
	if network != "tcp" && network != "tcp4" {
		return nil, ErrorWrongNetWork
	}

	c, err = s.dialer.Dial(network, s.url.Host)
	if err != nil {
		return nil, ErrDialFailed.New().Wrap(err)
	}
	// close connection later if got an error
	defer func() {
		if err != nil && c != nil {
			c.Close()
		}
	}()

	host, port, err := s.parseAddr(addr)
	if err != nil {
		return nil, ErrWrongAddr.New(addr).Wrap(err)
	}

	ip := net.IPv4(0, 0, 0, 1)
	if !s.isSocks4a() {
		if ip, err = s.lookupAddr(host); err != nil {
			return nil, ErrorHostUnknown.New(host).Wrap(err)
		}
	}

	req, err := request{Host: host, Port: port, IP: ip, Is4a: s.isSocks4a()}.Bytes()
	if err != nil {
		return nil, ErrBuffer.New().Wrap(err)
	}

	var i int
	i, err = c.Write(req)
	if err != nil {
		return c, ErrIO.New().Wrap(err)
	} else if i < minRequestLen {
		return c, ErrIO.New().Wrap(io.ErrUnexpectedEOF)
	}

	var resp [8]byte
	i, err = c.Read(resp[:])
	if err != nil && err != io.EOF {
		return c, ErrIO.New().Wrap(err)
	} else if i != 8 {
		return c, ErrIO.New().Wrap(err)
	}

	switch resp[i] {
	case accessGranted:
		return c, nil
	case accessIdentFailed, accessIdentRequired:
		return c, ErrIndentRequired
	case accessRejected:
		return c, ErrConnRejected
	default:
		return c, ErrInvalidResponce.New(resp[1])
	}
}

func (s socks4) parseAddr(addr string) (host string, iport int, err error) {
	var port string
	host, port, err = net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	iport, err = strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}

	return
}

func (s socks4) isSocks4a() bool {
	return s.url.Scheme == "socks4a"
}

func (s socks4) lookupAddr(host string) (net.IP, error) {
	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return net.IP{}, err
	}

	return ip.IP.To4(), nil
}
