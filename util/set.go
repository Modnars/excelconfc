package util

// Set 是一个泛型集合结构
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet 创建一个新的集合
func NewSet[T comparable](args ...T) *Set[T] {
	newSet := &Set[T]{elements: make(map[T]struct{})}
	for _, arg := range args {
		newSet.Add(arg)
	}
	return newSet
}

// Add 向集合中添加一个元素
func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

// Remove 从集合中移除一个元素
func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Contains 检查集合中是否包含某个元素
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size 返回集合的大小
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Elements 返回集合中的所有元素
func (s *Set[T]) Elements() []T {
	keys := make([]T, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	return keys
}
