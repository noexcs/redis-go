package datastruct

import (
	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type List struct {
	list *doublylinkedlist.List
}

func (l *List) PushLeft(element string) {
	l.list.Insert(0, element)
}

func (l *List) Size() int64 {
	return int64(l.list.Size())
}

func (l *List) PushRight(element string) {
	l.list.Append(element)
}

func (l *List) PopLeft() (string, bool) {
	v, exist := l.list.Get(0)
	if exist {
		l.list.Remove(0)
		return v.(string), exist
	}
	return "", exist
}

func (l *List) PopRight() (string, bool) {
	lastIndex := l.list.Size() - 1
	v, exist := l.list.Get(lastIndex)
	if exist {
		l.list.Remove(lastIndex)
		return v.(string), exist
	}
	return "", exist
}

func NewList() *List {
	return &List{list: doublylinkedlist.New()}
}
