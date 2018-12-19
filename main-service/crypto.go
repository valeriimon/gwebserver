package mainservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

//secretKey must be 16, 24 or 32 length of bytes
var secretKey string

func initCryptoFuncs() {
	var configMap map[string]interface{}
	configFileStream, err := os.OpenFile("config.json", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer configFileStream.Close()

	jsonDecoder := json.NewDecoder(configFileStream)
	if err = jsonDecoder.Decode(&configMap); err != nil {
		panic(err)
	}
	secretKey = configMap["secretKey"].(string)
}

// Encrypt value with cipher; will be used to get application token
func encryptValueWithCipher(value string) string {
	key, err := hex.DecodeString(secretKey) // get hexadecimal array of bytes by secretKey
	valueBuf := []byte(value)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key) // create new aes block
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(valueBuf)) // returns array of bytes ([0,0, ...])
	iv := ciphertext[:aes.BlockSize]                        // iv must be the same length as block length (16)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil { // get array of bytes with random values ([126, 140, ...])
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	//XORs each byte in the given slice with a byte from the
	// cipher's key stream. Dst and src must overlap entirely or not at all.
	stream.XORKeyStream(ciphertext[aes.BlockSize:], valueBuf)

	return fmt.Sprintf("%x", ciphertext)
}

// Decrypt value with cipher; will be used to decrypt application token (see encryptValueWithCipher func for explanation what is going on here)
func decryptValueWithCipher(value string) string {

	key, err := hex.DecodeString(secretKey)
	if err != nil {
		panic(err)
	}
	ciphertext, err := hex.DecodeString(value)
	if err != nil || len(ciphertext) < aes.BlockSize {
		log.Fatalln("Error on DecodeString: ", err.Error(), "Error on cipherText length: ", len(ciphertext) < aes.BlockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func EncryptValueByHex(value string) string {
	result := hex.EncodeToString([]byte(value))
	return fmt.Sprintf("%s", result)
}

func DecryptValueByHex(value string) string {
	result, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", result)
}
