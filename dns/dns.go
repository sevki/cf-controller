package dns

import (
	"crypto/sha512"
	"io"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Record cloudflare.DNSRecord

func (r Record) Hash() []byte {
	h := sha512.New()
	io.WriteString(h, r.Content)
	io.WriteString(h, r.Name)
	return h.Sum(nil)
}
