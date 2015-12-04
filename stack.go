package main

// Stack is a FIFO for ints
type Stack struct {
	s []parseState
}

// NewStack creates a stack
func NewStack() *Stack {
	return &Stack{make([]parseState, 0)}
}

// Push an int onto the stack
func (s *Stack) Push(v parseState) {
	s.s = append(s.s, v)
}

// Pop an element from the stack
func (s *Stack) Pop() parseState {
	l := len(s.s)
	if l == 0 {
		panic("Empty stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res
}
