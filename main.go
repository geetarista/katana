package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	Version = "0.1.0"
)

var (
	dieSize      int64 = 6
	scanChan     chan bool
	wordCount    *int
	wordLanguage *string
	wordMap      map[int]string
)

func init() {
	help := flag.Bool("h", false, "show usage for katana")
	wordLanguage = flag.String("l", "en", "choose language to use")
	version := flag.Bool("v", false, "prints current  version")
	wordCount = flag.Int("w", 5, "word count")

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	scanChan = make(chan bool)
	go scanWords(scanChan)
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	flag.PrintDefaults()
}

func roll() (i int, err error) {
	var b *big.Int
	b, err = rand.Int(rand.Reader, big.NewInt(dieSize))
	if err != nil {
		return 0, err
	}

	// Since rand.Int always includes 0, we have to increment
	i = int(b.Uint64()) + 1

	return
}

// roll n dice
func chuck() (n int, err error) {
	rolls := make([]string, *wordCount)

	for range rolls {
		n, err = roll()
		if err != nil {
			return
		}
		rolls = append(rolls, strconv.Itoa(n))
	}

	return strconv.Atoi(strings.Join(rolls, ""))
}

// build a dictionary of the words
func scanWords(done chan bool) {
	wordMap = make(map[int]string)
	file, err := os.Open(fmt.Sprintf("wordlists/%s.asc", *wordLanguage))
	if err != nil {
		fmt.Printf("[ERR] %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fields []string
	var text string
	var line int
	for scanner.Scan() {
		line++

		// Skip PGP content
		if line < 3 || line > 7778 {
			continue
		}

		text = scanner.Text()
		fields = strings.Fields(text)
		i, _ := strconv.Atoi(fields[0])
		wordMap[i] = fields[1]
	}

	done <- true
}

func main() {
	rolls := []int{}
	words := []string{}

	var r int
	var err error

	for i := 0; i < *wordCount; i++ {
		r, err = chuck()
		if err != nil {
			fmt.Println("[ERR] " + err.Error())
		}
		rolls = append(rolls, r)
	}

	<-scanChan

	for _, roll := range rolls {
		word, ok := wordMap[roll]
		if !ok {
			fmt.Errorf("Unable to find word for index %d", roll)
		}
		words = append(words, word)
	}

	fmt.Println(strings.Join(words, " "))
}
