package dfnpf

import (
	"fmt"
	"io"
	"log"
	"os"
)

// there's a constant for providing characters to hex-view
const (
	hexChars  = "abcdef0123456789"
	MinMsgLen = 0
	MaxMsgLen = 65535
	KeyLen    = 32
)

// there are some different global variables
var (
	LogInfo          *log.Logger // info about execution
	LogSuccess       *log.Logger // if test is succeded
	LogError         *log.Logger // if something is wrong or unexpected
	LogToFile        *log.Logger // info is written to the log file to avoid verbose output
	LogToTerm        *log.Logger //printing to terminal
	HandshakePattern string      // using Handshake Pattern
	Test             string      // msgtest or cktest for providing tests
	IterNum          int         // number of iterations in one of the tests
	Prog1            string      // path to the first executable program
	Prog2            string      // path to the second executable program
)

// MultiError allows to store multiple errors
type MultiError []error

// (MultiError) Error implements the error interface for our MultiError type
func (m MultiError) Error() string {
	var str string
	var num uint
	for _, err := range m {
		if err != nil {
			str += "\n" + err.Error()
			num++
		}
	}
	return fmt.Sprintf("(%d errors)%s", num, str)
}

//PatternKeyInfo contains information about keys
type PatternKeyInfo struct {
	initstat, respstat, respremotestat, initremotestat, respephem bool
}

//HandshakePatternsList is a list of handshake patternsList
var HandshakePatternsList = []string{"NN", "KN", "NK", "KK", "NX", "KX", "XN", "IN", "XK", "IK", "XX", "IX", "N", "K", "X"}

//PatternKeys contains information about keys for each Handshake Pattern
var PatternKeys = make(map[string]PatternKeyInfo)

// InitKeys initializes information about keys in each pattern
func InitKeys() {
	for _, h := range HandshakePatternsList {
		var k PatternKeyInfo
		if len(h) == 1 {
			switch h {
			case "N":
				k.respstat = true
				k.initremotestat = true
			case "K":
				k.initstat = true
				k.respremotestat = true
				k.respstat = true
				k.initremotestat = true
			case "X":
				k.initstat = true
				k.respstat = true
				k.initremotestat = true
			}
		} else {
			k.respephem = true
			switch h[0] {
			case 'X', 'I':
				k.initstat = true
			case 'K':
				k.initstat = true
				k.respremotestat = true
			}
			switch h[1] {
			case 'K':
				k.respstat = true
				k.initremotestat = true
			case 'X', 'R':
				k.respstat = true
			}
		}
		PatternKeys[h] = k
	}
}

// InitLog initializes the logging system to output data to the specified file
func InitLog(logFile *os.File) {
	if logFile != nil {
		multiOut := io.MultiWriter(logFile, os.Stdout)
		multiErr := io.MultiWriter(logFile, os.Stderr)

		LogInfo = log.New(multiOut, "\033[mINFO:\033[0m ", 0)
		LogSuccess = log.New(multiOut, "\033[0;32mSUCCESS:\033[0m", 0)
		LogError = log.New(multiErr, "\033[0;31mERROR:\033[0m", log.Lshortfile)
		LogToFile = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		LogToTerm = log.New(os.Stdout, "", 0)

		log.SetOutput(multiOut)
		LogToFile.Println("Initializing logs: done")
	}
}
