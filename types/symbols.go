package types

import "git.woa.com/modnarshen/excelconfc/util"

const (
	L_SQUARE_BRACKET = "["
	R_SQUARE_BRACKET = "]"
	L_CURLY_BRACKET  = "{"
	R_CURLY_BRACKET  = "}"
)

type Mark = string
type Token = string

const (
	MARK_TYPE_BOOL   Mark = "bool"
	MARK_TYPE_INT32  Mark = "int32"
	MARK_TYPE_UINT32 Mark = "uint32"
	MARK_TYPE_INT64  Mark = "int64"
	MARK_TYPE_UINT64 Mark = "uint64"
	MARK_TYPE_STRING Mark = "string"

	MARK_DESC_VECTOR   Mark = "vector"
	MARK_DESC_ARRAY    Mark = "array"
	MARK_DESC_DATETIME Mark = "D"
	MARK_DESC_ENUM     Mark = "E"
	MARK_DESC_KEY      Mark = "K"

	MARK_VAL_TRUE  Mark = "true"
	MARK_VAL_FALSE Mark = "false"
)

const (
	TOK_NONE = ""

	TOK_LF_SQ_BRACKET Token = "[" // left square bracket
	TOK_RG_SQ_BRACKET Token = "]" // right square bracket
	TOK_LF_CR_BRACKET Token = "{" // left curly bracket
	TOK_RG_CR_BRACKET Token = "}" // right curly bracket

	TOK_TYPE_BOOL        Token = "bool"
	TOK_TYPE_INT32       Token = "int32"
	TOK_TYPE_UINT32      Token = "uint32"
	TOK_TYPE_INT64       Token = "int64"
	TOK_TYPE_UINT64      Token = "uint64"
	TOK_TYPE_STRING      Token = "string"
	TOK_TYPE_STRUCT      Token = "struct"
	TOK_TYPE_DATETIME    Token = "DateTime" // for golang
	TOK_TYPE_ENUM        Token = "@Enum"
	TOK_TYPE_VECTOR      Token = "@Vector"
	TOK_TYPE_VEC_STRUCT  Token = "@VecStruct"
	TOK_TYPE_ROOT_STRUCT Token = "@RootStruct"

	TOK_VAL_TRUE  Token = "true"
	TOK_VAL_FALSE Token = "false"

	MID_NODE_FIELDS        Token = "Node@FIELDS"
	MID_NODE_ADT           Token = "Node@ADT"
	MID_NODE_BDT           Token = "Node@BDT"
	MID_NODE_ARRAY         Token = "Node@ARRAY"
	MID_NODE_VEC           Token = "Node@VEC"
	MID_NODE_STRUCT        Token = "Node@STRUCT"
	MID_NODE_VEC_ADT_ITEMS Token = "Node@VEC_ADT_ITEMS"
	MID_NODE_VEC_BDT_ITEMS Token = "Node@VEC_BDT_ITEMS"
)

var (
	_BOOL_TYPES   = util.NewSet(TOK_TYPE_BOOL)
	_INT_TYPES    = util.NewSet(TOK_TYPE_INT32, TOK_TYPE_INT64, TOK_TYPE_UINT32, TOK_TYPE_UINT64)
	_STRING_TYPES = util.NewSet(TOK_TYPE_STRING)
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

func IsRealStruct(tok Token) bool {
	return tok == TOK_TYPE_STRUCT
}
