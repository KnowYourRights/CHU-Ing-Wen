package converter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	termFilePath     = "./data/term.txt"
	charFilePath     = "./data/char.txt"
	anthologyDirPath = "./data/anthology/"
	errString        = "豪像粗問題ㄌ，等等在４４看！"
)

var termMap map[string]string
var charMap map[string]string
var anthologyFilenames []string

func init() {
	termMap, charMap = make(map[string]string), make(map[string]string)

	// Terms.
	f1, err := os.Open(termFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	scanner := bufio.NewScanner(f1)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), "\t")
		termMap[splitted[0]] = splitted[1]
	}

	// Single characters.
	f2, err := os.Open(charFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	scanner = bufio.NewScanner(f2)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), "\t")
		for _, ch := range splitted[1] {
			charMap[string(ch)] = splitted[0]
		}
	}

	// Anthology files.
	files, err := ioutil.ReadDir(anthologyDirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		anthologyFilenames = append(anthologyFilenames, f.Name())
	}
}

// Convert converts a good ole article into CHU, Ing-Wen.
// It first replaces the multi-character terms and than perform single-term replacements by looping through the corpus character by character.
func Convert(goodOleArticle string) string {
	goodOleArticle = replaceTerms(goodOleArticle)

	var builder strings.Builder
	for _, ch := range goodOleArticle {
		builder.WriteString(mapChar(string(ch)))
	}

	str := builder.String()
	fmt.Println(str)
	return str
}

// GetRandomFromAnthology returns a converted random article from the anthology directory.
func GetRandomFromAnthology() string {
	// Open a random file.
	file, err := os.Open(anthologyDirPath + anthologyFilenames[rand.Intn(len(anthologyFilenames))])
	if err != nil {
		log.Println(err)
		return errString
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return errString
	}

	return Convert(string(content))
}

// replaceTerms replace all matching terms in the corpus with the terms loaded in the term map.
func replaceTerms(corpus string) string {
	// This takes O((len(corpus)+maxlen(term))*len(termMap)), which is super bad, but I'm too lazy to write a look ahead parse or something like that, so don't judge me.
	for key, val := range termMap {
		corpus = strings.ReplaceAll(corpus, key, val)
	}

	return corpus
}

// MapChar maps a chinese character to a string in the char map.
// If the chinese character doesn't exist in the char map, the function will return the input string.
func mapChar(r string) string {
	if val, ok := charMap[r]; ok {
		return val
	}

	return r
}
