package conn

import (
	"net"
	"net/netip"
	"syscall"
)

type Addr struct {
	IP   [4]byte
	Port uint16
}

func (addr *Addr) ToNetIp() netip.AddrPort {
	return netip.AddrPortFrom(netip.AddrFrom4(addr.IP), addr.Port)
}

func (addr *Addr) FromNetIp(addrPort netip.AddrPort) {
	addr.IP = addrPort.Addr().As4()
	addr.Port = addrPort.Port()
}

func (addr *Addr) FromString(str string) error {
	var addrPort, err = netip.ParseAddrPort(str)
	if err != nil {
		return err
	}
	addr.FromNetIp(addrPort)
	return nil
}

func (addr *Addr) FromBytes(buf []byte) error {
	return addr.FromString(string(buf))
}

func (addr *Addr) String() string {
	return addr.ToNetIp().String()
}

func (addr *Addr) Bytes() []byte {
	return []byte(addr.String())
}

func (addr *Addr) ToUDPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   addr.IP[:],
		Port: int(addr.Port),
	}
}

func (addr *Addr) ToSysInet4Addr() *syscall.SockaddrInet4 {
	return &syscall.SockaddrInet4{
		Port: int(addr.Port),
		Addr: addr.IP,
	}
}
