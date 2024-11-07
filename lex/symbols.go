package lex

import "git.woa.com/modnarshen/excelconfc/util"

const (
	L_SQUARE_BRACKET = "["
	R_SQUARE_BRACKET = "]"
	L_CURLY_BRACKET  = "{"
	R_CURLY_BRACKET  = "}"
)

type Mark = string
type Token = string
type LexMark = string

const (
	TOK_NONE = ""

	TOK_TYPE_BOOL        Token = "bool"
	TOK_TYPE_INT32       Token = "int32"
	TOK_TYPE_UINT32      Token = "uint32"
	TOK_TYPE_INT64       Token = "int64"
	TOK_TYPE_UINT64      Token = "uint64"
	TOK_TYPE_STRING      Token = "string"
	TOK_TYPE_FSTRING     Token = "FString"
	TOK_TYPE_FTEXT       Token = "FText"
	TOK_TYPE_STRUCT      Token = "struct"
	TOK_TYPE_DATETIME    Token = "DateTime" // for golang
	TOK_TYPE_ENUM        Token = "@Enum"
	TOK_TYPE_VECTOR      Token = "@Vector"
	TOK_TYPE_VEC_STRUCT  Token = "@VecStruct"
	TOK_TYPE_ROOT_STRUCT Token = "@RootStruct"

	TOK_VAL_TRUE  Token = "true"
	TOK_VAL_FALSE Token = "false"

	TOK_DESC_VECTOR   Token = "vector"
	TOK_DESC_ARRAY    Token = "array"
	TOK_DESC_DATETIME Token = "D"
	TOK_DESC_ENUM     Token = "E"
	TOK_DESC_KEY      Token = "K"

	MID_NODE_FIELDS        LexMark = "Node@FIELDS"
	MID_NODE_ADT           LexMark = "Node@ADT"
	MID_NODE_BDT           LexMark = "Node@BDT"
	MID_NODE_ARRAY         LexMark = "Node@ARRAY"
	MID_NODE_VEC           LexMark = "Node@VEC"
	MID_NODE_STRUCT        LexMark = "Node@STRUCT"
	MID_NODE_VEC_ADT_ITEMS LexMark = "Node@VEC_ADT_ITEMS"
	MID_NODE_VEC_BDT_ITEMS LexMark = "Node@VEC_BDT_ITEMS"

	LEX_BOOL   LexMark = "bool"
	LEX_ENUM   LexMark = "enum"
	LEX_ARRAY  LexMark = "array"
	LEX_INT    LexMark = "int"
	LEX_STRING LexMark = "string"
	LEX_ID     LexMark = "id"
)

var (
	_BOOL_TYPES   = util.NewSet(TOK_TYPE_BOOL)
	_INT_TYPES    = util.NewSet(TOK_TYPE_INT32, TOK_TYPE_INT64, TOK_TYPE_UINT32, TOK_TYPE_UINT64)
	_STRING_TYPES = util.NewSet(TOK_TYPE_STRING)

	_LEX_BDT = util.NewSet(LEX_BOOL, LEX_ENUM, LEX_INT, LEX_STRING)
)

func IsBasicType(tok Token) bool {
	return _BOOL_TYPES.Contains(tok) || _INT_TYPES.Contains(tok) || _STRING_TYPES.Contains(tok)
}

func IsIntType(tok Token) bool {
	return _INT_TYPES.Contains(tok)
}

func IsStringType(tok Token) bool {
	return _STRING_TYPES.Contains(tok)
}

func IsRepeatedLex(lexVal LexMark) bool {
	return lexVal == MID_NODE_VEC || lexVal == LEX_ARRAY
}
