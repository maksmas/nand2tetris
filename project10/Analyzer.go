package main

import (
	"fmt"
	"strings"
)

type XmlElement struct {
	tag         string
	textContent string
	singleEl    bool

	childes []XmlElement
}

const UnaryOperators = "-~"
const BinaryOperators = "+-*/&|<>+="

func process(ti *TokenIterator) XmlElement {
	return compileClass(ti)
}

func compileClass(ti *TokenIterator) XmlElement {
	classTag := XmlElement{
		tag:      "class",
		singleEl: false,
		childes:  make([]XmlElement, 3),
	}

	classTag.childes[0] = eat(ti.next(), "keyword", "class")
	classTag.childes[1] = eatByType(ti.next(), "identifier")
	classTag.childes[2] = eatSymbol(ti.next(), "{")

	for ti.seeNext().token == "static" || ti.seeNext().token == "field" {
		classTag.childes = append(classTag.childes, compileClassVar(ti))
	}

	for ti.seeNext().token == "constructor" || ti.seeNext().token == "function" || ti.seeNext().token == "method" {
		classTag.childes = append(classTag.childes, compileSubroutineDec(ti))
	}

	classTag.childes = append(classTag.childes, eatSymbol(ti.next(), "}"))

	return classTag
}

func compileClassVar(ti *TokenIterator) XmlElement {
	classVarDec := XmlElement{
		tag:      "classVarDec",
		singleEl: false,
		childes:  make([]XmlElement, 3),
	}

	classVarDec.childes[0] = eatByType(ti.next(), "keyword")
	varType := ti.next()

	if varType.tokeType == "keyword" {
		classVarDec.childes[1] = eatByType(varType, "keyword")
	} else {
		classVarDec.childes[1] = eatByType(varType, "identifier")
	}
	classVarDec.childes[2] = eatByType(ti.next(), "identifier")

	for ti.seeNext().token == "," {
		classVarDec.childes = append(classVarDec.childes, eatSymbol(ti.next(), ","))
		classVarDec.childes = append(classVarDec.childes, eatByType(ti.next(), "identifier"))
	}

	classVarDec.childes = append(classVarDec.childes, eatSymbol(ti.next(), ";"))
	return classVarDec
}

func compileSubroutineDec(ti *TokenIterator) XmlElement {
	subroutineDec := XmlElement{
		tag:      "subroutineDec",
		singleEl: false,
		childes:  make([]XmlElement, 7),
	}

	subroutineDec.childes[0] = eatByType(ti.next(), "keyword")
	if ti.seeNext().tokeType == "keyword" {
		subroutineDec.childes[1] = eatByType(ti.next(), "keyword")
	} else {
		subroutineDec.childes[1] = eatByType(ti.next(), "identifier")
	}
	subroutineDec.childes[2] = eatByType(ti.next(), "identifier")
	subroutineDec.childes[3] = eatSymbol(ti.next(), "(")
	subroutineDec.childes[4] = compileParamList(ti)
	subroutineDec.childes[5] = eatSymbol(ti.next(), ")")
	subroutineDec.childes[6] = compileSubroutineBody(ti)

	return subroutineDec
}

func compileParamList(ti *TokenIterator) XmlElement {
	paramList := XmlElement{
		tag:      "parameterList",
		singleEl: false,
		childes:  make([]XmlElement, 0),
	}

	for ti.seeNext().token != ")" {
		paramList.childes = append(paramList.childes, eatByType(ti.next(), "keyword"))
		paramList.childes = append(paramList.childes, eatByType(ti.next(), "identifier"))

		if ti.seeNext().token == "," {
			paramList.childes = append(paramList.childes, eatSymbol(ti.next(), ","))
		}
	}

	return paramList
}

func compileSubroutineBody(ti *TokenIterator) XmlElement {
	body := XmlElement{
		tag:      "subroutineBody",
		singleEl: false,
		childes:  make([]XmlElement, 1),
	}

	body.childes[0] = eatSymbol(ti.next(), "{")

	for ti.seeNext().token == "var" {
		body.childes = append(body.childes, compileVarDec(ti))
	}

	body.childes = append(body.childes, compileStatements(ti))
	body.childes = append(body.childes, eatSymbol(ti.next(), "}"))

	return body
}

func compileVarDec(ti *TokenIterator) XmlElement {
	varDec := XmlElement{
		tag:      "varDec",
		singleEl: false,
		childes:  make([]XmlElement, 3),
	}

	varDec.childes[0] = eat(ti.next(), "keyword", "var")
	if ti.seeNext().tokeType == "keyword" {
		varDec.childes[1] = eatByType(ti.next(), "keyword")
	} else {
		varDec.childes[1] = eatByType(ti.next(), "identifier")
	}

	varDec.childes[2] = eatByType(ti.next(), "identifier")

	for ti.seeNext().token == "," {
		varDec.childes = append(varDec.childes, eatSymbol(ti.next(), ","))
		varDec.childes = append(varDec.childes, eatByType(ti.next(), "identifier"))
	}

	varDec.childes = append(varDec.childes, eatSymbol(ti.next(), ";"))

	return varDec
}

func compileStatements(ti *TokenIterator) XmlElement {
	statements := XmlElement{
		tag:      "statements",
		singleEl: false,
		childes:  make([]XmlElement, 0),
	}

	for ti.seeNext().token != "}" {
		if ti.seeNext().token == "let" {
			statements.childes = append(statements.childes, compileLetStatement(ti))
		} else if ti.seeNext().token == "if" {
			statements.childes = append(statements.childes, compileIfStatement(ti))
		} else if ti.seeNext().token == "while" {
			statements.childes = append(statements.childes, compileWhileStatement(ti))
		} else if ti.seeNext().token == "do" {
			statements.childes = append(statements.childes, compileDoStatement(ti))
		} else if ti.seeNext().token == "return" {
			statements.childes = append(statements.childes, compileReturnStatement(ti))
		} else {
			panic("Wrong statement " + ti.seeNext().token)
		}
	}

	return statements
}

func compileLetStatement(ti *TokenIterator) XmlElement {
	let := XmlElement{
		tag:      "letStatement",
		singleEl: false,
		childes:  make([]XmlElement, 2),
	}

	let.childes[0] = eat(ti.next(), "keyword", "let")
	let.childes[1] = eatByType(ti.next(), "identifier")

	if ti.seeNext().token == "[" {
		let.childes = append(let.childes, eatSymbol(ti.next(), "["))
		let.childes = append(let.childes, compileExpression(ti))
		let.childes = append(let.childes, eatSymbol(ti.next(), "]"))
	}

	let.childes = append(let.childes, eatSymbol(ti.next(), "="))
	let.childes = append(let.childes, compileExpression(ti))
	let.childes = append(let.childes, eatSymbol(ti.next(), ";"))

	return let
}

func compileIfStatement(ti *TokenIterator) XmlElement {
	ifStatement := XmlElement{
		tag:      "ifStatement",
		singleEl: false,
		childes:  make([]XmlElement, 7, 11),
	}

	ifStatement.childes[0] = eat(ti.next(), "keyword", "if")
	ifStatement.childes[1] = eatSymbol(ti.next(), "(")
	ifStatement.childes[2] = compileExpression(ti)
	ifStatement.childes[3] = eatSymbol(ti.next(), ")")
	ifStatement.childes[4] = eatSymbol(ti.next(), "{")
	ifStatement.childes[5] = compileStatements(ti)
	ifStatement.childes[6] = eatSymbol(ti.next(), "}")

	if ti.seeNext().token == "else" {
		ifStatement.childes = append(ifStatement.childes, eat(ti.next(), "keyword", "else"))
		ifStatement.childes = append(ifStatement.childes, eatSymbol(ti.next(), "{"))
		ifStatement.childes = append(ifStatement.childes, compileStatements(ti))
		ifStatement.childes = append(ifStatement.childes, eatSymbol(ti.next(), "}"))
	}

	return ifStatement
}

func compileWhileStatement(ti *TokenIterator) XmlElement {
	while := XmlElement{
		tag:         "whileStatement",
		singleEl:    false,
		childes:     make([]XmlElement, 7, 7),
	}

	while.childes[0] = eat(ti.next(), "keyword", "while")
	while.childes[1] = eatSymbol(ti.next(), "(")
	while.childes[2] = compileExpression(ti)
	while.childes[3] = eatSymbol(ti.next(), ")")
	while.childes[4] = eatSymbol(ti.next(), "{")
	while.childes[5] = compileStatements(ti)
	while.childes[6] = eatSymbol(ti.next(), "}")

	return while
}

func compileDoStatement(ti *TokenIterator) XmlElement {
	do := XmlElement{
		tag:      "doStatement",
		singleEl: false,
		childes:  make([]XmlElement, 1),
	}

	do.childes[0] = eat(ti.next(), "keyword", "do")
	ti.next()
	do.childes = append(do.childes, addSubroutineCall(ti)...)
	do.childes = append(do.childes, eatSymbol(ti.next(), ";"))

	return do
}

func compileReturnStatement(ti *TokenIterator) XmlElement {
	returnStatement := XmlElement{
		tag:      "returnStatement",
		singleEl: false,
		childes:  make([]XmlElement, 1),
	}

	returnStatement.childes[0] = eat(ti.next(), "keyword", "return")

	if ti.seeNext().token != ";" {
		returnStatement.childes = append(returnStatement.childes, compileExpression(ti))
	}

	returnStatement.childes = append(returnStatement.childes, eatSymbol(ti.next(), ";"))

	return returnStatement
}

func compileExpression(ti *TokenIterator) XmlElement {
	expr := XmlElement{
		tag:      "expression",
		singleEl: false,
		childes:  make([]XmlElement, 1),
	}

	expr.childes[0] = compileTerm(ti)

	for (len(ti.seeNext().token) == 1 && strings.Contains(BinaryOperators, ti.seeNext().token)) || ti.seeNext().token == "&lt;" || ti.seeNext().token == "&gt;" || ti.seeNext().token == "&amp;" {
		expr.childes = append(expr.childes, eatByType(ti.next(), "symbol"))
		expr.childes = append(expr.childes, compileTerm(ti))
	}

	return expr
}

func compileTerm(ti *TokenIterator) XmlElement {
	term := XmlElement{
		tag:      "term",
		singleEl: false,
		childes:  make([]XmlElement, 0),
	}

	token := ti.next()

	if token.tokeType == "integerConstant" || token.tokeType == "stringConstant" || token.tokeType == "keyword" {
		term.childes = append(term.childes, eatByType(token, token.tokeType))
	} else if token.tokeType == "identifier" {
		next := ti.seeNext()
		if next.token == "." || next.token == "(" {
			// func
			term.childes = append(term.childes, addSubroutineCall(ti)...)
		} else {
			// var
			term.childes = append(term.childes, eatByType(token, "identifier"))

			if ti.seeNext().token == "[" {
				term.childes = append(term.childes, eatSymbol(ti.next(), "["))
				term.childes = append(term.childes, compileExpression(ti))
				term.childes = append(term.childes, eatSymbol(ti.next(), "]"))
			}
		}
	} else if token.token == "(" {
		term.childes = append(term.childes, eatSymbol(token, "("))
		term.childes = append(term.childes, compileExpression(ti))
		term.childes = append(term.childes, eatSymbol(ti.next(), ")"))
	} else if len(token.token) == 1 && strings.Contains(UnaryOperators, token.token) {
		term.childes = append(term.childes, eatByType(token, "symbol"))
		term.childes = append(term.childes, compileTerm(ti))
	} else {
		panic("Wrong term")
	}

	return term
}

func compileExpressionList(ti *TokenIterator) XmlElement {
	elements := XmlElement{
		tag:      "expressionList",
		singleEl: false,
		childes:  make([]XmlElement, 0),
	}

	if ti.seeNext().token != ")" {
		elements.childes = append(elements.childes, compileExpression(ti))

		for ti.seeNext().token == "," {
			elements.childes = append(elements.childes, eatSymbol(ti.next(), ","))
			elements.childes = append(elements.childes, compileExpression(ti))
		}
	}

	return elements
}

func addSubroutineCall(ti *TokenIterator) []XmlElement {
	els := make([]XmlElement, 0)

	if ti.seeNext().token == "." {
		els = append(els, eatByType(ti.current(), "identifier"))
		els = append(els, eatSymbol(ti.next(), "."))
		ti.next()
	}

	els = append(els, eatByType(ti.current(), "identifier"))
	els = append(els, eatSymbol(ti.next(), "("))
	els = append(els, compileExpressionList(ti))
	els = append(els, eatSymbol(ti.next(), ")"))
	return els
}

func eat(token Token, expectedType string, expectedToken string) XmlElement {
	if token.token != expectedToken || token.tokeType != expectedType {
		panic(fmt.Sprintf("Expected %s %s, but got %s %s\n", expectedToken, expectedType, token.token, token.tokeType))
	}

	return XmlElement{
		tag:         expectedType,
		textContent: expectedToken,
		singleEl:    true,
	}
}

func eatByType(token Token, expectedType string) XmlElement {
	if token.tokeType != expectedType {
		panic("Expected " + expectedType + ", got " + token.tokeType)
	}

	return XmlElement{
		tag:         expectedType,
		textContent: token.token,
		singleEl:    true,
	}
}

func eatSymbol(token Token, expectedSymbol string) XmlElement {
	if token.tokeType != "symbol" || token.token != expectedSymbol {
		panic(fmt.Sprintf("Expected %s symbol, got %s %s", expectedSymbol, token.token, token.tokeType))
	}

	return XmlElement{
		tag:         "symbol",
		textContent: expectedSymbol,
		singleEl:    true,
	}
}
