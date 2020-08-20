package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Token struct {
	token string
	tokeType string
}

func main() {
	filenames := filenamesToProcess()

	for _, filename := range filenames {
		fmt.Println("Processing " + filename)

		tokens := tokenize(filename)
		tokenOutputName := strings.Replace(filename, ".jack", "T.xml", 1)
		writeTokens(tokens, tokenOutputName)

		fmt.Println("Finished generating tokens")

		ti := NewTokenIterator(tokens)
		xmlEl := process(&ti)

		treeOutputFile := strings.Replace(filename, ".jack", ".xml", 1)
		writeSyntaxTree(xmlEl, treeOutputFile)
	}
}

func filenamesToProcess() []string {
	target := os.Args[1]
	info, err := os.Stat(target)

	if err != nil {
		panic(err)
	}

	filenames := make([]string, 0)

	if info.Mode().IsDir() {
		err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".jack") {
				filenames = append(filenames, path)
			}
			return nil
		})
	} else {
		filenames = append(filenames, target)
	}


	return filenames
}

func writeTokens(tokens []Token, outputFileName string) {
	fout, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	_, _ = writer.WriteString("<tokens>\n")

	for _, token := range tokens {
		_, _ = writer.WriteString(fmt.Sprintf("<%s> %s </%s>\n", token.tokeType, token.token, token.tokeType))
	}

	_, _ = writer.WriteString("</tokens>\n")
}

func writeSyntaxTree(element XmlElement, outputFileName string) {
	fout, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	writeXmlTag(element, 0, writer)
}

func writeXmlTag(element XmlElement, prefix int, writer *bufio.Writer) {
	if element.singleEl {
		_, _ = writer.WriteString(fmt.Sprintf("%s<%s> %s </%s>\n", buildPrefix(prefix), element.tag, element.textContent, element.tag))
	} else {
		p := buildPrefix(prefix)
		_, _ = writer.WriteString(fmt.Sprintf("%s<%s>\n", p, element.tag))

		for _, child := range element.childes {
			writeXmlTag(child, prefix + 2, writer)
		}

		_, _ = writer.WriteString(fmt.Sprintf("%s</%s>\n", p, element.tag))
	}
}

func buildPrefix(len int) string {
	r := make([]rune, len, len)

	for i := 0; i < len; i++ {
		r[i] = ' '
	}

	return string(r)
}
