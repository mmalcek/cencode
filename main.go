package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	removeBom := flag.Bool("r", false, "Remove BOM at the begining of the file")
	removeBomAll := flag.Bool("rab", false, "Remove all occurrences of BOM")
	listCodings := flag.Bool("l", false, "List all supported encodings")
	decodeFrom := flag.String("d", "", `Decode to UTF-8. Parameter should contain encoding name (e.g. "ISO 8859-1")`)
	inFile := flag.String("i", "", "Input file")
	outFile := flag.String("o", "", "Output file")
	flag.Parse()

	// Print flags if no flags provided
	if len(os.Args) < 2 {
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
