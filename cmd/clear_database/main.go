package main

import (
	db "github.com/replit/database-go"
	"log"
)

func main() {
	keys, err := db.ListKeys("")
	if err != nil {
		panic(err)
	}
	if len(keys) > 0 {
		for _, key := range keys {
			err := db.Delete(key)
			if err != nil {
				log.Fatalf("failed to delete key %s: %s", key, err)
			}
		}
		log.Println("all keys deleted")
	} else {
		log.Println("no keys found")
	}
}
