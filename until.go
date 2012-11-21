package main

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func peerId() string {
	sid := "-tt" + strconv.Itoa(os.Getpid()) + "_" + strconv.FormatInt(rand.Int63(), 10)
	log.Println("peerId() ", sid[0:20])
	return sid[0:20]
}

var kBitTorrentHeader = []byte{'\x13', 'B', 'i', 't', 'T', 'o', 'r',
	'r', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l'}

func string2Bytes(s string) []byte { return bytes.NewBufferString(s).Bytes() }

func GetLocalIP() (string, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, t := range ifs {
		addrsList, err := t.Addrs()
		if err != nil {
			return "", err
		}
		for _, a := range addrsList {
			/*
			   ipnet, ok := a.(*net.IPNet) 
			   log.Printf("%+v\n", ipnet)
			   if !ok { 
			           continue 
			   } 

			   v4 := ipnet.IP.To4() 
			   log.Printf("%+v\n", ipnet.IP.To4())
			   if v4 == nil || v4[0] == 127 { // loopback address 
			           continue 
			   } 
			   return v4, nil 
			*/
			return a.String(), nil
		}
	}
	return "", errors.New("cannot find local IP address")
}

type CacheValue struct {
	buf []byte
}

func (self *CacheValue) Size() int {
	return len(self.buf)
}

func min(a, b int) int {
	if a > b {
		return a
	}

	return b	
}

type CacheCounter struct{
	hint int64
	miss int64
}

func (c* CacheCounter) IncHint() {
	c.hint++
}

func (c* CacheCounter) IncMiss() {
	c.miss++
}

func (c* CacheCounter) HintRate() float64 {
	return float64(c.hint) / float64(c.hint + c.miss)
}