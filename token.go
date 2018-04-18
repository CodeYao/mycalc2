package main

const MAX_TOKENSIZE int = 100

type TokenKind int

type TokenType int

type StateType int

const (
	BAD_TOKEN TokenKind = iota
	NUMBER_TOKEN
	PARAM_TOKEN
	STATE_TOKEN
	STATE_TYPE_TOKEN
	TOKEN_TYPE_TOKEN
	ADD_OPERATOR_TOKEN
	SUB_OPERATOR_TOKEN
	MUL_OPERATOR_TOKEN
	DIV_OPERATOR_TOKEN
	ASS_OPERATOR_TOKEN
	LEFT_PAREN_TOKEN
	RIGHT_PAREN_TOKEN
	END_OF_LINE_TOKEN
)

const (
	INT8 TokenType = iota
	INT16
	INT32
	INT64
	UINT8
	UINT16
	UINT32
	UINT64
	BOOL
	FLOAT32
	FLOAT64
	STRING
	CHAR
)

const (
	LET StateType = iota
	SET
)

var KeyWords = []string{"int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "bool", "float32", "float64", "string", "char"}
var StatementWords = []string{"let", "set"}

type Token struct {
	kind      TokenKind
	value     float32 //interface{}
	str       string
	tokenType TokenType
	stateType StateType
}
