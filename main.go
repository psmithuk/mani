package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	inputFilename  string
	outputFilename string
	inputfile      io.Reader
	outputfile     *bufio.Writer
	prefix         string
	showversion    bool
)

const VERSION = "0.0.1"

type ManifestEntry struct {
	Url       string `json:"url"`
	Mandatory bool   `json:"mandatory"`
}

type Manifest struct {
	Entries []ManifestEntry `json:"entries"`
}

func NewManifest() (m *Manifest) {
	m = &Manifest{}
	m.Entries = make([]ManifestEntry, 0, 0)
	return m
}

func (m *Manifest) AddEntry(url string) {
	entry := ManifestEntry{url, true}
	m.Entries = append(m.Entries, entry)
}

func (m *Manifest) String() string {
	s := ""
	for _, e := range m.Entries {
		s = fmt.Sprintf("%s%s\n", s, e.Url)
	}
	return s
}

func init() {

	flag.BoolVar(&showversion, "version", false, "print version string")

	flag.StringVar(&inputFilename, "i", "", "input filename (stdout if none provided)")
	flag.StringVar(&outputFilename, "o", "", "output filename (stdout if none provided)")
	flag.StringVar(&prefix, "p", "", "s3 bucketname or prefix")

}

func main() {
	flag.Usage = usage
	flag.Parse()

	if showversion {
		version()
		return
	}

	if inputFilename != "" {
		file, err := os.Open(inputFilename)
		if err != nil {
			exitWithError(fmt.Errorf("Unable to open input %s: %s", inputFilename, err))
		}
		defer file.Close()
		inputfile = bufio.NewReader(file)
	} else {
		inputfile = bufio.NewReader(os.Stdin)
	}

	// output to Stdout if no file given
	if outputFilename != "" {
		file, err := os.Create(outputFilename)
		if err != nil {
			exitWithError(fmt.Errorf("Unable to create output %s: %s", outputFilename, err))
		}
		defer file.Close()
		outputfile = bufio.NewWriter(file)
	} else {
		outputfile = bufio.NewWriter(os.Stdout)
	}
	defer outputfile.Flush()

	manifest := NewManifest()

	scanner := bufio.NewScanner(inputfile)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if prefix != "" {
			line = fmt.Sprintf("%s%s", prefix, line)
		}
		manifest.AddEntry(line)
	}

	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		exitWithError(fmt.Errorf("Unable to create JSON manifest: %s", err))
	}

	fmt.Fprintf(outputfile, "%s", b)
}

// display usage message
func usage() {
	fmt.Fprintf(os.Stderr, "usage: mani [flags]\n")
	flag.PrintDefaults()
}

// display error and exit
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

// print application version
func version() {
	fmt.Printf("v%s\n", VERSION)
}
