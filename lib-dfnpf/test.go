package dfnpf

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//TestNoise implements tests for different
// HandshakePattern patterns and their properties
//(at least it should)
func TestNoise() error {
	LogInfo.Printf("Testing %s of %s HandshakePattern of Noise Protocol Framework\n", Test, HandshakePattern)
	rand.Seed(time.Now().UnixNano())
	failed := false

	if len(HandshakePattern) > 1 || Test == "msgtest" {
		// test for 1 random length msg
		if err := testOneMessage(HandshakePattern, rand.Intn(MaxMsgLen)); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for 1 random message is succeceed")
		}

		//tests for special cases: with empty payload and maximum length payload
		if err := testOneMessage(HandshakePattern, 0); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for empty string is succeceed")
		}

		if err := testOneMessage(HandshakePattern, MaxMsgLen); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for maximum message's length is succeceed")
		}

		// test for psk
		if err := testOneMessageForPsk(rand.Intn(MaxMsgLen)); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for 1 random message with psk-modifier is succeceed")
		}

		//special cases for psk

		if err := testOneMessageForPsk(0); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for empty string with psk-modifier is succeceed")
		}

		if err := testOneMessageForPsk(MaxMsgLen); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for maximum message's length with psk-modifier is succeceed")
		}

		//the long ones
		if err := testEachByteMessage(); err != nil {
			failed = true
			LogError.Printf("Error: %v\n", err)
		} else {
			LogSuccess.Println("Test for consiquencies of messages is succeded")
		}

		if err := testPskEachByteMessage(); err != nil {
			failed = true
			LogError.Printf("Error: %v\n", err)
		} else {
			LogSuccess.Println("Test for consiquencies of messages with psk-modifier is succeded")
		}

	}
	if Test == "cktest" && len(HandshakePattern) == 1 {
		if err := testOneMessage(HandshakePattern, 100); err != nil {
			failed = true
			LogError.Printf("Error:%v\n", err)
		} else {
			LogSuccess.Println("Test for 1 random message is succeceed")
		}
	}

	if failed {
		return errors.New("At least one test failed")
	}

	return nil
}

func testOneMessage(HandshakePat string, lenmsg int) error {
	msg := randomHex(lenmsg)
	LogInfo.Printf("Testing %s of len %d of HandshakePattern pattern %s\n", Test, lenmsg, HandshakePat)

	args := initHandsh(HandshakePat)
	args = append(args, msg)

	out1, err1 := execProg(Prog1, Test+HandshakePat, args)
	if err1 != nil {
		LogError.Printf("Error: %v", err1)
		return err1
	}

	out2, err2 := execProg(Prog2, Test+HandshakePat, args)
	if err2 != nil {
		LogError.Printf("Error: %v", err2)
		return err2
	}

	if out1 != out2 {
		LogError.Printf("Test failed, there are differencies in %v and %v:", Prog1, Prog2)
		LogToFile.Printf("\n%s\n%s\n with args %v", out1, out2, args)
	}

	return nil
}

func testEachByteMessage() error {
	msg := randomHex(MaxMsgLen)

	LogInfo.Printf("Testing %s with different message's length of %s\n", Test, HandshakePattern)
	arghandshakePat := initHandsh(HandshakePattern)
	errSlice := testConsistency(HandshakePattern, msg, arghandshakePat)
	return errSlice
}

func testConsistency(HandshakePat string, message string, arghandshakePat []string) error {
	LogInfo.Printf("Testing %s consistency of HandshakePattern pattern %s of len: %d\n", Test, HandshakePat, IterNum)

	msgs := make(chan string)
	errs := make(chan error, IterNum)

	var wgroup sync.WaitGroup
	wgroup.Add(1)
	go func() {
		for msg := range msgs {
			id := Test + "cons" + HandshakePat + strconv.Itoa(len(msg))
			argForRunningProg := append(arghandshakePat, msg)
			out1, err1 := execProg(Prog1, id, argForRunningProg)
			out2, err2 := execProg(Prog2, id, argForRunningProg)

			if out1 != out2 {
				LogError.Printf("Failed to run %s programs: %v and %v on length %d\n", Test, Prog1, Prog2, len(msg))
				LogToFile.Printf(" with args %v \nwhere results are\n%v\n%v", argForRunningProg, out1, out2)
				errs <- fmt.Errorf("error on length %d", len(msg))
				log.Fatalln(err1, err2)
			}
		}
		wgroup.Done()
	}()

	for i := MinMsgLen + 1; i < IterNum; i++ {
		LogToTerm.Printf("\033[1F\033[2K%d / %d", i+1, IterNum)
		LogToFile.Printf("%d / %d ", i+1, IterNum)
		msgs <- message[:i*2]
	}
	close(msgs)
	//wait until wgroups ending and then prepare output if there were any mistakes
	wgroup.Wait()

	if len(errs) > 0 {
		var errSlice MultiError
		for len(errs) > 0 {
			err := <-errs
			errSlice = append(errSlice, err)
		}
		close(errs)
		return errSlice
	}
	return nil
}

func testOneMessageForPsk(msglen int) error {
	var errSlice MultiError
	if len(HandshakePattern) == 1 || HandshakePattern == "NN" || HandshakePattern == "NK" || HandshakePattern == "KN" || HandshakePattern == "KK" {
		if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(0), msglen); err != nil {
			errSlice = append(errSlice, err)
		}
	}
	if HandshakePattern == "X" || HandshakePattern == "IN" || HandshakePattern == "IK" {
		if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(1), msglen); err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if HandshakePattern[0] != 'X' && len(HandshakePattern) == 2 {
		if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(2), msglen); err != nil {
			errSlice = append(errSlice, err)
		}
	}
	if HandshakePattern[0] == 'X' && len(HandshakePattern) == 2 {
		if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(3), msglen); err != nil {
			errSlice = append(errSlice, err)
		}
	}
	//it's not really needed in practice
	if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(0), msglen); err != nil {
		errSlice = append(errSlice, err)
	}
	if err := testOneMessage(HandshakePattern+"psk"+strconv.Itoa(1), msglen); err != nil {
		errSlice = append(errSlice, err)
	}

	if errSlice != nil {
		return errSlice
	}
	return nil
}

func testPskEachByteMessage() error {
	msg := randomHex(MaxMsgLen)
	var errSlice MultiError

	if len(HandshakePattern) == 1 || HandshakePattern == "NN" || HandshakePattern == "NK" || HandshakePattern == "KN" || HandshakePattern == "KK" {
		LogInfo.Println("testing with psk-modifier 0th")
		arghandshakePat := initHandsh(HandshakePattern + "psk" + strconv.Itoa(0))
		if err := testConsistency(HandshakePattern+"psk"+strconv.Itoa(0), msg, arghandshakePat); err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if HandshakePattern == "X" || HandshakePattern == "IN" || HandshakePattern == "IK" {
		LogInfo.Println("testing with psk-modifier 1st")
		arghandshakePat := initHandsh(HandshakePattern + "psk" + strconv.Itoa(1))
		if err := testConsistency(HandshakePattern+"psk"+strconv.Itoa(0), msg, arghandshakePat); err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if HandshakePattern[0] != 'X' && len(HandshakePattern) == 2 {
		LogInfo.Println("testing with psk-modifier 2nd")
		arghandshakePat := initHandsh(HandshakePattern + "psk" + strconv.Itoa(2))
		if err := testConsistency(HandshakePattern+"psk"+strconv.Itoa(0), msg, arghandshakePat); err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if HandshakePattern[0] == 'X' && len(HandshakePattern) == 2 {
		LogInfo.Println("testing with psk-modifier 3rd")
		arghandshakePat := initHandsh(HandshakePattern + "psk" + strconv.Itoa(3))
		if err := testConsistency(HandshakePattern+"psk"+strconv.Itoa(0), msg, arghandshakePat); err != nil {
			errSlice = append(errSlice, err)
		}
	}
	if errSlice != nil {
		return errSlice
	}
	return nil
}
