package iniths

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"

	"strconv"
	"strings"

	noise "github.com/flynn/noise"
)

var cipherfunction = map[string]noise.CipherFunc{
	"AESGCM":     noise.CipherAESGCM,
	"ChaChaPoly": noise.CipherChaChaPoly,
}

var hashfunction = map[string]noise.HashFunc{
	"SHA256":  noise.HashSHA256,
	"SHA512":  noise.HashSHA512,
	"BLAKE2b": noise.HashBLAKE2b,
	"BLAKE2s": noise.HashBLAKE2s,
}

var handshakepatterns = map[string]noise.HandshakePattern{
	"N":          noise.HandshakeN,
	"K":          noise.HandshakeK,
	"X":          noise.HandshakeX,
	"NN":         noise.HandshakeNN,
	"KN":         noise.HandshakeKN,
	"NK":         noise.HandshakeNK,
	"KK":         noise.HandshakeKK,
	"NX":         noise.HandshakeNX,
	"KX":         noise.HandshakeKX,
	"XN":         noise.HandshakeXN,
	"IN":         noise.HandshakeIN,
	"XK":         noise.HandshakeXK,
	"IK":         noise.HandshakeIK,
	"XX":         noise.HandshakeXX,
	"XXfallback": noise.HandshakeXXfallback,
	"IX":         noise.HandshakeIX,
}

func bytefromhex(hexstr string) []byte {
	res, err := hex.DecodeString(hexstr)
	if err != nil {
		panic(err)
	}
	return res
}

func hexReader(s string) io.Reader {
	return bytes.NewBuffer(bytefromhex(s))

}

//InitHandshake initializes Handshake state from command line
func InitHandshake(args []string) (*noise.HandshakeState, *noise.HandshakeState, []byte) {

	var prologue, presharedkey, initpeerstatic, respeerstatic []byte
	var presharedkeyplacement int
	var patternName string
	var ephemInit, ephemResp io.Reader
	var initstatickey, initephemeralkey, respstatickey, respephemeralkey noise.DHKey

	handshakeComponents := strings.Split(args[0], "_")

	handshakePatternComponents := strings.Split(handshakeComponents[1], "psk")
	if len(handshakePatternComponents) == 2 {
		presharedkeyplacement, _ = strconv.Atoi(handshakePatternComponents[1])
	}
	patternName = handshakePatternComponents[0]
	if strings.HasSuffix(patternName, "+") {
		patternName = patternName[:len(patternName)-1]
	}

	handshakeparams := strings.Split(args[1], "_")

	for i, value := range handshakeparams {
		switch value {
		case "staticInitKey":
			statickey, err1 := noise.DH25519.GenerateKeypair(hexReader(args[i+2]))
			if err1 != nil {
				log.Fatal(err1)
			}
			initstatickey = statickey
		case "staticRespKey":
			statickey, err2 := noise.DH25519.GenerateKeypair(hexReader(args[i+2]))
			if err2 != nil {
				log.Fatal(err2)
			}
			respstatickey = statickey

		case "remoteInitStatKey":
			initpeerstatic = respstatickey.Public
		case "remoteRespStatKey":

			respeerstatic = initstatickey.Public

		case "ephemInitKey":
			ephemInit = hexReader(args[i+2])
		case "ephemRespKey":
			ephemResp = hexReader(args[i+2])

		case "prologue":
			prologue = bytefromhex(args[i+2])
		case "presharedKey":
			presharedkey = bytefromhex(args[i+2])
		}

	}

	var payload []byte
	if len(handshakeparams) < len(args)-2 {
		payload = bytefromhex(args[len(args)-1])
	}

	ciphersuite := noise.NewCipherSuite(noise.DH25519, cipherfunction[handshakeComponents[3]], hashfunction[handshakeComponents[4]])

	handshInit, _ := initializationHS(patternName, ciphersuite, ephemInit, true,
		prologue, presharedkey, presharedkeyplacement, initstatickey,
		initephemeralkey, initpeerstatic)
	handshResp, _ := initializationHS(patternName, ciphersuite, ephemResp, false,
		prologue, presharedkey, presharedkeyplacement, respstatickey,
		respephemeralkey, respeerstatic)

	return handshInit, handshResp, payload
}

func initializationHS(handshpat string, cs noise.CipherSuite, rng io.Reader,
	initiatorb bool, prologue []byte, preshk []byte, preshkpos int,
	statkeyp noise.DHKey, ephkeyp noise.DHKey,
	peerstat []byte) (*noise.HandshakeState, error) {

	handshakepattern := handshakepatterns[handshpat]
	handshst, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:           cs,
		Random:                rng,
		Pattern:               handshakepattern,
		Initiator:             initiatorb,
		Prologue:              prologue,
		PresharedKey:          preshk,
		PresharedKeyPlacement: preshkpos,
		StaticKeypair:         statkeyp,
		EphemeralKeypair:      ephkeyp,
		PeerStatic:            peerstat,
	})
	return handshst, err
}
