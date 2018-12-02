package zero

import (
	"fmt"
	"testing"
)

func TestPackage_Pack(t *testing.T) {

	pack := &Package{}
	pack.Version = [2]byte{'V', '1'}
	pack.Data = []byte("hello")
	pack.Length = uint16(len(pack.Data))

	buffer := make([]byte, pack.Length+4)
	pack.Pack(buffer[:])
	fmt.Println(buffer[0:2])

	upack := &Package{}
	upack.Unpack(buffer)
	fmt.Println(upack.Version, upack.Length)
}
