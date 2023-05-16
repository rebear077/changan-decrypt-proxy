package structure

import (
	"sync"
)

type Set struct {
	m   map[interface{}]struct{}
	len int
	sync.RWMutex
}

func NewSet() *Set {
	temp := make(map[interface{}]struct{})
	return &Set{
		m: temp,
	}
}

// 增加一个元素
func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = struct{}{} // 实际往字典添加这个键
	s.len = len(s.m)       // 重新计算元素数量
}

// 移除一个元素
func (s *Set) Remove(item interface{}) {
	s.Lock()
	defer s.Unlock()
	// 集合没元素直接返回
	if s.len == 0 {
		return
	}
	delete(s.m, item) // 实际从字典删除这个键
	s.len = len(s.m)  // 重新计算元素数量
}

// 查看是否存在元素
func (s *Set) Has(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// 转化为列表
func (s *Set) List() interface{} {
	s.RLock()
	defer s.RUnlock()
	list := make([]interface{}, 0, s.len)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
