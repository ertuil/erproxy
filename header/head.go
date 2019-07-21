package header

// AddrInfo .
type AddrInfo struct{
	From string 
	Host string
	Port string
	Atyp byte // 0x01 IPv4 0x03 Daemon 0x04 IPv6
	CMD  byte // 0x01 TCP 0x02 UDP
}

// GetInfo Unbox AddrInfo
func (ad AddrInfo) GetInfo() (string,string,string,byte,byte) {
	return ad.From,ad.Host,ad.Port,ad.Atyp,ad.CMD
}

// SetInfo BoxAddrInfo
func (ad *AddrInfo) SetInfo(from,host,port string,atype,cmd byte) {
	ad.From = from
	ad.Host = host
	ad.Port = port
	ad.Atyp = atype
	ad.CMD = cmd
}