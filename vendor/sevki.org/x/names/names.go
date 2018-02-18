package names

import (
	"crypto/sha512"
	"encoding/binary"
	"io"
	"math/rand"
)

// Random returns a raondom name
func Random() string {
	return names["all"][rand.Intn(len(names["all"]))]
}

// For returns a name for a given string by calculating its hash
func For(name string) string {
	h := sha512.New()
	io.WriteString(h, name)
	u, _ := binary.Uvarint(h.Sum(nil))
	nAll := uint64(len(names["all"]))
	return names["all"][int(u%nAll)]
}
