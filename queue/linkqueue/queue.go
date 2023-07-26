package linkqueue

import (
	"errors"
	"sync"
)

type LinkElement struct {
	item interface{}
	next *LinkElement
}

type LinkQueue struct {
	lock   sync.Mutex
	head   *LinkElement
	tail   *LinkElement
	length int
}

func (t *LinkQueue) GetLength() int {
	return t.length
}

func (t *LinkQueue) Put(item interface{}) error {
	if t.head != nil {
		if t.tail == nil {
			return errors.New("tail == nil")
		}
		if t.tail.next != nil {
			return errors.New("tail.next != nil")
		}
	}

	elem := &LinkElement{
		item: item,
	}

	if t.head == nil {
		t.head = elem
		t.tail = elem
		t.length = 1
		return nil
	}

	tailElem := t.tail
	tailElem.next = elem
	t.tail = elem
	t.length += 1
	return nil
}

func (t *LinkQueue) PutToHead(item interface{}) error {
	elem := &LinkElement{
		item: item,
	}

	if t.head == nil {
		t.head = elem
		t.tail = elem
		t.length = 1
		return nil
	}

	previousHead := t.head
	t.head = elem
	t.head.next = previousHead
	t.length += 1

	return nil
}

func (t *LinkQueue) Puts(items ...interface{}) (int, error) {
	c := 0
	for n, item := range items {
		if err := t.Put(item); err != nil {
			return c, err
		}
		c = n + 1
	}

	return c, nil
}

func (t *LinkQueue) PutsToHead(items ...interface{}) (int, error) {
	c := 0
	for n, item := range items {
		if err := t.Put(item); err != nil {
			return c, err
		}
		c = n + 1
	}

	return c, nil
}

func (t *LinkQueue) SafetyPut(item interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Put(item)
}

func (t *LinkQueue) SafetyPutToHead(item interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.PutToHead(item)
}

func (t *LinkQueue) SafetyPuts(items ...interface{}) (int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Puts(items...)
}

func (t *LinkQueue) SafetyPutsToHead(items ...interface{}) (int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.PutsToHead(items...)
}

func (t *LinkQueue) Pop() (interface{}, bool, error) {
	if t.length == 0 {
		return nil, false, nil
	}

	headElem := t.head
	if headElem == nil {
		return nil, false, errors.New("LinkQueue length != 0, head node == nil")
	}

	t.head = headElem.next
	retItem := headElem.item
	headElem.next = nil
	headElem.item = nil
	headElem = nil

	if t.length == 1 {
		t.tail = nil
	}

	t.length -= 1
	return retItem, true, nil
}

func (t *LinkQueue) Pops(count int) (retList []interface{}, retErr error) {
	for i := 0; i < count; i++ {
		item, ok, err := t.Pop()
		if err != nil {
			retErr = err
			return
		}
		if !ok {
			return
		}
		retList = append(retList, item)
	}

	return
}

func (t *LinkQueue) SafetyPop() (interface{}, bool, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Pop()
}

func (t *LinkQueue) SafetyPops(count int) ([]interface{}, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Pops(count)
}
