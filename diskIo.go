package main

import (
	"github.com/petar/GoLLRB/llrb"
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

func ioLessFun(a, b interface{}) bool {
	return a.(*IoArgs).offset < b.(*IoArgs).offset
}

func IoRoutine(request <-chan *IoArgs, responce chan<- interface{}) {
	log.Println("start IoRoutine")

	sortPieces := llrb.New(ioLessFun)

	for arg := range request {
		//todo: sort by offset and batch process io
		cnt := len(request)

		//batch get request, then sort by offset
		sortPieces.InsertNoReplace(arg)
		for i := 0; i < cnt; i++ {
			a := <-request
			sortPieces.InsertNoReplace(a)
			if a.ioMode == MODE_WRITE {
				break
			}
		}

		//batch process
		if cnt > cap(request) / 2 {
			log.Println("io is busy, batch io count", cnt)
		}

		for {
			min := sortPieces.DeleteMin()
			if min == nil {
				break
			}

			//log.Println("io offset", min.(*IoArgs).offset)

			HandleIo(min.(*IoArgs))
			responce <- min.(*IoArgs).context
		}
	}

	log.Println("exit IoRoutine")
}

func HandleIo(arg *IoArgs) {
	var realLength = 0
	var err error = nil
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
	} else {
		if arg.ioMode == MODE_READ {
			if c, ok := arg.context.(*ReadContext); ok {
				c.realLength = c.length
			} else {
				panic("")
			}
		} else { //write
			if c, ok := arg.context.(*WriteContext); ok {
				c.realLength = c.length
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
		log.Printf("warning, disk io too slow, use %v seconds, mod:%v, offset:%v\n", sec, mod, arg.offset)
	}
}
