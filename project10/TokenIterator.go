package main

type TokenIterator struct {
	tokens []Token
	cur int
}

func NewTokenIterator(tokens []Token) TokenIterator {
	return TokenIterator{
		tokens: tokens,
		cur:    -1,
	}
}

func (this TokenIterator) hasNext() bool {
	return this.cur < len(this.tokens) - 1
}

func (this *TokenIterator) next() Token {
	if this.hasNext() {
		this.cur += 1
		return this.tokens[this.cur]
	}

	return Token{}
}

func (this TokenIterator) seeNext() Token {
	return this.tokens[this.cur + 1]
}

func (this TokenIterator) current() Token {
	if this.cur == -1 {
		return this.tokens[0]
	}
	return this.tokens[this.cur]
}