package cache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}
	l.pushFront(item)
	return item
}

func (l *list) pushFront(item *ListItem) {
	item.Next = l.front
	if l.front != nil {
		l.front.Prev = item
	}
	l.front = item
	if l.back == nil {
		l.back = item
	}
	l.len++
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	if l.back != nil {
		l.back.Next = item
	}
	l.back = item
	if l.front == nil {
		l.front = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFront(i)
}
