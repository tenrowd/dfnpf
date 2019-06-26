package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"dfnpf/examples/iniths"
)

func main() {

	flag.Parse()

	if len(flag.Args()) < 4 || len(flag.Args()) > 11 {
		log.Fatalln("Please provide Noise Protocol name, initiator's static and ephemeral keys,",
			"\n responder's static and ephemeral keys, remote key, prologue, preshared key, message")
	}

	handshInit, handshResp, payload := iniths.InitHandshake(flag.Args())

	var err error
	var msg []byte

	//var csWrite0, csWrite1, csRead0, csRead1 *noise.CipherState
	msg, _, _, _ = handshInit.WriteMessage(nil, payload)
	_, _, _, err = handshResp.ReadMessage(nil, msg)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", hex.EncodeToString(msg))
}
