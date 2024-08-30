package util

import "strings"

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

func SnakeToPascal(snakeStr string) string {
	// 将字符串按下划线分割
	parts := strings.Split(snakeStr, "_")
	// 处理每个部分
	for i := range parts {
		if len(parts[i]) > 0 {
			// 将每个部分的首字母大写，其余部分保持不变
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}
	// 将所有部分连接起来
	pascalStr := strings.Join(parts, "")
	return pascalStr
}

func GetPackageName(pkgDecl string) string {
	splitCh := ';'
	index := 0
	if strings.ContainsRune(pkgDecl, splitCh) {
		index = strings.IndexRune(pkgDecl, splitCh) + 1
	} else {
		index = strings.IndexRune(pkgDecl, '/') + 1
	}
	return pkgDecl[index:]
}
