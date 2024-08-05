package util

// 定义一个集合类型
type Set map[string]bool

// 初始化集合
func NewSet(elements ...string) Set {
	s := make(Set)
	for _, elem := range elements {
		s[elem] = true
	}
	return s
}
