package ds

type (
	//Queue 队列
	Queue struct {
		top    *node
		rear   *node
		length int
	}
	//双向链表节点
	node struct {
		pre   *node
		next  *node
		value any
	}
)

// NewQueue Create a new queue
func NewQueue() *Queue {
	return &Queue{nil, nil, 0}
}

// Len 获取队列长度
func (q *Queue) Len() int {
	return q.length
}

// Any 返回true队列不为空
func (q *Queue) Any() bool {
	return q.length > 0
}

// Peek 返回队列顶端元素
func (q *Queue) Peek() any {
	if q.top == nil {
		return nil
	}
	return q.top.value
}

// Rear 返回队列尾端元素
func (q *Queue) Rear() any {
	if q.rear == nil {
		return nil
	}
	return q.rear.value
}

// Push 入队操作
func (q *Queue) Push(v any) {
	n := &node{nil, nil, v}
	if q.length == 0 {
		q.top = n
		q.rear = q.top
	} else {
		n.pre = q.rear
		if q.rear != nil {
			q.rear.next = n
		}
		q.rear = n
	}
	q.length++
}

// Pop 出队操作
func (q *Queue) Pop() any {
	if q.length == 0 {
		return nil
	}
	n := q.top
	if q.top.next == nil {
		q.top = nil
	} else {
		q.top = q.top.next
		q.top.pre.next = nil
		q.top.pre = nil
	}
	q.length--
	return n.value
}

// RearRange 末尾遍历操作
func (q *Queue) RearRange(max int, handler func(item any)) {
	if q.length == 0 {
		return
	}
	if q.length < max {
		max = q.length
	}
	total := 0
	rear := q.rear
	for total < max && rear != nil {
		handler(rear.value)
		rear = rear.pre
		total += 1
	}
}

// Range 从头开始遍历
func (q *Queue) Range(handler func(item any)) {
	if q.length == 0 {
		return
	}
	curr := q.top
	for curr != nil {
		handler(curr.value)
		curr = curr.next
	}
}

// RangePop 从头开始遍历
func (q *Queue) RangePop(handler func(item any) bool) {
	if q.length == 0 {
		return
	}
	curr := q.top
	for curr != nil {
		pop := handler(curr.value)
		if pop {
			pre := curr.pre
			next := curr.next
			if pre != nil {
				pre.next = next
				curr.pre = nil
			} else {
				q.top = next
			}

			if next != nil {
				next.pre = pre
				curr.next = nil
			} else {
				q.rear = pre
			}
			curr = next
			q.length--
		} else {
			curr = curr.next
		}
	}
}

// RangePopMax 从头开始遍历,限制最大数量
func (q *Queue) RangePopMax(max int, handler func(item any) bool) {
	if q.length == 0 {
		return
	}
	total := 0
	curr := q.top
	for total < max && curr != nil {
		pop := handler(curr.value)
		if pop {
			total += 1
			pre := curr.pre
			next := curr.next
			if pre != nil {
				pre.next = next
				curr.pre = nil
			} else {
				q.top = next
			}

			if next != nil {
				next.pre = pre
				curr.next = nil
			} else {
				q.rear = pre
			}
			curr = next
			q.length--
		} else {
			curr = curr.next
		}
	}
}
