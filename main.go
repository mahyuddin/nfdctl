package main

import (
	"github.com/davecgh/go-spew/spew"
	//"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-ndn/log"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/packet"
	//"github.com/go-ndn/persist"
	"github.com/go-ndn/tlv"
)

var (
	command = flag.String("command", "status", "nfdctl command")
	configPath = flag.String("config", "client.json", "config path")
)

var (
	key ndn.Key
)

func main() {

	flag.Parse()
	
	// connect to nfd
	conn, err := packet.Dial("tcp", ":6363")
	if err != nil {
		log.Fatalln(err)
	}

	// read client key
	pem, err := os.Open("key/default.pri")
	if err != nil {
		log.Fatalln(err)
	}
	defer pem.Close()
	
	key, err = ndn.DecodePrivateKey(pem)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("key", key.Locator())

	// create a new face
	//recv := make(chan *ndn.Interest)
	//face := ndn.NewFace(conn, recv)
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	// create a data fetcher
	f := mux.NewFetcher()

	f.Use(mux.ChecksumVerifier)
	// 2. add the data to the in-memory cache
	f.Use(mux.Cacher)
	// 3. logging
	f.Use(mux.Logger)
	// see producer
	// 4. assemble segments if the content has multiple segments
	// 5. decrypt
	//dec := mux.AESDecryptor([]byte("example key 1234"))

	var faces []ndn.FaceEntry
	tlv.UnmarshalByte(f.Fetch(face,
		&ndn.Interest{
			Name: ndn.NewName("/localhost/nfd/faces/list"),
			Selectors: ndn.Selectors{
				MustBeFresh: true,
			},
		}, mux.Assembler),
		&faces,
		128,
	)

	fmt.Println()
	fmt.Println("Face List")
	fmt.Println("---------")
	spew.Dump(faces)

	var fib []ndn.FIBEntry
	tlv.UnmarshalByte(f.Fetch(face,
		&ndn.Interest{
			Name: ndn.NewName("/localhop/nfd/fib/list"),
			Selectors: ndn.Selectors{
				MustBeFresh: true,
			},
		}, mux.Assembler),
		&fib,
		128,
	)

	fmt.Println()
	fmt.Println("FIB")
	fmt.Println("---")
	spew.Dump(fib)


	var rib []ndn.RIBEntry
	tlv.UnmarshalByte(f.Fetch(face,
		&ndn.Interest{
			Name: ndn.NewName("/localhop/nfd/rib/list"),
			Selectors: ndn.Selectors{
				MustBeFresh: true,
			},
		}, mux.Assembler),
		&rib,
		128,
	)

	fmt.Println()
	fmt.Println("RIB")
	fmt.Println("---")
	spew.Dump(rib)

}