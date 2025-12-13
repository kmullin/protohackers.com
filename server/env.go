package server

import "os"

const defaultAddr = ":8080"

func GetListenAddr() string {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = defaultAddr
	}
	return addr
}
