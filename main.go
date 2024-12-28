package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ultra-supara/cookie/HackChromeData/browsingdata"
	"github.com/ultra-supara/cookie/HackChromeData/masterkey"
)

func main() {
	// Parse cli options
	targetpath := flag.String("targetpath", "", "File path of the kind (Cookies or Login Data)")
	kind := flag.String("kind", "", "cookie or logindata")
	localState := flag.String("localstate", "", "(optional) Chrome Local State file path")
	sessionstorage := flag.String("sessionstorage", "", "(optional) Chrome Sesssion Storage on Keychain (Mac only)")

	flag.Parse()
	if *targetpath == "" || *kind == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Get Chrome's master key
	var decryptedKey string
	if *sessionstorage == "" {
		// Default path to get master key
		k, err := masterkey.GetMasterKey(*localState)
		if err != nil {
			log.Fatalf("Failed to get master key: %v", err)
		}
		decryptedKey = base64.StdEncoding.EncodeToString(k)
	} else if runtime.GOOS == "darwin" {
		// User input seed key in keychain
		b, err := masterkey.KeyGeneration([]byte(*sessionstorage))
		if err != nil {
			log.Fatalf("Failed to get master key: %v", err)
		}
		decryptedKey = base64.StdEncoding.EncodeToString(b)
	}
	fmt.Println("Master Key: " + decryptedKey)

	// Get Decrypted Data
	log.SetOutput(os.Stderr)
	switch *kind {
	case "cookie":
		c, err := browsingdata.GetCookie(decryptedKey, *targetpath)
		if err != nil {
			log.Fatalf("Failed to get logain data: %v", err)
		}
		output := struct {
			Cookies []browsingdata.Cookie `json:"cookies"`
		}{
			Cookies: c,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			log.Fatalf("Failed to encode cookie data: %v", err)
		}

	case "logindata":
		ld, err := browsingdata.GetLoginData(decryptedKey, *targetpath)
		if err != nil {
			log.Fatalf("Failed to get login data: %v", err)
		}
		for _, v := range ld {
			j, _ := json.Marshal(v)
			fmt.Println(string(j))
		}

	default:
		fmt.Println("Failed to get kind")
		os.Exit(1)
	}
}
