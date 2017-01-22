package util_test

import (
	"fmt"

	"github.com/noypi/util"
)

func ExamplePrivKey() {
	// sender's keys
	privkSender, _ := util.GenPrivKey(2048)
	pubkSender := privkSender.PubKey()

	// receiver's keys
	privkRecvr, _ := util.GenPrivKey(2048)
	pubkRecvr := privkRecvr.PubKey()

	const message = "hello youtube"

	// encrypt message
	bbCipher, _ := pubkRecvr.Encrypt([]byte(message))
	signatureSender, _ := privkSender.Sign(bbCipher)

	// sending [cipher, signature]

	// verify / decrypt message
	bbPlain, _ := privkRecvr.Decrypt(bbCipher)
	bValid := (nil == pubkSender.Verify(bbCipher, signatureSender))

	if bValid {
		fmt.Println(string(bbPlain))
	}

	// Output: hello youtube
}
