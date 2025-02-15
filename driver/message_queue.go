package driver

type message struct {
	id   uint
	argc uint
	argv [3]uint
}

var queue = make([]message, 0)

func (m *message) getId() uint {
	return m.id
}

func (m *message) getArgc() uint {
	return m.argc
}

func (m *message) getArgv() [3]uint {
	return m.argv
}

func Enq(m message) {
	queue = append(queue, m)
}

func Deq() message {
	m := queue[0]
	queue = queue[1:]
	return m
}
