// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// FIFO data structure.
//
//=====================================================================================================================

package utils


type Queue struct {
	values	[]interface{}
	head	int
	tail	int
	size	int
	Count	int
}


func CreateQueue(size int) *Queue {
	return &Queue{
		values:		make([]interface{}, size),
		size:		size,
	}
}


func (q *Queue) Push(n interface{}) {
	if q.head == q.tail && q.Count > 0 {
		values := make([]interface{}, len(q.values) + q.size)
		copy(values, q.values[q.head:])
		copy(values[len(q.values)-q.head:], q.values[:q.head])
		q.head = 0
		q.tail = len(q.values)
		q.values = values
	}
	q.values[q.tail] = n
	q.tail = (q.tail + 1) % len(q.values)
	q.Count++
}


func (q *Queue) Pop() interface{} {
	if q.Count != 0 {
		value := q.values[q.head]
		q.head = (q.head + 1) % len(q.values)
		q.Count--
		return value
	}
	return nil
}