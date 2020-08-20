package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var keywords = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}
const symbols = "{}()[].,;+-*/&|<>=~"

func tokenize(filename string) []Token {
	fin, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	reader := bufio.NewReader(fin)
	return tokenizeFile(reader)
}

func tokenizeFile(reader *bufio.Reader) []Token {
	tokens := make([]Token, 0)

	for ;; {
		token, tokenType, eof := readToken(reader)

		if eof {
			break
		} else if token != "" {
			tokens = append(tokens, Token{token, tokenType})
		}
	}
	

	return tokens
}

func readToken(reader *bufio.Reader) (token string, tokenType string, eof bool) {
	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			return "", "", true
		} else if err != nil {
			panic(err)
		}

		if r == '/' {
			if isComment := readComment(reader); !isComment {
				return "/", "symbol", false
			}
		} else if r >= '0' && r <= '9' {
			return readIntegerConstant(r, reader), "integerConstant", false
		} else if strings.Contains(symbols, string(r)) {
			return escapeXmlString(string(r)), "symbol", false
		} else if r == '"' {
			return readStringToken(reader), "stringConstant", false
		} else if isChar(r) {
			str := readChars(r, reader)
			strType := "identifier"

			if isKeyword(str) {
				strType = "keyword"
			}

			return str, strType, false
		}
	}

	return token, tokenType,false
}

func escapeXmlString(s string) string {
	if s == "<" {
		return "&lt;"
	} else if s == ">" {
		return "&gt;"
	} else if s == "&" {
		return "&amp;"
	}

	return s
}

func readComment(reader *bufio.Reader) (isComment bool) {
	r, _, err := reader.ReadRune()
	if err != nil {
		panic(err)
	}

	if r == '/' {
		_, _, _ = reader.ReadLine()
		return true
	} else if r == '*' {
		_, _, _ = reader.ReadRune()

		for ;; {
			r, _, err := reader.ReadRune()
			if err != nil {
				panic(err)
			}

			if r == '*' {
				r, _, err = reader.ReadRune()
				if err != nil {
					panic(err)
				}

				if r == '/' {
					return true
				}
			}
		}

		return true
	}

	return false
}

func readStringToken(reader *bufio.Reader) string {
	token := make([]rune, 0)

	for ;; {
		rune, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if rune == '"' {
			break
		}

		token = append(token, rune)
	}

	return string(token)
}

func readIntegerConstant(r rune, reader *bufio.Reader) string {
	token := make([]rune, 0)
	for err := error(nil); r >= '0' && r <= '9'; r, _, err = reader.ReadRune() {
		if err != nil {
			panic(err)
		}

		token = append(token, r)
	}

	err := reader.UnreadRune()
	if err != nil {
		panic(err)
	}

	return string(token)
}

func readChars(r rune, reader *bufio.Reader) string {
	token := make([]rune, 0)
	for err := error(nil); isChar(r); r, _, err = reader.ReadRune() {
		if err != nil {
			panic(err)
		}

		token = append(token, r)
	}

	err := reader.UnreadRune()
	if err != nil {
		panic(err)
	}

	return string(token)
}

func isChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isKeyword(s string) bool {
	for _, kwd := range keywords {
		if kwd == s {
			return true
		}
	}

	return false
}
