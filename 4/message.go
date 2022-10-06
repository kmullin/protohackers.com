package main

import "strings"

const keySeparator = "="

type msg string

func (m msg) IsInsert() bool {
	return strings.Contains(string(m), keySeparator)
}

func (m msg) KV() (string, string) {
	if !m.IsInsert() {
		return "", ""
	}
	k := strings.SplitN(string(m), keySeparator, 2)
	return k[0], k[1]
}
