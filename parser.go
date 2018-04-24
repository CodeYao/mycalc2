package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

var st_look_ahead_token Token
var st_look_ahead_token_exists int

var stl_bak int = 0

/*0表示不存在,1表示存在,2表示存在且逻辑表达式为true,3表示逻辑表达式为false*/
var if_type int = 0
var else_type int = 0

//var if_st_lines [][]rune

var paramList map[string]Token //变量列表

var fi *os.File
var inputReader *bufio.Reader

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

func parse_primary_expression() interface{} {
	var token Token
	var value interface{}
	var minus_flages int = 0

	my_get_token(&token)
	if token.kind == SUB_OPERATOR_TOKEN {
		minus_flages = 1
	} else {
		unget_token(&token)
	}

	my_get_token(&token)

	//判断是否声明变量
	if token.kind == STATE_TYPE_TOKEN {
		state_token := token
		my_get_token(&token)
		//获取变量名，放入map
		if token.kind == STATE_TOKEN {
			if state_token.str == "let" {
				token.stateType = LET
			} else if state_token.str == "set" {
				token.stateType = SET
			}
			stk := token
			my_get_token(&token)
			//获取变量类型
			if token.kind == TOKEN_TYPE_TOKEN {
				stk.tokenType = getTokenType(token.str)
				if stk.tokenType == ERRORTYPE {
					fmt.Println("error type : ", token.str)
					os.Exit(1)
				}

				my_get_token(&token)

				//变量后续是否为赋值操作
				if token.kind == ASS_OPERATOR_TOKEN {
					value = parse_expression()
					value = getValue(stk.tokenType, value, minus_flages)
					// tokentype := getTokenType(reflect.TypeOf(value).String())
					// if stk.tokenType != tokentype {
					// 	fmt.Println("The type of variable assignment is not consistent : ", token.str)
					// 	os.Exit(1)
					// }

					if _, ok := paramList[stk.str]; ok {
						fmt.Println("error the variable is existed : ", stk.str)
						os.Exit(1)
					} else {
						stk.value = value
						paramList[stk.str] = stk
					}
				} else {
					unget_token(&token)
					paramList[stk.str] = stk
					return value
				}
			}
		}
	}

	/*获取变量，返回value*/
	if token.kind == STATE_TOKEN {
		//判断是否已经声明
		if tk, ok := paramList[token.str]; ok {
			my_get_token(&token)
			//变量后续是否为赋值操作
			if token.kind == ASS_OPERATOR_TOKEN {
				value = parse_expression()
				//value = getValue(tk.tokenType, value, minus_flages)
				tk.value = getValue(tk.tokenType, value, minus_flages)
				paramList[tk.str] = tk
				fmt.Println(tk.str, " :: ", tk.value)
			} else {
				unget_token(&token)
				if t, ok := paramList[tk.str]; ok {
					fmt.Println(t.str, " : ", t.value)
					//遗留问题
					// if minus_flages == 1 {
					// 	return getValue(t.tokenType, value, minus_flages)
					// }
					value = t.value
				} else {
					fmt.Println("Undeclared variables : ", tk.str)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("an undeclared variable : ", token.str)
			os.Exit(1)
		}
	}
	//如果是else
	if token.kind == ELSE_TOKEN {
		value = true
		if else_type == 2 && if_type == 3 {
			my_get_token(&token)
			if token.kind == LEFT_BRACES_TOKEN {
				getCode()
				else_type = 0
				if_type = 0
			} else if token.kind == IF_TOKEN {
				//unget_token(&token)
				else_type = 0
			}
		} else if else_type == 2 && if_type == 0 {
			fmt.Println("Miss if error")
			os.Exit(1)
		}
	}
	//如果是if
	if token.kind == IF_TOKEN {
		if_type = 1
		my_get_token(&token)
		fmt.Println("if token", token.str)
		if token.kind == LEFT_PAREN_TOKEN {
			value = parse_logic_expression()
			my_get_token(&token)
			if token.kind != RIGHT_PAREN_TOKEN {
				fmt.Println("missing ')' error.")
				os.Exit(1)
			}
			if value.(bool) {
				fmt.Println("if true")
				if_type = 2
				my_get_token(&token)
				if token.kind == LEFT_BRACES_TOKEN {
					getCode()
				}
				if_type = 0
			} else {
				if_type = 3
				else_type = 2
				fmt.Println("if false")
				skipCode()
			}
		} else {
			fmt.Println("missing '(' error.")
			os.Exit(1)
		}
	}
	//如果是常量
	if token.kind == NUMBER_TOKEN {
		value = token.value
	} else if token.kind == LEFT_PAREN_TOKEN {
		value = parse_expression()
		my_get_token(&token)
		if token.kind != RIGHT_PAREN_TOKEN {
			fmt.Println("missing ')' error.")
			os.Exit(1)
		}
	} else if token.kind == CHAR_TOKEN || token.kind == STRING_TOKEN {
		value = string([]rune(token.str)[1 : len(token.str)-1])
		fmt.Println(token.str, " ----------- ", value.(string))

	} else if token.kind == BOOL_TOKEN {
		if token.str == "true" {
			value = true
		} else if token.str == "false" {
			value = false
		}
	} else {
		unget_token(&token)
	}
	//fmt.Println(reflect.ValueOf(value), "************")
	if reflect.TypeOf(value).String() == "float32" || reflect.TypeOf(value).String() == "float64" {
		value = reflect.ValueOf(value).Float()
		if minus_flages == 1 {
			value = -reflect.ValueOf(value).Float()
		}
	} else if reflect.TypeOf(value).String() == "int8" || reflect.TypeOf(value).String() == "int16" || reflect.TypeOf(value).String() == "int32" || reflect.TypeOf(value).String() == "int64" {
		value = reflect.ValueOf(value).Int()
		if minus_flages == 1 {
			value = -reflect.ValueOf(value).Int()
		}
		value = float64(value.(int64))
	} else if reflect.TypeOf(value).String() == "uint8" || reflect.TypeOf(value).String() == "uint16" || reflect.TypeOf(value).String() == "uint32" || reflect.TypeOf(value).String() == "uint64" {
		value = reflect.ValueOf(value).Uint()
		if minus_flages == 1 {
			value = -reflect.ValueOf(value).Uint()
		}
		value = float64(value.(uint64))
	} // else {
	// 	fmt.Println("These Type can not be negative ", reflect.TypeOf(value).String())
	// }
	return value
}

func parse_term() interface{} {
	var v1 interface{}
	var v2 interface{}
	var token Token

	v1 = parse_primary_expression()
	for {
		my_get_token(&token)
		if token.kind != DIV_OPERATOR_TOKEN && token.kind != MUL_OPERATOR_TOKEN && token.kind != MOD_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_primary_expression()

		if reflect.TypeOf(v1).String() == reflect.TypeOf(v2).String() {
			if token.kind == MUL_OPERATOR_TOKEN {
				if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
					v1 = reflect.ValueOf(v1).Float() * reflect.ValueOf(v2).Float()
				} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
					v1 = reflect.ValueOf(v1).Int() * reflect.ValueOf(v2).Int()
				} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
					v1 = reflect.ValueOf(v1).Uint() * reflect.ValueOf(v2).Uint()
				} else if reflect.TypeOf(v1).String() == "rune" {
					v1 = v1.(rune) * v2.(rune)
				} else {
					fmt.Println("These Type can not add ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
			} else if token.kind == DIV_OPERATOR_TOKEN {
				if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
					v1 = reflect.ValueOf(v1).Float() / reflect.ValueOf(v2).Float()
				} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
					v1 = reflect.ValueOf(v1).Int() / reflect.ValueOf(v2).Int()
				} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
					v1 = reflect.ValueOf(v1).Uint() / reflect.ValueOf(v2).Uint()
				} else if reflect.TypeOf(v1).String() == "rune" {
					v1 = v1.(rune) / v2.(rune)
				} else {
					fmt.Println("These Type can not sub ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
			} else if token.kind == MOD_OPERATOR_TOKEN {
				//fmt.Println("MOD_OPERATOR_TOKEN")
				if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
					v1 = int64(reflect.ValueOf(v1).Float()) % int64(reflect.ValueOf(v2).Float())
				} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
					v1 = reflect.ValueOf(v1).Int() % reflect.ValueOf(v2).Int()
				} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
					v1 = reflect.ValueOf(v1).Uint() % reflect.ValueOf(v2).Uint()
				} else {
					fmt.Println("These Type can not mod ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
				v1 = float64(v1.(int64))
			} else {
				unget_token(&token)
			}
		} else {
			if token.kind == MOD_OPERATOR_TOKEN {
				if !strings.Contains(reflect.ValueOf(v1).String(), ".") && !strings.Contains(reflect.ValueOf(v2).String(), ".") {
					if reflect.TypeOf(v1).String() == "float64" {
						v1 = int64(v1.(float64))
					} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
						v1 = reflect.ValueOf(v1).Int()
					} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
						v1 = reflect.ValueOf(v1).Uint()
						v1 = int64(v1.(uint64))
					}
					if reflect.TypeOf(v2).String() == "float64" {
						v2 = int64(v2.(float64))
					} else if reflect.TypeOf(v2).String() == "int8" || reflect.TypeOf(v2).String() == "int16" || reflect.TypeOf(v2).String() == "int32" || reflect.TypeOf(v2).String() == "int64" {
						v2 = reflect.ValueOf(v2).Int()
					} else if reflect.TypeOf(v2).String() == "uint8" || reflect.TypeOf(v2).String() == "uint16" || reflect.TypeOf(v2).String() == "uint32" || reflect.TypeOf(v2).String() == "uint64" {
						v2 = reflect.ValueOf(v2).Uint()
						v2 = int64(v2.(uint64))
					}
					v1 = v1.(int64) % v2.(int64)
					v1 = float64(v1.(int64))
				} else {
					fmt.Println("These Type can not mod ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
			} else {
				fmt.Println("Type inconsistency ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
			}
		}
	}
	//fmt.Println("v1...", v1)
	return v1
}

func parse_expression() interface{} {
	var v1 interface{}
	var v2 interface{}
	var token Token

	v1 = parse_term()
	for {
		my_get_token(&token)
		if token.kind != ADD_OPERATOR_TOKEN && token.kind != SUB_OPERATOR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_term()
		if reflect.TypeOf(v1).String() == reflect.TypeOf(v2).String() {
			if token.kind == ADD_OPERATOR_TOKEN {
				//fmt.Println("======", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
					v1 = reflect.ValueOf(v1).Float() + reflect.ValueOf(v2).Float()
				} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
					v1 = reflect.ValueOf(v1).Int() + reflect.ValueOf(v2).Int()
				} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
					v1 = reflect.ValueOf(v1).Uint() + reflect.ValueOf(v2).Uint()
				} else if reflect.TypeOf(v1).String() == "rune" {
					v1 = v1.(rune) + v2.(rune)
				} else if reflect.TypeOf(v1).String() == "string" {
					v1 = v1.(string) + v2.(string)
				} else {
					fmt.Println("These Type can not add ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
			} else if token.kind == SUB_OPERATOR_TOKEN {
				if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
					v1 = reflect.ValueOf(v1).Float() - reflect.ValueOf(v2).Float()
				} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
					v1 = reflect.ValueOf(v1).Int() - reflect.ValueOf(v2).Int()
				} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
					v1 = reflect.ValueOf(v1).Uint() - reflect.ValueOf(v2).Uint()
				} else if reflect.TypeOf(v1).String() == "rune" {
					v1 = v1.(rune) - v2.(rune)
				} else {
					fmt.Println("These Type can not sub ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
				}
			} else {
				unget_token(&token)
			}
		} else {
			fmt.Println("Type inconsistency ", reflect.TypeOf(v1).String(), " : ", reflect.TypeOf(v2).String())
		}
	}
	return v1
}

func parse_relation_expression() interface{} {
	var v1 interface{}
	var v2 interface{}
	var token Token

	v1 = parse_expression()
	for {
		my_get_token(&token)
		if token.kind != EQ_TOKEN && token.kind != GE_TOKEN && token.kind != GT_TOKEN && token.kind != LT_TOKEN && token.kind != LE_TOKEN && token.kind != NE_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_expression()
		if token.kind == EQ_TOKEN {
			v1 = (v1 == v2)
			//fmt.Println(v1.(bool))
		} else if token.kind == GE_TOKEN {
			//fmt.Println(reflect.ValueOf(v1), "GEGEGE...", reflect.ValueOf(v2))
			if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
				v1 = (reflect.ValueOf(v1).Float() >= reflect.ValueOf(v2).Float())
			} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
				v1 = (reflect.ValueOf(v1).Int() >= reflect.ValueOf(v2).Int())
			} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
				v1 = (reflect.ValueOf(v1).Uint() >= reflect.ValueOf(v2).Uint())
			} else if reflect.TypeOf(v1).String() == "rune" {
				v1 = (v1.(rune) >= v2.(rune))
			} else if reflect.TypeOf(v1).String() == "string" {
				v1 = (v1.(string) >= v2.(string))
			} else {
				fmt.Println("These Type can not >= between ", reflect.TypeOf(v1).String(), " and ", reflect.TypeOf(v2).String())
			}
			//fmt.Println(v1.(bool))
		} else if token.kind == GT_TOKEN {
			//fmt.Println(reflect.ValueOf(v1), "GEGEGE...", reflect.ValueOf(v2))
			if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
				v1 = (reflect.ValueOf(v1).Float() > reflect.ValueOf(v2).Float())
			} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
				v1 = (reflect.ValueOf(v1).Int() > reflect.ValueOf(v2).Int())
			} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
				v1 = (reflect.ValueOf(v1).Uint() > reflect.ValueOf(v2).Uint())
			} else if reflect.TypeOf(v1).String() == "rune" {
				v1 = (v1.(rune) > v2.(rune))
			} else if reflect.TypeOf(v1).String() == "string" {
				v1 = (v1.(string) > v2.(string))
			} else {
				fmt.Println("These Type can not >= between ", reflect.TypeOf(v1).String(), " and ", reflect.TypeOf(v2).String())
			}
			//fmt.Println(v1.(bool))
		} else if token.kind == LE_TOKEN {
			//fmt.Println(reflect.ValueOf(v1), "GEGEGE...", reflect.ValueOf(v2))
			if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
				v1 = (reflect.ValueOf(v1).Float() <= reflect.ValueOf(v2).Float())
			} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
				v1 = (reflect.ValueOf(v1).Int() <= reflect.ValueOf(v2).Int())
			} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
				v1 = (reflect.ValueOf(v1).Uint() <= reflect.ValueOf(v2).Uint())
			} else if reflect.TypeOf(v1).String() == "rune" {
				v1 = (v1.(rune) <= v2.(rune))
			} else if reflect.TypeOf(v1).String() == "string" {
				v1 = (v1.(string) <= v2.(string))
			} else {
				fmt.Println("These Type can not >= between ", reflect.TypeOf(v1).String(), " and ", reflect.TypeOf(v2).String())
			}
			//fmt.Println(v1.(bool))
		} else if token.kind == LT_TOKEN {
			//fmt.Println(reflect.ValueOf(v1), "GEGEGE...", reflect.ValueOf(v2))
			if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
				v1 = (reflect.ValueOf(v1).Float() < reflect.ValueOf(v2).Float())
			} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
				v1 = (reflect.ValueOf(v1).Int() < reflect.ValueOf(v2).Int())
			} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
				v1 = (reflect.ValueOf(v1).Uint() < reflect.ValueOf(v2).Uint())
			} else if reflect.TypeOf(v1).String() == "rune" {
				v1 = (v1.(rune) < v2.(rune))
			} else if reflect.TypeOf(v1).String() == "string" {
				v1 = (v1.(string) < v2.(string))
			} else {
				fmt.Println("These Type can not >= between ", reflect.TypeOf(v1).String(), " and ", reflect.TypeOf(v2).String())
			}
			//fmt.Println(v1.(bool))
		} else if token.kind == NE_TOKEN {
			//fmt.Println(reflect.ValueOf(v1), "GEGEGE...", reflect.ValueOf(v2))
			if reflect.TypeOf(v1).String() == "float32" || reflect.TypeOf(v1).String() == "float64" {
				v1 = (reflect.ValueOf(v1).Float() != reflect.ValueOf(v2).Float())
			} else if reflect.TypeOf(v1).String() == "int8" || reflect.TypeOf(v1).String() == "int16" || reflect.TypeOf(v1).String() == "int32" || reflect.TypeOf(v1).String() == "int64" {
				v1 = (reflect.ValueOf(v1).Int() != reflect.ValueOf(v2).Int())
			} else if reflect.TypeOf(v1).String() == "uint8" || reflect.TypeOf(v1).String() == "uint16" || reflect.TypeOf(v1).String() == "uint32" || reflect.TypeOf(v1).String() == "uint64" {
				v1 = (reflect.ValueOf(v1).Uint() != reflect.ValueOf(v2).Uint())
			} else if reflect.TypeOf(v1).String() == "rune" {
				v1 = (v1.(rune) != v2.(rune))
			} else if reflect.TypeOf(v1).String() == "string" {
				v1 = (v1.(string) != v2.(string))
			} else {
				fmt.Println("These Type can not >= between ", reflect.TypeOf(v1).String(), " and ", reflect.TypeOf(v2).String())
			}
			//fmt.Println(v1.(bool))
		}
	}
	return v1
}

func parse_logic_expression() interface{} {
	var v1 interface{}
	var v2 interface{}
	var token Token
	v1 = parse_relation_expression()
	for {
		my_get_token(&token)
		if token.kind != LOGICAL_AND_TOKEN && token.kind != LOGICAL_OR_TOKEN {
			unget_token(&token)
			break
		}
		v2 = parse_relation_expression()
		if token.kind == LOGICAL_AND_TOKEN {
			//fmt.Println(v1.(bool), "ANDAND", v2.(bool))
			v1 = (v1.(bool) && v2.(bool))
		} else if token.kind == LOGICAL_OR_TOKEN {
			//fmt.Println(v1.(bool), "OROROR", v2.(bool))
			v1 = (v1.(bool) || v2.(bool))
		}
	}
	return v1
}

//获取变量类型
func getTokenType(str string) TokenType {
	if str == "int8" {
		return INT8
	} else if str == "int16" {
		return INT16
	} else if str == "int32" {
		return INT32
	} else if str == "int64" {
		return INT64
	} else if str == "uint8" {
		return UINT8
	} else if str == "uint16" {
		return UINT16
	} else if str == "uint32" {
		return UINT32
	} else if str == "uint64" {
		return UINT64
	} else if str == "bool" {
		return BOOL
	} else if str == "float32" {
		return FLOAT32
	} else if str == "float64" {
		return FLOAT64
	} else if str == "string" {
		return STRING
	} else if str == "rune" {
		return CHAR
	} else {
		return ERRORTYPE
	}
}

//获取value类型
func getValue(t TokenType, value interface{}, minus_flages int) interface{} {
	if t == INT8 {
		if minus_flages == 1 {
			value = -int8(value.(float64))
		} else {
			value = int8(value.(float64))
		}
	} else if t == INT16 {
		if minus_flages == 1 {
			value = -int16(value.(float64))
		} else {
			value = int16(value.(float64))
		}
	} else if t == INT32 {
		if minus_flages == 1 {
			value = -int32(value.(float64))
		} else {
			value = int32(value.(float64))
		}
	} else if t == INT64 {
		if minus_flages == 1 {
			value = -int64(value.(float64))
		} else {
			value = int64(value.(float64))
		}
	} else if t == FLOAT32 {
		if minus_flages == 1 {
			value = -float32(value.(float64))
		} else {
			value = float32(value.(float64))
		}
	} else if t == FLOAT64 {
		if minus_flages == 1 {
			value = -value.(float64)
		} else {
			value = value.(float64)
		}
	} else if t == UINT8 {
		if minus_flages == 1 {
			value = uint8(-value.(float64))
		} else {
			value = uint8(value.(float64))
		}
	} else if t == UINT16 {
		if minus_flages == 1 {
			value = uint16(-value.(float64))
		} else {
			value = uint16(value.(float64))
		}
	} else if t == UINT32 {
		if minus_flages == 1 {
			value = uint32(-value.(float64))
		} else {
			value = uint32(value.(float64))
		}
	} else if t == UINT64 {
		if minus_flages == 1 {
			value = uint64(-value.(float64))
		} else {
			value = uint64(value.(float64))
		}
	} else if t == CHAR {
		r1 := []rune(value.(string))
		//暂时没有实现转义字符
		value = r1[0]
	} else if t == STRING {
		value = value.(string)
	} else if t == BOOL {
		value = value.(bool)
	}
	return value
}
func getCode() {
	var value interface{}
	for {
		input, _, c := inputReader.ReadLine()
		if c == io.EOF {
			fmt.Println("if error")
			fmt.Println("missing '}' error.")
			os.Exit(1)
		}
		fmt.Println(len(input), string(input), "get")

		if len(input) == 0 || strings.Replace(string(input), " ", "", -1) == "" { //跳过空行
			continue
		} else {
			line := string(input) + "\n"
			//fmt.Println("get", line)
			set_line([]rune(line))
			if strings.Replace(string(input), " ", "", -1) == "}" {
				break
			}
			value = parse_line()
			fmt.Println(">>", reflect.ValueOf(value))
		}
	}
}

//跳过{...}
func skipCode() {
	var braces_num int = 1
	for {
		input, _, c := inputReader.ReadLine()
		if c == io.EOF {
			fmt.Println("if error")
			fmt.Println("missing '}' error.")
			os.Exit(1)
		}
		//fmt.Println(len(input), string(input), "kkk")

		if len(input) == 0 || strings.Replace(string(input), " ", "", -1) == "" { //跳过空行
			continue
		} else {

			line := string(input) + "\n"
			fmt.Println("sk", line)
			if strings.Contains(string(input), "}") {
				braces_num--
				if strings.Contains(string(input), "else") && braces_num == 0 {
					set_line([]rune(strings.TrimSpace(string(input)))[1:])
					fmt.Println("bak", string([]rune(strings.TrimSpace(string(input)))[1:]))
					stl_bak = 1
					break
				}
			}
			if strings.Contains(string(input), "{") {
				braces_num++
			}
			if braces_num == 0 {
				break
			}
		}
	}
}

func excutes() {
	var value interface{}

	paramList = make(map[string]Token) //变量列表
	// for {
	// 	inputReader := bufio.NewReader(os.Stdin)
	// 	fmt.Println("please input:")
	// 	input, err := inputReader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("There ware errors reading,exiting program.")
	// 		return
	// 	}
	// 	set_line([]rune(input))
	// 	value = parse_line()
	// 	fmt.Println(">>", reflect.ValueOf(value))
	// }

	for {
		//fmt.Println("please input:")
		if stl_bak == 0 {
			input, _, c := inputReader.ReadLine()
			if c == io.EOF {
				break
			}
			if len(input) == 0 { //跳过空行
				continue
			} else {
				line := string(input) + "\n"
				fmt.Println(line)
				set_line([]rune(line))
				value = parse_line()
				fmt.Println(">>", reflect.ValueOf(value))
			}
		} else {
			stl_bak = 0
			value = parse_line()
			fmt.Println(">>", reflect.ValueOf(value))
		}

	}
}

func parse_line() interface{} {
	var value interface{}

	st_look_ahead_token_exists = 0
	value = parse_logic_expression()
	return value
}

func main() {
	var err error
	//fi, err = os.Open("D:\\work\\TongJi\\go_work\\src\\src\\github.com\\fate\\mycalc2\\test.fate")
	fi, err = os.Open("D:\\0_chenyao\\git\\src\\fate\\mycalc2\\test.fate")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	inputReader = bufio.NewReader(fi)
	excutes()
}
