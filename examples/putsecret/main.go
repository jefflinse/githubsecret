package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jefflinse/githubsecret"
)

type publicKey struct {
	Key   string `json:"key"`
	KeyID string `json:"key_id"`
}

func main() {
	dryRun := flag.Bool("dry", false, "dry run (encrypt and print only)")
	flag.Parse()

	if len(flag.Args()) < 3 {
		fmt.Printf("usage: %s owner repo secret-id secret-value", os.Args[0])
		os.Exit(1)
	}

	owner := flag.Args()[0]
	repo := flag.Args()[1]
	secretID := flag.Args()[2]
	secretValue := flag.Args()[3]

	gh := github{
		username: os.Getenv("GITHUB_USERNAME"),
		token:    os.Getenv("GITHUB_TOKEN"),
	}

	key, err := gh.getPublicKey(owner, repo)
	if err != nil {
		log.Fatalf("couldn't obtain public key: %v", err)
	}

	if *dryRun {
		encrypted, err := githubsecret.Encrypt(key.Key, secretValue)
		if err != nil {
			log.Fatalf("couldn't encrypt secret: %v", err)
		}
		fmt.Println(encrypted)
	} else {
		if err := gh.storeSecret(owner, repo, key, secretID, secretValue); err != nil {
			log.Fatalf("couldn't store secret: %v", err)
		}
	}
}
