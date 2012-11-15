package main

import "log"
import "strconv"

// BitTorrent message types. Sources:
// http://bittorrent.org/beps/bep_0003.html
// http://wiki.theory.org/BitTorrentSpecification
const (
	CHOKE = iota
	UNCHOKE
	INTERESTED
	NOT_INTERESTED
	HAVE
	BITFIELD
	REQUEST
	PIECE
	CANCEL
	PORT // Not implemented. For DHT support.
)

const STORAGE_BLOCK_SIZE = 32 * 1024
const BLOCK_META_SIZE = 24

func chooseListenPort() (listenPort int, err error) {
	listenPort = cfg.port
	if cfg.useUPnP {
		log.Println("Using UPnP to open port.")
		// TODO: Look for ports currently in use. Handle collisions.
		var nat NAT
		nat, err = Discover()
		if err != nil {
			log.Println("Unable to discover NAT:", err)
			return
		}
		// TODO: Check if the port is already mapped by someone else.
		err2 := nat.DeletePortMapping("TCP", listenPort)
		if err2 != nil {
			log.Println("Unable to delete port mapping", err2)
		}
		err = nat.AddPortMapping("TCP", listenPort, listenPort,
			"Taipei-Torrent port "+strconv.Itoa(listenPort), 0)
		if err != nil {
			log.Println("Unable to forward listen port", err)
			return
		}
	}
	return
}
