package main

import (
	"fmt"

	"github.com/witwoywhy/go-cores/cryptos"
)

func main() {
	key := "12345678901234561234567890123456"
	iv := "1234567812345678"

	data := []interface{}{
		"HELLO",
		map[string]interface{}{
			"HELLO": "WORLD",
		},
		123456,
		3.14,
	}

	aes := cryptos.NewAesCompatibleCryptoJS(key, []byte(iv))
	for i, v := range data {
		enc, err := aes.Encrypt(v)
		if err != nil {
			fmt.Printf("%d: %v encrypt failed: %v\n", i, v, err)
			continue
		}

		fmt.Printf("%d: %v\n", i, enc)

		dec, err := aes.Decrypt(enc)
		if err != nil {
			fmt.Printf("%d: %v decrypt failed: %v\n", i, v, err)
			continue
		}

		fmt.Printf("%d: %v\n", i, dec)
	}
}
