package util

const (
	indentSpaces   = "                                        " // len(indentSpaces) == 40
	indentSpaceNum = 4
)

func IndentSpace(indent int) string {
	if indent < 0 {
		return ""
	}
	return indentSpaces[:indent*indentSpaceNum]
}
