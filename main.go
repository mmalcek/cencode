package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

const version = "0.2.0"

func main() {
	removeBom := flag.Bool("rb", false, "Remove BOM at the begining of the file")
	removeBomAll := flag.Bool("rab", false, "Remove all occurrences of BOM")
	removeCR := flag.Bool("rcr", false, "Replace all CRLF with LF")
	addCR := flag.Bool("acr", false, "Replace all LF with CRLF")
	listCodings := flag.Bool("l", false, "List all supported encodings")
	decodeFrom := flag.String("e", "", `Input Encoding (e.g. "ISO 8859-1") - output will be UTF-8. If not specified encoding will not be changed`)
	inFile := flag.String("i", "", "Input file")
	outFile := flag.String("o", "", "Output file")
	flag.Parse()

	// Print flags if no flags provided
	if len(os.Args) < 2 {
		fmt.Printf("Change encoding tool version: %s\n\n", version)
		flag.PrintDefaults()
		return
	}
	// List available encodings
	if *listCodings {
		fmt.Println("Supported encodings:")
		for i := range charmap.All {
			cmap, ok := charmap.All[i].(*charmap.Charmap)
			if ok {
				fmt.Println(cmap.String())
			}
		}
		return
	}

	if *inFile == "" {
		log.Fatal("Input file is required")
	}
	if *outFile == "" {
		log.Fatal("Output file is required")
	}

	f, err := ioutil.ReadFile(*inFile)
	if err != nil {
		log.Fatalln(err)
	}

	if *removeBom {
		if len(f) < 3 { // BOM is 3 bytes so file must be at least 3 bytes
			log.Fatal("Input file is too short")
		}
		if f[0] == 0xEF && f[1] == 0xBB && f[2] == 0xBF {
			f = f[3:]
		}
	}

	if *removeBomAll {
		if len(f) < 3 { // BOM is 3 bytes so file must be at least 3 bytes
			log.Fatal("Input file is too short")
		}
		for i := range f {
			if i == len(f)-3 {
				break
			}
			if f[i] == 0xEF && f[i+1] == 0xBB && f[i+2] == 0xBF {
				f = append(f[:i], f[i+3:]...)
			}
		}
	}

	if *removeCR {
		f = []byte(strings.ReplaceAll(string(f), "\r\n", "\n"))
	}

	if *addCR {
		f = []byte(strings.ReplaceAll(string(f), "\n", "\r\n"))
	}

	if *decodeFrom != "" && decodeFrom != nil {
		inputEnc := ""
		if decodeFrom != nil { // Without this check, compiler complains "possible nil pointer dereference (SA5011)"
			inputEnc = *decodeFrom
		}
		maps := make(map[string]*charmap.Charmap, len(charmap.All))
		if decodeFrom != nil {
			for i := range charmap.All {
				cmap, ok := charmap.All[i].(*charmap.Charmap)
				if ok {
					maps[cmap.String()] = charmap.All[i].(*charmap.Charmap)
				}
			}
		}
		if maps[inputEnc] == nil {
			log.Fatalf("Unknown encoding: %s", inputEnc)
			return
		}
		f, _ = maps[inputEnc].NewDecoder().Bytes(f)
	}
	ioutil.WriteFile(*outFile, f, 0644)
}
