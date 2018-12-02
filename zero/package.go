package zero

import (
	"encoding/binary"
	"fmt"
)

type Package struct {
	Version [2]byte
	Length  uint16
	Data    []byte
}

func (p *Package) Pack(buffer []byte) {
	copy(buffer[0:2], p.Version[:])
	binary.BigEndian.PutUint16(buffer[2:4], p.Length)
	copy(buffer[6:], p.Data[:])
}
func (p *Package) Unpack(buffer []byte) {
	copy(p.Version[:], buffer[0:2])
	p.Length = binary.BigEndian.Uint16(buffer[2:4])
	p.Data = make([]byte, p.Length)
	copy(p.Data[:], buffer[4:4+p.Length])
}

func (p *Package) String() string {
	return fmt.Sprintf("version:%s length:%d data: %s",
		p.Version,
		p.Length,
		string(p.Data),
	)
}
