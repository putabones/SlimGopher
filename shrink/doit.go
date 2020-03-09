package main

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"os"
	"strings"
)

// parse arguments
func args() (*string, *string) {
	//
	// cCreate parser
	parse := argparse.NewParser("shrink", "Shrink compresses and encrypts your embedded, and also creates the variables need for your payload")

	// argument(s)
	p := parse.String("p", "path", &argparse.Options{
		Required: true,
		Validate: nil,
		Help:     "Path where you want the variable files saved (REQUIRED)",
		Default:  "/tmp",
	})

	e := parse.String("e", "embed", &argparse.Options{
		Required: true,
		Validate: nil,
		Help:     "Executable you want to embed (REQUIRED)",
		Default:  nil,
	})

	// catch errors
	err := parse.Parse(os.Args)
	if err != nil{
		fmt.Println(parse.Usage(os.Stdout))
		return nil, nil
	} else {
		// return argument(s)
		return p, e
	}
}

// compression function
func compress(e *string) []byte {
	// byte buffer
	b := bytes.Buffer{}

	// compression writer
	w := zlib.NewWriter(&b)

	// read the bytes from exeecutable
	r, err := ioutil.ReadFile(*e)
	if err != nil{
		fmt.Println("[-] ioutil.ReadFile() Error: ", err)
	}

	_, err = w.Write(r)
	if err != nil{
		fmt.Println("[-] w.Write() Error: ", err)
	}

	// close writer
	err = w.Close()
	if err != nil{
		fmt.Println("[-] w.Close() Error: ", err)
	} else {
		fmt.Println("[+] compressed []byte len: ", b.Len(), "\n[+] compression done...")
	}

	return b.Bytes()
}

func encrypt(data []byte, k []byte, i []byte) []byte {
	b, err := aes.NewCipher(k)
	if err != nil{
		fmt.Println("[-] aes.NewCipher() Error: ", err)
	}

	// fun...
	ase := cipher.NewCFBEncrypter(b, i)
	ciphertext := make([]byte, aes.BlockSize+len(data))

	ase.XORKeyStream(ciphertext, data)

	fmt.Println("[+] encrypted []byte len: ", len(ciphertext), "\n[+] encryption done...")

	return ciphertext
}

// generate the []byte of the executable
func generate(data []byte, path *string, name string) {
	// create []string for writing
	var dataSlice []string
	var fullpath = *path + "/" + name + ".go"

	// create variable file
	file, err := os.Create(fullpath)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("[+] Creating: ", fullpath)
	}
	defer file.Close()

	// write in package and data var
	_, err = file.Write([]byte("package main\n\nvar (\n\t" + name + " = []byte{"))
	if err != nil{
		fmt.Println("[-] file.Write(package) Error: ", err)
	} else {
		fmt.Println("[+] Writing variable name")
	}

	// append the bytes into the []string
	for _, b := range data {
		bString := fmt.Sprintf("%v", b)
		dataSlice = append(dataSlice, bString)
	}

	// write literal bytes
	_, err = file.WriteString(strings.Join(dataSlice, ", "))
	if err != nil{
		fmt.Println("[-] file.WriteString(dataString) Error: ", err)
	} else {
		fmt.Println("[+] Data written to variable")
	}

	// close variable
	_, err = file.Write([]byte("}\n)"))
	if err != nil{
		fmt.Println("[-] file.Write([]byte(\"}\n)\")) Error: ", err)
	} else {
		fmt.Println("[+] Variable written to file")
	}
}

// create encryption key and iv
func keyiv() ([] byte, []byte) {
	// create blank []byte for key and iv
	key := make([]byte, 32)
	iv := make([]byte, 16)

	// read in random bytes to key
	_, err := rand.Read(key)
	if err != nil{
		fmt.Println("[-] rand.Read(key) Error: ", err)
	} else {
		fmt.Println("[+] Random Key: ", key)
	}

	// read in random bytes to iv
	_, err = rand.Read(iv)
	if err != nil{
		fmt.Println("[-] rand.Read(iv) Error: ", err)
	} else {
		fmt.Println("[+] Random IV: ", iv)
	}

	return key, iv
}

// main function
func main(){
	// argparse
	path, exe := args()

	// check for args
	if path != nil && exe != nil {

		key, iv := keyiv()
		compdata := compress(exe)

		// create key and iv

		// compress the exe

		// encrypt the exe
		encdata := encrypt(compdata, key, iv)

		// generate the variable files
		generate(encdata, path, "data")
		generate(key, path, "key")
		generate(iv, path, "iv")
		fmt.Println("[+] Command below to change directory and list\n\ncd", *path, "&& ls -l")
	}
}