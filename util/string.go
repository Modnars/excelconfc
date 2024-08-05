package util

import "strings"

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
