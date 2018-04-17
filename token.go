package main

const MAX_TOKENSIZE int = 100

type TokenKind int

const (
	BAD_TOKEN TokenKind = iota
	NUMBER_TOKEN
	ADD_OPERATOR_TOKEN
	SUB_OPERATOR_TOKEN
	MUL_OPERATOR_TOKEN
	DIV_OPERATOR_TOKEN
	LEFT_PAREN_TOKEN
	RIGHT_PAREN_TOKEN
	END_OF_LINE_TOKEN
)

type Token struct {
	kind  TokenKind
	value float32
	str   string
}
