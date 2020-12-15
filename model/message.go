package model

import (
	"fmt"
	"strings"
)

type Prefix string

func Extract(s string) (Prefix, string) {
	for k, v := range ValidPrefixMap {
		if v {
			if strings.HasPrefix(s, k.String()) {
				return k, s[len(k):]
			}
		}
	}
	return Prefix(s), ""
}

func (p Prefix) AssembleMessage(s string) string {
	return fmt.Sprintf("%s%s\n", p.String(), s)
}

func (m Prefix) String() string {
	return string(m)
}

