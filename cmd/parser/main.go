package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rejlersembriq/gs2"
	"log"
	"os"
	"time"
)

var filename = flag.String("file", "", "File to parse")

func main() {
	flag.Parse()

	if *filename == "" {
		log.Fatal("Need to speficy filename with -file=<filename>")
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	start := time.Now()

	result, err := gs2.NewDecoder(file).Decode()
	if err != nil {
		log.Fatal(err)
	}

	took := time.Since(start)

	indent, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println(string(indent))

	fmt.Printf("Parsing took: %v\n", took)
}
