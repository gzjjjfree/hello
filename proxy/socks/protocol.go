//go:build !confonly
// +build !confonly

package socks

import (
	"encoding/binary"
	"io"
	"fmt"
	"errors"

	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/buf"
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/common/protocol"
)

const (
	socks5Version = 0x05
	socks4Version = 0x04

	cmdTCPConnect    = 0x01
	cmdTCPBind       = 0x02
	cmdUDPAssociate  = 0x03
	cmdTorResolve    = 0xF0
	cmdTorResolvePTR = 0xF1

	socks4RequestGranted  = 90
	socks4RequestRejected = 91

	authNotRequired = 0x00
	// authGssAPI           = 0x01
	authPassword         = 0x02
	authNoMatchingMethod = 0xFF

	statusSuccess       = 0x00
	statusCmdNotSupport = 0x07
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),
	protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain),
)

type ServerSession struct {
	config        *ServerConfig
	address       net.Address
	port          net.Port
	clientAddress net.Address
}

func (s *ServerSession) handshake4(cmd byte, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	fmt.Println("in proxy-socks-protocol.go func (s *ServerSession) handshake4")
	if s.config.AuthType == AuthType_PASSWORD {
		writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0))
		return nil, errors.New("socks 4 is not allowed when auth is required")
	}

	var port net.Port
	var address net.Address

	{
		buffer := buf.StackNew()
		if _, err := buffer.ReadFullFrom(reader, 6); err != nil {
			buffer.Release()
			return nil, errors.New("insufficient header")
		}
		port = net.PortFromBytes(buffer.BytesRange(0, 2))
		address = net.IPAddress(buffer.BytesRange(2, 6))
		buffer.Release()
	}

	if _, err := ReadUntilNull(reader); /* user id */ err != nil {
		return nil, err
	}
	if address.IP()[0] == 0x00 {
		domain, err := ReadUntilNull(reader)
		if err != nil {
			return nil, errors.New("failed to read domain for socks 4a")
		}
		address = net.DomainAddress(domain)
	}

	switch cmd {
	case cmdTCPConnect:
		request := &protocol.RequestHeader{
			Command: protocol.RequestCommandTCP,
			Address: address,
			Port:    port,
			Version: socks4Version,
		}
		if err := writeSocks4Response(writer, socks4RequestGranted, net.AnyIP, net.Port(0)); err != nil {
			return nil, err
		}
		return request, nil
	default:
		writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0))
		return nil, fmt.Errorf("unsupported command: %v", cmd)
	}
}

func (s *ServerSession) auth5(nMethod byte, reader io.Reader, writer io.Writer) (username string, err error) {
	fmt.Println("in proxy-socks-protocol.go func (s *ServerSession) auth5")
	buffer := buf.StackNew()
	defer buffer.Release()
// 客户端请求第一个数据包，第二个字节是 NMETHODS： 表示客户端支持的认证方法的数量，根据 nMethod 读取相应的字节到 buffer
	if _, err = buffer.ReadFullFrom(reader, int32(nMethod)); err != nil {
		return "", errors.New("failed to read auth methods")
	}

	var expectedAuth byte = authNotRequired
	// 认证方式，这里只为 noauth 匿名认证
	if s.config.AuthType == AuthType_PASSWORD {
		expectedAuth = authPassword
	}
    // 因为匿名认证，所以下面的判定也注释掉
	if !hasAuthMethod(expectedAuth, buffer.BytesRange(0, int32(nMethod))) {
		writeSocks5AuthenticationResponse(writer, socks5Version, authNoMatchingMethod)
		return "", errors.New("no matching auth method")
	}

	// 先把版本及认证方式发送回去
	if err := writeSocks5AuthenticationResponse(writer, socks5Version, expectedAuth); err != nil {
		return "", errors.New("failed to write auth response")
	}

	// 下面是用户密码认证，不用所以注释掉
	
	if expectedAuth == authPassword {
		username, password, err := ReadUsernamePassword(reader)
		if err != nil {
			return "", errors.New("failed to read username and password for authentication")
		}

		if !s.config.HasAccount(username, password) {
			writeSocks5AuthenticationResponse(writer, 0x01, 0xFF)
			return "", errors.New("invalid username or password")
		}

		if err := writeSocks5AuthenticationResponse(writer, 0x01, 0x00); err != nil {
			return "", errors.New("failed to write auth response")
		}
		return username, nil
	}
	

	return "", nil
}

func (s *ServerSession) handshake5(nMethod byte, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	fmt.Println("in proxy-socks-protocol.go func (s *ServerSession) handshake5")
	var (
		username string
		err      error
	)
	if username, err = s.auth5(nMethod, reader, writer); err != nil {
		return nil, err
	}

	var cmd byte
	{
		buffer := buf.StackNew()
		if _, err := buffer.ReadFullFrom(reader, 3); err != nil {
			buffer.Release()
			return nil, errors.New("failed to read request")
		}
		cmd = buffer.Byte(1)
		buffer.Release()
	}

	request := new(protocol.RequestHeader)
	if username != "" {
		fmt.Println("in proxy-socks-protocol.go func (s *ServerSession) handshake5 username != nil : ", username)
		request.User = &protocol.MemoryUser{Email: username}
	}
	switch cmd {
	case cmdTCPConnect, cmdTorResolve, cmdTorResolvePTR:
		fmt.Println("in  handshake5 cmd = cmdTCPConnect, cmdTorResolve, cmdTorResolvePTR")
		// We don't have a solution for Tor case now. Simply treat it as connect command.
		request.Command = protocol.RequestCommandTCP
	case cmdUDPAssociate:
		fmt.Println("in  handshake5 cmd = cmdUDPAssociate")
		if !s.config.UdpEnabled {
			writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
			return nil, errors.New("UDP is not enabled")
		}
		request.Command = protocol.RequestCommandUDP
	case cmdTCPBind:
		fmt.Println("in  handshake5 cmd = cmdTCPBind")
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
		return nil, errors.New("TCP bind is not supported")
	default:
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
		return nil, errors.New("unknown command ")
	}

	request.Version = socks5Version
	
	addr, port, err := addrParser.ReadAddressPort(nil, reader)
	if err != nil {
		return nil, errors.New("failed to read address")
	}
	request.Address = addr
	request.Port = port

	responseAddress := s.address
	responsePort := s.port
	//nolint:gocritic // Use if else chain for clarity
	if request.Command == protocol.RequestCommandUDP {
		
		if s.config.Address != nil {
			fmt.Println("in  handshake5 protocol.RequestCommandUDP s.config.Address != nil")
			// Use configured IP as remote address in the response to UdpAssociate
			// 在对 UdpAssociate 的响应中使用配置的 IP 作为远程地址
			responseAddress = s.config.Address.AsAddress()
		} else if s.clientAddress == net.LocalHostIP || s.clientAddress == net.LocalHostIPv6 {
			fmt.Println("in  handshake5 protocol.RequestCommandUDP s.clientAddress == net.LocalHostIP || s.clientAddress == net.LocalHostIPv6")
			// For localhost clients use loopback IP
			responseAddress = s.clientAddress
		} else {
			fmt.Println("in  handshake5 protocol.RequestCommandUDP responseAddress = s.address")
			// For non-localhost clients use inbound listening address
			responseAddress = s.address
		}
	}
	if err := writeSocks5Response(writer, statusSuccess, responseAddress, responsePort); err != nil {
		return nil, err
	}
	fmt.Println("in  handshake5 request is : ", request)
	return request, nil
}

// Handshake performs a Socks4/4a/5 handshake.
// 握手执行 Socks4/4a/5 握手。直接按 Socks5 执行
func (s *ServerSession) Handshake(reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	fmt.Println("in proxy-socks-protocol.go func (s *ServerSession) Handshake")
	// buffer 取得一个 sync.Pool 的自定义内存池
	buffer := buf.StackNew()
	// 从 reader(里面是原conn) 中读取 6 个 字节 长度到 buffer ，并把 buffer 的标记 end 设为 6 
	// 其中第一个表示版本 Socks5 ，第二个表示认证方法数量，第三个表示认证方法，第5个表示连接类型，第七个是地址类型，后面的地址
	if _, err := buffer.ReadFullFrom(reader, 3); err != nil {
		buffer.Release()
		return nil, errors.New("insufficient header")
	}
// 读取的第一位 byte 数为 version 类型是 []byte
	version := buffer.Byte(0)
	if version != socks5Version {
		return nil, errors.New("Socks version is NOT Socks5")
	}
	var expectedAuth byte = authNotRequired
	// 先返回 Socks5 类型，不认证
	if err := writeSocks5AuthenticationResponse(writer, socks5Version, expectedAuth); err != nil {
		return nil, errors.New("failed to write auth response")
	}	
	//buffer.Release()

	//switch version {
	//case socks4Version:
	//	return s.handshake4(cmd, reader, writer)
	//case socks5Version:
	//	return s.handshake5(cmd, reader, writer)
	//default:
	//	return nil, errors.New("unknown Socks version: ")
	//}

	
	//var cmd byte
	//{
	//	buffer := buf.StackNew()
		//if _, err := buffer.ReadFullFrom(reader, 3); err != nil {
		//	buffer.Release()
		//	return nil, errors.New("failed to read request")
		//}
		//cmd = buffer.Byte(1)
		//buffer.Release()
	//}
	if _, err := buffer.ReadFullFrom(reader, 3); err != nil {
		buffer.Release()
		return nil, errors.New("insufficient header")
	}
	request := new(protocol.RequestHeader)
	//fmt.Println("in  handshake5 cmd: ", buffer.Byte(0), buffer.Byte(1), buffer.Byte(2), buffer.Byte(3), buffer.Byte(4), buffer.Byte(5), buffer.Byte(6) )
	cmd := buffer.Byte(4)
	buffer.Release()
	switch cmd {
	case cmdTCPConnect, cmdTorResolve, cmdTorResolvePTR:
		fmt.Println("in  handshake5 cmd = cmdTCPConnect, cmdTorResolve, cmdTorResolvePTR")
		// We don't have a solution for Tor case now. Simply treat it as connect command.
		request.Command = protocol.RequestCommandTCP
	case cmdUDPAssociate:
		fmt.Println("in  handshake5 cmd = cmdUDPAssociate")
		if !s.config.UdpEnabled {
			writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
			return nil, errors.New("UDP is not enabled")
		}
		request.Command = protocol.RequestCommandUDP
	case cmdTCPBind:
		fmt.Println("in  handshake5 cmd = cmdTCPBind")
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
		return nil, errors.New("TCP bind is not supported")
	default:
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
		return nil, errors.New("unknown command ")
	}

	request.Version = socks5Version
//根据第一个字节选择类型解析地址
	addr, port, err := addrParser.ReadAddressPort(nil, reader)
	if err != nil {
		return nil, errors.New("failed to read address")
	}
	// request 目标地址
	request.Address = addr
	request.Port = port
// response 来源地址
	responseAddress := s.address
	responsePort := s.port
// 先返回请求成功信号 statusSuccess 及 本代理地址，这里是 127.0.0.1:54321
	if err := writeSocks5Response(writer, statusSuccess, responseAddress, responsePort); err != nil {
		return nil, err
	}
	fmt.Println("in  handshake5 request is : ", request)
	return request, nil
	
}

// ReadUsernamePassword reads Socks 5 username/password message from the given reader.
// +----+------+----------+------+----------+
// |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
// +----+------+----------+------+----------+
// | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
// +----+------+----------+------+----------+
func ReadUsernamePassword(reader io.Reader) (string, string, error) {
	fmt.Println("in proxy-socks-protocol.go func ReadUsernamePassword")
	buffer := buf.StackNew()
	defer buffer.Release()

	if _, err := buffer.ReadFullFrom(reader, 2); err != nil {
		return "", "", err
	}
	nUsername := int32(buffer.Byte(1))

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, nUsername); err != nil {
		return "", "", err
	}
	username := buffer.String()

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return "", "", err
	}
	nPassword := int32(buffer.Byte(0))

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, nPassword); err != nil {
		return "", "", err
	}
	password := buffer.String()
	return username, password, nil
}

// ReadUntilNull reads content from given reader, until a null (0x00) byte.
func ReadUntilNull(reader io.Reader) (string, error) {
	fmt.Println("in proxy-socks-protocol.go func ReadUntilNull")
	b := buf.StackNew()
	defer b.Release()

	for {
		_, err := b.ReadFullFrom(reader, 1)
		if err != nil {
			return "", err
		}
		if b.Byte(b.Len()-1) == 0x00 {
			b.Resize(0, b.Len()-1)
			return b.String(), nil
		}
		if b.IsFull() {
			return "", errors.New("buffer overrun")
		}
	}
}

func hasAuthMethod(expectedAuth byte, authCandidates []byte) bool {
	for _, a := range authCandidates {
		if a == expectedAuth {
			return true
		}
	}
	return false
}

func writeSocks5AuthenticationResponse(writer io.Writer, version byte, auth byte) error {
	return buf.WriteAllBytes(writer, []byte{version, auth})
}

func writeSocks5Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	fmt.Println("in proxy-socks-protocol.go func writeSocks5Response")
	buffer := buf.New()
	defer buffer.Release()

	common.Must2(buffer.Write([]byte{socks5Version, errCode, 0x00 /* reserved */}))
	if err := addrParser.WriteAddressPort(buffer, address, port); err != nil {
		return err
	}

	return buf.WriteAllBytes(writer, buffer.Bytes())
}

func writeSocks4Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	fmt.Println("in proxy-socks-protocol.go func writeSocks4Response")
	buffer := buf.StackNew()
	defer buffer.Release()

	common.Must(buffer.WriteByte(0x00))
	common.Must(buffer.WriteByte(errCode))
	portBytes := buffer.Extend(2)
	binary.BigEndian.PutUint16(portBytes, port.Value())
	common.Must2(buffer.Write(address.IP()))
	return buf.WriteAllBytes(writer, buffer.Bytes())
}

func DecodeUDPPacket(packet *buf.Buffer) (*protocol.RequestHeader, error) {
	fmt.Println("in proxy-socks-protocol.go func DecodeUDPPacket")
	if packet.Len() < 5 {
		return nil, errors.New("insufficient length of packet.")
	}
	request := &protocol.RequestHeader{
		Version: socks5Version,
		Command: protocol.RequestCommandUDP,
	}

	// packet[0] and packet[1] are reserved
	if packet.Byte(2) != 0 /* fragments */ {
		return nil, errors.New("discarding fragmented payload.")
	}

	packet.Advance(3)

	addr, port, err := addrParser.ReadAddressPort(nil, packet)
	if err != nil {
		return nil, errors.New("failed to read UDP header")
	}
	request.Address = addr
	request.Port = port
	return request, nil
}

func EncodeUDPPacket(request *protocol.RequestHeader, data []byte) (*buf.Buffer, error) {
	fmt.Println("in proxy-socks-protocol.go func EncodeUDPPacket")
	b := buf.New()
	common.Must2(b.Write([]byte{0, 0, 0 /* Fragment */}))
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		b.Release()
		return nil, err
	}
	common.Must2(b.Write(data))
	return b, nil
}

type UDPReader struct {
	reader io.Reader
}

func NewUDPReader(reader io.Reader) *UDPReader {
	return &UDPReader{reader: reader}
}

func (r *UDPReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b := buf.New()
	if _, err := b.ReadFrom(r.reader); err != nil {
		return nil, err
	}
	if _, err := DecodeUDPPacket(b); err != nil {
		return nil, err
	}
	return buf.MultiBuffer{b}, nil
}

type UDPWriter struct {
	request *protocol.RequestHeader
	writer  io.Writer
}

func NewUDPWriter(request *protocol.RequestHeader, writer io.Writer) *UDPWriter {
	return &UDPWriter{
		request: request,
		writer:  writer,
	}
}

// Write implements io.Writer.
func (w *UDPWriter) Write(b []byte) (int, error) {
	eb, err := EncodeUDPPacket(w.request, b)
	if err != nil {
		return 0, err
	}
	defer eb.Release()
	if _, err := w.writer.Write(eb.Bytes()); err != nil {
		return 0, err
	}
	return len(b), nil
}

func ClientHandshake(request *protocol.RequestHeader, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	fmt.Println("in proxy-socks-protocol.go func ClientHandshake")
	authByte := byte(authNotRequired)
	if request.User != nil {
		authByte = byte(authPassword)
	}

	b := buf.New()
	defer b.Release()

	common.Must2(b.Write([]byte{socks5Version, 0x01, authByte}))
	if authByte == authPassword {
		account := request.User.Account.(*Account)

		common.Must(b.WriteByte(0x01))
		common.Must(b.WriteByte(byte(len(account.Username))))
		common.Must2(b.WriteString(account.Username))
		common.Must(b.WriteByte(byte(len(account.Password))))
		common.Must2(b.WriteString(account.Password))
	}

	if err := buf.WriteAllBytes(writer, b.Bytes()); err != nil {
		return nil, err
	}

	b.Clear()
	if _, err := b.ReadFullFrom(reader, 2); err != nil {
		return nil, err
	}

	if b.Byte(0) != socks5Version {
		return nil, errors.New("unexpected server version: ")
	}
	if b.Byte(1) != authByte {
		return nil, errors.New("auth method not supported.")
	}

	if authByte == authPassword {
		b.Clear()
		if _, err := b.ReadFullFrom(reader, 2); err != nil {
			return nil, err
		}
		if b.Byte(1) != 0x00 {
			return nil, errors.New("server rejects account: ")
		}
	}

	b.Clear()

	command := byte(cmdTCPConnect)
	if request.Command == protocol.RequestCommandUDP {
		command = byte(cmdUDPAssociate)
	}
	common.Must2(b.Write([]byte{socks5Version, command, 0x00 /* reserved */}))
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		return nil, err
	}

	if err := buf.WriteAllBytes(writer, b.Bytes()); err != nil {
		return nil, err
	}

	b.Clear()
	if _, err := b.ReadFullFrom(reader, 3); err != nil {
		return nil, err
	}

	resp := b.Byte(1)
	if resp != 0x00 {
		return nil, errors.New("server rejects request: ")
	}

	b.Clear()

	address, port, err := addrParser.ReadAddressPort(b, reader)
	if err != nil {
		return nil, err
	}

	if request.Command == protocol.RequestCommandUDP {
		udpRequest := &protocol.RequestHeader{
			Version: socks5Version,
			Command: protocol.RequestCommandUDP,
			Address: address,
			Port:    port,
		}
		return udpRequest, nil
	}

	return nil, nil
}
