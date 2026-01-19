package main

import (
	"context"
	"log"
	"os"
)

func main() {
	cmd := Commands()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
