package types

import "testing"

type A struct {
	name string
	age uint64
}

type WA A

type B struct {
	wa *WA
	score uint64
}

func TestA(t *testing.T) {
	a := &A{
		name: "xiaming",
		age: 14,
	}
	b := &B{}
	b.wa = (*WA)(a)
	t.Log(b.wa.name)
	t.Log(b.wa.age)
}
