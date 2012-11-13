package main

import (
	"log"
	"time"
)

const (
	MODE_READ  = 1
	MODE_WRITE = 2
)

type IoArgs struct {
	f       FileStore
	ioMode  int
	buf     []byte
	offset  int64
	context interface{}
}

type WriteContext struct {
	peer       *peerState
	whichPiece uint32 //piece index
	begin      uint32
	length     int
	realLength int
}

type ReadContext struct {
	peer         *peerState
	msgBuf       []byte
	globalOffset int64
	length       int
	realLength   int
}

func IoRoutine(request <-chan *IoArgs, responce chan<- interface{}) {
	log.Println("start IoRoutine")
	var realLength = 0
	var err error = nil
	for arg := range request {
		start := time.Now()
		if cfg.doRealReadWrite {
			if arg.ioMode == MODE_READ {
				//log.Println("read offset", arg.offset, "bufffer size", len(arg.buf))
				realLength, err = arg.f.ReadAt(arg.buf, arg.offset)
				if err != nil || realLength < 0 {
					panic("")
				}

				if c, ok := arg.context.(*ReadContext); ok {
					c.realLength = realLength
					//log.Println("read", c.realLength, "offset", arg.offset, "bufffer size", len(arg.buf))
				} else {
					panic("")
				}
			} else { //write
				realLength, err = arg.f.WriteAt(arg.buf, arg.offset)
				if err != nil || realLength < 0 {
					panic("")
				}

				if c, ok := arg.context.(*WriteContext); ok {
					c.realLength = realLength
					//log.Println("write", c.realLength, "offset", arg.offset)
				} else {
					panic("")
				}
			}
		}

		sec := time.Now().Sub(start).Seconds()
		if sec >= 2 {
			var mod = "READ"
			if arg.ioMode == MODE_WRITE {
				mod = "WRITE"
			}
			log.Printf("\nwarning, disk io too slow, use %v seconds, mod:%v, offset:%v\n", sec, mod, arg.offset)
		}

		responce <- arg.context
	}
	log.Println("exit IoRoutine")
}
