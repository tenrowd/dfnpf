package dfnpf

import (
	"bytes"
	"math/rand"
	"os/exec"
	"strings"
)

// randomHex generates len*2 random hex string to have random bytes of len length
func randomHex(len int) string {
	charArray := make([]byte, len*2)
	for i := 0; i < len*2; i++ {
		charArray[i] = hexChars[(rand.Uint32())%16]
	}
	return string(charArray)
}

// execProg allows to run the program with specific arguments
func execProg(prog, id string, args []string) (string, error) {
	cmd := exec.Command(prog, args...)

	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Start()
	if err != nil {
		LogError.Fatalln("Couldn't start execution of Cmd ", err)
	}
	if err = cmd.Wait(); err != nil {
		LogToTerm.Printf("Error on %s with %s Program", id, prog)
		LogToFile.Printf("Error on %s with %s Program returned: %s", id, prog, out.String()+outErr.String())
	} else {
		LogToFile.Printf("%s with %s runned successfully, it returned: %s", prog, id, out.String())
	}
	return strings.ToLower(strings.TrimSpace(out.String() + outErr.String())), err

}

//var ciphersuite for random calling
var ciphfunc = map[int]string{
	0: "AESGCM",
	1: "ChaChaPoly",
}

var hashfunc = map[int]string{
	0: "SHA256",
	1: "SHA512",
	2: "BLAKE2b",
	3: "BLAKE2s",
}

// initHandsh sets arguments for initializing handshakestate
func initHandsh(handsh string) []string {
	args := []string{"Noise_" + handsh + "_25519_" + ciphfunc[rand.Intn(2)] + "_" + hashfunc[rand.Intn(4)],
		"ephemInitKey"}

	args = append(args, randomHex(KeyLen))

	if PatternKeys[HandshakePattern].respephem {
		args[1] += "_ephemRespKey"
		args = append(args, randomHex(KeyLen))
	}

	if PatternKeys[HandshakePattern].initstat {
		args[1] += "_staticInitKey"
		args = append(args, randomHex(KeyLen))
	}

	if PatternKeys[HandshakePattern].respstat {
		args[1] += "_staticRespKey"
		args = append(args, randomHex(KeyLen))
	}

	if PatternKeys[HandshakePattern].initremotestat {
		args[1] += "_remoteInitStatKey"
		args = append(args, randomHex(KeyLen))
	}

	if PatternKeys[HandshakePattern].respremotestat {
		args[1] += "_remoteRespStatKey"
		args = append(args, randomHex(KeyLen))
	}

	if rand.Intn(2) > 0 {
		args[1] += "_prologue"
		args = append(args, randomHex(KeyLen))
	}
	if len(strings.Split(handsh, "psk")) == 2 {
		args[1] += "_presharedKey"
		args = append(args, randomHex(KeyLen))
	}

	return args
}
