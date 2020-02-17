package main

import "strconv"

// SrvConfig that packs all info needed to build a server
type SrvConfig struct {
	Host string
	Port int
}

// Addr returns server's host and port, separated by ":"
func (sc *SrvConfig) Addr() string {
	return sc.Host + ":" + strconv.Itoa(sc.Port)
}
