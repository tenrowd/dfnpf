package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	dfnpf "dfnpf/lib-dfnpf"
)

// usage is called when input doesn't seem right or help is needed
func usage() {
	flag.Usage()
	fmt.Println("To perform the tests: \ndfnpf HandshakePattern (or all) (msgtest or cktest) path/to/program1 path/to/program2")
}

func init() {
	flag.Parse()

	// check args, their correctness and existing
	if len(flag.Args()) != 4 {
		usage()
		os.Exit(1)
	}

	if flag.Arg(0) == "all" {
		dfnpf.HandshakePattern = "all"
	} else {
		for _, handshPat := range dfnpf.HandshakePatternsList {
			if flag.Arg(0) == handshPat {
				dfnpf.HandshakePattern = handshPat
				break
			}
		}
	}

	if dfnpf.HandshakePattern == "" {
		log.Fatalln("Invalid interface")
	}

	// check correctness of called test and programs' paths
	if flag.Arg(1) == "msgtest" || flag.Arg(1) == "cktest" {
		dfnpf.Test = flag.Arg(1)
	} else {
		log.Fatalln("This test isn't implemented")
	}
	dfnpf.Prog1 = flag.Arg(2)
	dfnpf.Prog2 = flag.Arg(3)
	if _, err := os.Stat(dfnpf.Prog1); os.IsNotExist(err) {
		log.Fatalln("This file doesn't exist: ", dfnpf.Prog1)
	}
	if _, err := os.Stat(dfnpf.Prog2); os.IsNotExist(err) {
		log.Fatalln("This file doesn't exist:", dfnpf.Prog2)
	}
}

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	// close logFile on exit checking for its error to ensure everything get written.
	defer func() {
		if err := logFile.Close(); err != nil {
			panic(err)
		}
	}()

	dfnpf.InitLog(logFile)
	dfnpf.InitKeys()
	dfnpf.LogInfo.Println("Running dfnpf:")

	// run tests for provided handshake. set number of iterations in
	// test with each-byte sequences of message
	if dfnpf.HandshakePattern != "all" {
		dfnpf.IterNum = 100
		err = dfnpf.TestNoise()
		if err != nil {
			dfnpf.LogError.Println(err)

		} else {
			dfnpf.LogSuccess.Printf("%s test completed without error for %s\n", dfnpf.HandshakePattern, dfnpf.Test)
		}
	} else {
		for _, handshake := range dfnpf.HandshakePatternsList {
			dfnpf.HandshakePattern = handshake
			dfnpf.IterNum = 100
			err = dfnpf.TestNoise()
			if err != nil {
				dfnpf.LogError.Printf("%s test for %s with err %v", dfnpf.Test, dfnpf.HandshakePattern, err)
			} else {
				dfnpf.LogSuccess.Printf("%s test completed without error for %s", dfnpf.Test, dfnpf.HandshakePattern)
			}
		}

	}

	dfnpf.LogInfo.Println("Test are done")
}
