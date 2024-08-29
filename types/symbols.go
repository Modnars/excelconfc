package types

import "git.woa.com/modnarshen/excelconfc/util"

const (
	L_SQUARE_BRACKET = "["
	R_SQUARE_BRACKET = "]"
	L_CURLY_BRACKET  = "{"
	R_CURLY_BRACKET  = "}"
)

type Token = string

const (
	TOK_LF_SQ_BRACKET Token = "[" // left square bracket
	TOK_RG_SQ_BRACKET Token = "]" // right square bracket
	TOK_LF_CR_BRACKET Token = "{" // left curly bracket
	TOK_RG_CR_BRACKET Token = "}" // right curly bracket

	TOK_TYPE_BOOL   Token = "bool"
	TOK_TYPE_INT32  Token = "int32"
	TOK_TYPE_UINT32 Token = "uint32"
	TOK_TYPE_INT64  Token = "int64"
	TOK_TYPE_UINT64 Token = "uint64"
	TOK_TYPE_STRING Token = "string"

	TOK_DESC_VECTOR Token = "vector"
	TOK_DESC_ARRAY  Token = "array"
	TOK_DESC_STRUCT Token = "struct"
	TOK_DESC_ENUM   Token = "E"

	TOK_PARSED_TYPE_DATETIME Token = "DateTime"
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
