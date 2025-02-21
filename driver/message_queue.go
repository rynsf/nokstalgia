package driver

type message struct {
	id   uint32
	argc uint32
	argv [3]uint32
}

var queue = make([]message, 0)

func (m *message) getId() uint32 {
	return m.id
}

func (m *message) getArgc() uint32 {
	return m.argc
}

func (m *message) getArgv() [3]uint32 {
	return m.argv
}

func Enq(id, argc uint32, argv [3]uint32) {
	m := message{
		id,
		argc,
		argv,
	}
	queue = append(queue, m)
}

func Deq() (message, bool) {
	if len(queue) == 0 {
		return message{}, false
	}
	m := queue[0]
	queue = queue[1:]
	return m, true
}
