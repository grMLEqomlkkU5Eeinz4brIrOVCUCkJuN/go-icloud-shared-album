package utils

import (
	"strconv"
	"strings"
)

type Base62 struct {
	index map[byte]int
}

func NewBase62() *Base62 {
	cs := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	index := make(map[byte]int, len(cs))

	for i := 0; i < len(cs); i++ {
		index[cs[i]] = i
	}

	return &Base62{index: index}
}

func (b *Base62) Decode(s string) int {
	result := 0
	for i := 0; i < len(s); i++ {
		result = result*62 + b.index[s[i]]
	}
	return result
}

// GetBaseUrl constructs the initial iCloud base URL based on the album token.
func GetBaseUrl(token string) string {
	b62 := NewBase62()

	t := token[0]
	n := 0
	if t == 'A' {
		n = b62.Decode(string(token[1]))
	} else {
		n = b62.Decode(token[1:3])
	}

	i := strings.Index(token, ";")
	r := token
	s := ""

	if i >= 0 {
		s = token[i+1:]
		r = strings.Replace(r, ";"+s, "", 1)
	}

	serverPartition := n

	baseUrl := "https://p"

	if serverPartition < 10 {
		baseUrl += "0" + strconv.Itoa(serverPartition)
	} else {
		baseUrl += strconv.Itoa(serverPartition)
	}
	baseUrl += "-sharedstreams.icloud.com"
	baseUrl += "/"
	baseUrl += token // Using original token here as per TS, not 'r'
	baseUrl += "/sharedstreams/"

	return baseUrl
}
