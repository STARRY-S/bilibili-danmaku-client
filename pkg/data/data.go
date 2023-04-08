package data

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Protocol uint16
type Operation uint32

const (
	PV_0 Protocol = iota // JSON plantext
	PV_1                 // Int 32 Big Endian 人气值
	PV_2                 // Buffer, gzip compressed
	PV_3                 // Buffer brotli compressed
)

const (
	O_2 = Operation(2) // heart beat
	O_3 = Operation(3) // heart beat response (人气值)
	O_5 = Operation(5) // broadcast
	O_7 = Operation(7) // enter into the room
	O_8 = Operation(8) // entrance reply
)

var (
	PkgHeaderLen  uint32 = 16
	HeartBeatData        = []byte("[Object object]")
)

type Package struct {
	PackageLength   uint32    // 0-3,   4 bytes, package length (16 + data len)
	HeaderLength    uint16    // 4-5,   2 bytes, header length, always 16
	ProtocolVersion Protocol  // 6-7,   2 bytes, protocol version
	Operation       Operation // 8-11,  4 bytes, operation
	SequenceID      uint32    // 12-15, 4 bytes, always 1
	Data            []byte
}

func NewPackage(d []byte, p Protocol, o Operation) *Package {
	pkg := &Package{
		PackageLength:   uint32(16 + len(d)),
		HeaderLength:    uint16(16),
		ProtocolVersion: p,
		Operation:       o,
		SequenceID:      1,
		Data:            d,
	}

	return pkg
}

func (p *Package) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, p.PackageLength))
	binary.Write(buf, binary.BigEndian, p.PackageLength)
	binary.Write(buf, binary.BigEndian, p.HeaderLength)
	binary.Write(buf, binary.BigEndian, p.ProtocolVersion)
	binary.Write(buf, binary.BigEndian, p.Operation)
	binary.Write(buf, binary.BigEndian, p.SequenceID)
	binary.Write(buf, binary.BigEndian, p.Data)
	return buf.Bytes()
}

func (p *Package) DecodeHead(d []byte) error {
	if len(d) < 16 {
		return fmt.Errorf("invalid data length")
	}

	buf := bytes.NewReader(d)
	binary.Read(buf, binary.BigEndian, &p.PackageLength)
	if p.PackageLength < 16 {
		return fmt.Errorf("package length in data is less than 16")
	}
	binary.Read(buf, binary.BigEndian, &p.HeaderLength)
	if p.HeaderLength != 16 {
		return fmt.Errorf("header length in data is not 16")
	}
	binary.Read(buf, binary.BigEndian, &p.ProtocolVersion)
	binary.Read(buf, binary.BigEndian, &p.Operation)
	binary.Read(buf, binary.BigEndian, &p.SequenceID)

	return nil
}
