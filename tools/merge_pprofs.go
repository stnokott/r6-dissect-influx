package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/pprof/profile"
)

func main() {
	rootFolder := flag.String("root", "", "Root folder containing *.pprof files")
	outFile := flag.String("o", "", "Output (merged) file")
	flag.Parse()

	files, err := os.ReadDir(*rootFolder)
	if err != nil {
		log.Fatalf("could not read dir %s: %v", *rootFolder, err)
	}

	var profiles []*profile.Profile
	for _, f := range files {
		full := filepath.Join(*rootFolder, f.Name())
		if !strings.HasSuffix(full, ".pprof") {
			continue
		}
		log.Printf("processing %s...\n", full)
		file, errFiles := os.Open(full)
		if errFiles != nil {
			log.Fatalf("could not open %s: %v", full, errFiles)
		}
		defer func() {
			if errInner := file.Close(); errInner != nil {
				log.Printf("WARNING: error closing %s: %v", file.Name(), errInner)
			}
		}()
		p, errFiles := profile.Parse(file)
		if errFiles != nil {
			log.Fatalf("could not parse %s: %v", full, errFiles)
		}
		profiles = append(profiles, p)
	}

	log.Printf("merging %d files...\n", len(profiles))
	merged, err := profile.Merge(profiles)
	if err != nil {
		log.Fatalf("error merging: %v", err)
	}

	log.Printf("writing merged file to %s...\n", *outFile)
	out, err := os.OpenFile(*outFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("could not open file for writing: %v", err)
	}
	defer func() {
		if errInner := out.Close(); errInner != nil {
			log.Printf("WARNING: error closing %s: %v", out.Name(), errInner)
		}
	}()

	if err = merged.Write(out); err != nil {
		log.Fatalf("error writing: %v", err)
	}
}
