package attr

import (
	"strings"

	"golang.org/x/net/html"
)

type List []html.Attribute
type L = List

func (l L) Val(key string) (val string, ok bool) {
	for _, a := range l {
		if a.Key == key {
			return a.Val, true
		}
	}
	return
}

func (l L) ID() (string, bool) {
	return l.Val("id")
}

func (l L) ClassList() (classes []string, ok bool) {
	classStr, ok := l.Val("class")
	if !ok {
		return classes, false
	}
	return strings.Fields(classStr), ok
}

func (l L) Exists(key string) (ok bool) {
	for _, a := range l {
		if a.Key == key {
			return true
		}
	}
	return
}

func (l L) HasVal(key, val string) bool {
	foundVal, ok := l.Val(key)
	if !ok {
		return false
	}
	return val == foundVal
}

func (l L) HasWord(key, word string) bool {
	foundVal, ok := l.Val(key)
	if !ok {
		return false
	}
	for _, w := range strings.Fields(foundVal) {
		if w == word {
			return true
		}
	}
	return false
}

func (l L) HasID(id string) bool {
	return l.HasVal("id", id)
}

func (l L) HasClass(class string) bool {
	return l.HasWord("class", class)
}
