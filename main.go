package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func passToHash(hash *hash.Hash, filename string) (err error) {
	fmt.Println(filename)
	file, err := os.Open(filename)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		return
	}

	if _, err = io.Copy(*hash, bytes.NewBufferString(file.Name())); err != nil {
		return
	}

	_, err = io.Copy(*hash, file)

	return
}

func main() {
	hashSum := sha256.New()
	files, present := os.LookupEnv("PLUGIN_FILES")
	if !present {
		log.Fatal("PLUGIN_FILES not specified")
	}
	validHash, present := os.LookupEnv("PLUGIN_HASH")
	if !present {
		log.Fatal("PLUGIN_HASH not specified")
	}
	counter := 0
	for _, name := range strings.Split(files, ",") {
		info, err := os.Stat(name)
		if os.IsNotExist(err) {
			log.Fatal("File '" + name + "' does not exist.")
		}
		if info.IsDir() {
			err = filepath.Walk(name, func(path string, info fs.FileInfo, err error) error {
				if !info.IsDir() {
					counter++
					return passToHash(&hashSum, path)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			counter++
			err = passToHash(&hashSum, name)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("-----------------------")
	calculatedHash := fmt.Sprintf("%x", hashSum.Sum(nil))
	if calculatedHash != validHash {
		fmt.Printf("%d items processed\n\n", counter)
		fmt.Printf("Actual: %s\nExpected: %s\n", calculatedHash, validHash)
		fmt.Println("Validation failed")
		os.Exit(1)
	} else {
		fmt.Println("Integrity check passed")
	}
}
