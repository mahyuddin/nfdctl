package nfdctl

import (
	"os"
	"time"

	"github.com/go-ndn/log"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	"github.com/go-ndn/persist"
)

var (
	command = flag.String("command", "status", "nfdctl command")
)

var (
	key ndn.Key
)

func main() {
	
	// connect to nfd
	conn, err := packet.Dial("tcp", ":6363")
	if err != nil {
		log.Fatalln(err)
	}

	// create a new face
	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	// read producer key
	pem, err := os.Open("key/default.pri")
	if err != nil {
		log.Fatalln(err)
	}
	defer pem.Close()
	key, _ := ndn.DecodePrivateKey(pem)

}