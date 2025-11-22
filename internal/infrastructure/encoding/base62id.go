package encoding

import (
	"math/big"

	"github.com/jxskiss/base62"
)

func EncodeID(n int64) string {
	u := uint64(n)
	return string(base62.FormatUint(u))
}

func DecodeID(s string) int64 {
	if s == "" {
		return 0
	}

	v, err := base62.ParseUint([]byte(s))
	if err == nil {
		return int64(v)
	}

	b := new(big.Int)
	b.SetString(s, 62)
	if !b.IsUint64() {
		return ^int64(0)
	}
	return int64(b.Uint64())
}
