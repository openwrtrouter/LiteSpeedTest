package tools

import (
	"bytes"
	"math/rand"
	"sync"

	"github.com/xxf098/lite-proxy/common/pool"
)

var BufPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}

func AppendRandBytes(b *bytes.Buffer, length int) {
	randBytes := pool.Get(length)
	defer pool.Put(randBytes)
	rand.Read(randBytes)
	b.Write(randBytes)
}
