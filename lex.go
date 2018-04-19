package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var st_line []rune
var st_line_pos int

type LexerStatus int

const (
	INITIAL_STATUS LexerStatus = iota
	IN_INT_PART_STATUS
	DOT_STATUS
	IN_FRAC_PART_STATUS

	FIRST_PARAM_STATUS
	FOLLOW_PARAM_STATUS
)

func getToken(token *Token) {
	var out_pos int = 0
	var status LexerStatus = INITIAL_STATUS
	var current_char rune
	token.str = ""
	token.kind = BAD_TOKEN
	for {
		if st_line[st_line_pos] == '\000' {
			break
		}
		current_char = st_line[st_line_pos]
		//fmt.Println("current_char---(", string(current_char), ")")
		if (status == IN_INT_PART_STATUS || status == IN_FRAC_PART_STATUS) && !unicode.IsDigit(current_char) && current_char != '.' {
			token.kind = NUMBER_TOKEN
			if status == IN_INT_PART_STATUS {
				value, _ := strconv.ParseInt(token.str, 10, 64)
				token.value = int64(value)
			} else if status == IN_FRAC_PART_STATUS {
				value, _ := strconv.ParseFloat(token.str, 64)
				token.value = float64(value)
			}
			//fmt.Println("current_char---(", token.str, ")")
			return
		}

		if (status == FIRST_PARAM_STATUS || status == FOLLOW_PARAM_STATUS) && !unicode.IsDigit(current_char) && current_char != '_' && !unicode.IsLetter(current_char) {
			fmt.Println("current_char---(", token.str, ")")
			if IsKeyWord(token.str) {
				token.kind = TOKEN_TYPE_TOKEN
			} else if IsStatement(token.str) {
				token.kind = STATE_TYPE_TOKEN
			} else {
				token.kind = STATE_TOKEN
			}
			return
		}

		if unicode.IsSpace(current_char) {
			if current_char == '\n' {
				token.kind = END_OF_LINE_TOKEN
				return
			}
			st_line_pos++
			continue
		}
		if out_pos >= MAX_TOKENSIZE-1 {
			fmt.Println("token too long.")
			os.Exit(1)
		}

		token.str += string(st_line[st_line_pos])
		st_line_pos++
		out_pos++

		if current_char == '+' {
			token.kind = ADD_OPERATOR_TOKEN
			return
		} else if current_char == '-' {
			token.kind = SUB_OPERATOR_TOKEN
			return
		} else if current_char == '*' {
			token.kind = MUL_OPERATOR_TOKEN
			return
		} else if current_char == '/' {
			token.kind = DIV_OPERATOR_TOKEN
			return
		} else if current_char == '=' {
			token.kind = ASS_OPERATOR_TOKEN
			return
		} else if unicode.IsDigit(current_char) {
			if status == INITIAL_STATUS {
				status = IN_INT_PART_STATUS
			} else if status == DOT_STATUS {
				status = IN_FRAC_PART_STATUS
			} else if status == FIRST_PARAM_STATUS {
				status = FOLLOW_PARAM_STATUS
			}
		} else if current_char == '.' {
			if status == IN_INT_PART_STATUS {
				status = DOT_STATUS
			} else {
				fmt.Println("syntax error.")
				os.Exit(1)
			}
		} else if unicode.IsLetter(current_char) {
			if status == INITIAL_STATUS {
				status = FIRST_PARAM_STATUS
			}
		} else if current_char == '_' {
			if status == FIRST_PARAM_STATUS {
				status = FOLLOW_PARAM_STATUS
			} else {
				fmt.Println("error! a variable must begin with an alphabet")
				os.Exit(1)
			}
		} else if current_char == '(' {
			token.kind = LEFT_PAREN_TOKEN
			return
		} else if current_char == ')' {
			token.kind = RIGHT_PAREN_TOKEN
			return
		} else {
			fmt.Println("bad character(", current_char, ")")
			os.Exit(1)
		}
	}
}

func set_line(line []rune) {
	st_line = line
	st_line_pos = 0
}

func IsKeyWord(str string) bool {
	for _, keyword := range KeyWords {
		if keyword == str {
			return true
		}
	}
	return false
}

func IsStatement(str string) bool {
	for _, statementword := range StatementWords {
		if statementword == str {
			return true
		}
	}
	return false
}

// func parse_line(buf []rune) {
// 	var token Token

// 	set_line(buf)

// 	for {
// 		getToken(&token)
// 		if token.kind == END_OF_LINE_TOKEN {
// 			break
// 		} else {
// 			fmt.Println("kind...", token.kind, "str...", token.str)
// 		}
// 	}
// }
// func main() {
// 	inputReader := bufio.NewReader(os.Stdin)
// 	fmt.Println("please input:")
// 	input, err := inputReader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("There ware errors reading,exiting program.")
// 		return
// 	}
// 	parse_line([]rune(input))
// }
