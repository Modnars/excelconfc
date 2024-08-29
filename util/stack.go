package util

// Stack 是一个泛型栈结构
type Stack[T any] struct {
	elements []T
}

// NewStack 创建一个新的栈
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{elements: []T{}}
}

// Push 向栈中添加一个元素
func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

// Pop 从栈中移除并返回顶部元素
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element, true
}

// Peek 返回栈顶元素但不移除它
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

// PeekOrNil 返回栈顶元素但不移除它，如果是空栈，返回一个零值
func (s *Stack[T]) PeekOrZero() T {
	if s.IsEmpty() {
		var zero T
		return zero
	}
	return s.elements[len(s.elements)-1]
}

// IsEmpty 检查栈是否为空
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// Size 返回栈的大小
func (s *Stack[T]) Size() int {
	return len(s.elements)
}
