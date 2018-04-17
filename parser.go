package main

import (
	"bufio"
	"fmt"
	"os"
)

var st_look_ahead_token Token
var st_look_ahead_token_exists int

func my_get_token(token *Token) {
	if st_look_ahead_token_exists == 1 {
		*token = st_look_ahead_token
		st_look_ahead_token_exists = 0
	} else {
		getToken(token)
	}
}

func unget_token(token *Token) {
	st_look_ahead_token = *token
	st_look_ahead_token_exists = 1
}

func parse_primary_expression() float32 {
	var token Token
	var value float32 = 0.0
	var minus_flages int = 0

	my_get_token(&token)
	if token.kind == SUB_OPERATOR_TOKEN {
		minus_flages = 1
	} else {
		unget_token(&token)
	}

	my_get_token(&token)
	if token.kind == NUMBER_TOKEN {
		if minus_flages == 1 {
			return -token.value
		} else {
			fmt.Println("token.value...", token.value)
			return token.value
		}
	} else if token.kind == LEFT_PAREN_TOKEN {
		value = parse_expression()
		my_get_token(&token)
		if token.kind != RIGHT_PAREN_TOKEN {
			fmt.Println("missing ')' error.")
			os.Exit(1)
		}
		if minus_flages == 1 {
			return -value
		} else {
			return value
		}
	}
	fmt.Println("syntax error.")
	os.Exit(1)
	return 0.0 /* make compiler happy */
}

func parse_term() float32 {
	var v1 float32
	var v2 float32
	var token Token

	v1 = parse_primary_expression()
	for {
		my_get_token(&token)
		if token.kind != DIV_OPERATOR_TOKEN && token.kind != MUL_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_primary_expression()
		fmt.Println("kind...", token.kind, "str...", token.str)
		if token.kind == MUL_OPERATOR_TOKEN {
			v1 *= v2
		} else if token.kind == DIV_OPERATOR_TOKEN {
			v1 /= v2
		}
	}
	fmt.Println("v1...", v1)
	return v1
}

func parse_expression() float32 {
	var v1 float32
	var v2 float32
	var token Token

	v1 = parse_term()
	for {
		my_get_token(&token)
		if token.kind != ADD_OPERATOR_TOKEN && token.kind != SUB_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_term()
		if token.kind == ADD_OPERATOR_TOKEN {
			v1 += v2
		} else if token.kind == SUB_OPERATOR_TOKEN {
			v1 -= v2
		} else {
			unget_token(&token)
		}
	}
	return v1
}

func parse_line() float32 {
	var value float32

	st_look_ahead_token_exists = 0
	value = parse_expression()

	return value
}

func main() {
	var value float32

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("please input:")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println("There ware errors reading,exiting program.")
		return
	}
	set_line([]rune(input))
	value = parse_line()
	fmt.Println(">>", value)
}
