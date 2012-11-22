package main

import "flag"
import "math/rand"
import "time"

type Config struct {
	STANDARD_BLOCK_LENGTH      int
	MAX_NUM_PEERS              int
	TARGET_NUM_PEERS           int
	MAX_DOWNLOADING_CONNECTION int
	MAX_UPLOADING_CONNECTION   int
	MAX_OUR_REQUESTS           int
	MAX_PEER_REQUESTS          int
	port                       int
	useUPnP                    bool
	fileDir                    string
	useDHT                     bool
	trackerLessMode            bool
	noCheckSum                 bool
	doRealReadWrite            bool
	rechokeTick                int
	totalTransferSize          int64
	changeProtocolName         bool
	superSeeding               bool
}

var cfg Config

func init() {
	cfg = Config{MAX_NUM_PEERS: 200, TARGET_NUM_PEERS: 3,
		MAX_PEER_REQUESTS: 10,
	}

	flag.StringVar(&cfg.fileDir, "fileDir", ".", "path to directory where files are stored")
	// If the port is 0, picks up a random port - but the DHT will keep
	// running on port 0 because ListenUDP doesn't do that.
	// Don't use port 6881 which blacklisted by some trackers.
	flag.IntVar(&cfg.port, "port", 7777, "Port to listen on.")
	flag.BoolVar(&cfg.useUPnP, "useUPnP", false, "Use UPnP to open port in firewall.")
	flag.BoolVar(&cfg.useDHT, "useDHT", false, "Use DHT to get peers.")
	flag.BoolVar(&cfg.trackerLessMode, "trackerLessMode", false, "Do not get peers from the tracker. Good for "+
		"testing the DHT mode.")
	flag.BoolVar(&cfg.noCheckSum, "nochecksum", false, "do not use checksum for fast starting")
	rand.Seed(int64(time.Now().Nanosecond()))
	flag.BoolVar(&cfg.doRealReadWrite, "doRealReadWrite", true, "do not io disk, using memory instead")
	flag.IntVar(&cfg.STANDARD_BLOCK_LENGTH, "STANDARD_BLOCK_LENGTH", STORAGE_BLOCK_SIZE+BLOCK_META_SIZE, "stand block length")
	flag.IntVar(&cfg.MAX_OUR_REQUESTS, "MAX_OUR_REQUESTS", 8, "max our requests")
	flag.IntVar(&cfg.MAX_UPLOADING_CONNECTION, "MAX_UPLOADING_CONNECTION", 3, "max uploading connection")
	flag.IntVar(&cfg.MAX_DOWNLOADING_CONNECTION, "MAX_DOWNLOADING_CONNECTION", 3, "max downloading connection")
	flag.IntVar(&cfg.rechokeTick, "rechokeTick", 10, "rechoke tick seconds")
	flag.Int64Var(&cfg.totalTransferSize, "totalTransferSize", 1*1024*1024*1024, "total transfer size")
	flag.BoolVar(&cfg.changeProtocolName, "changeProtocolName", true, "change protocol name")
	flag.BoolVar(&cfg.superSeeding, "superSeeding", false, "as super seeding")

	if cfg.changeProtocolName {
		kBitTorrentHeader = []byte{'\x13', 'B', 'c', 't', 'T', 'o', 'l', 'l', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l'}
	}
}
