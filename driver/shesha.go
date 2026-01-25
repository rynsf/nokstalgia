package driver

import "log"

var (
	BOARD_WIDTH  = 23
	BOARD_HEIGHT = 13
	LEFT         = Cord{-1, 0}
	RIGHT        = Cord{1, 0}
	UP           = Cord{0, -1}
	DOWN         = Cord{0, 1}
)

var Directions = [4]Cord{UP, RIGHT, DOWN, LEFT}

var Hamilton298 = [13][23]int{
	{2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
	{2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 0, 3},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 0},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 3},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 0},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 3},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 0},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 3},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 0},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 3},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 0},
	{2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 3},
	{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0},
}

type Cord struct {
	x, y int
}

type State struct {
	Snake           []Segment
	Board           [][]bool
	Apple           Cord
	Dir             Cord
	Lost, Win, Done bool
	Mem             []byte
}

type Segment struct {
	pos Cord
}

func (game *State) GetHead() (int, int) {
	head := game.Snake[len(game.Snake)-1]
	return head.pos.x, head.pos.y
}

func NewSegment(pos Cord) Segment {
	return Segment{
		pos,
	}
}

func NewGame(mem []byte) State {
	board := [][]bool{}
	for range BOARD_HEIGHT {
		board = append(board, make([]bool, BOARD_WIDTH))
	}

	snake := make([]Segment, 7)
	snake[0] = NewSegment(Cord{5, 7})
	snake[1] = NewSegment(Cord{6, 7})
	snake[2] = NewSegment(Cord{7, 7})
	snake[3] = NewSegment(Cord{8, 7})
	snake[4] = NewSegment(Cord{9, 7})
	snake[5] = NewSegment(Cord{10, 7})
	snake[6] = NewSegment(Cord{11, 7})

	for _, s := range snake {
		board[s.pos.y][s.pos.x] = true
	}

	x := int(mem[0x20c4])
	y := int(mem[0x20c5])
	apple := Cord{x, y}

	return State{
		snake,
		board,
		apple,
		RIGHT,
		false,
		false,
		false,
		mem,
	}
}

func vectorAdd(a, b Cord) Cord {
	return Cord{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

func vectorEqual(a, b Cord) bool {
	if a.x == b.x && a.y == b.y {
		return true
	}
	return false
}

func (game *State) locateApple() {
	x := int(game.Mem[0x20c4])
	y := int(game.Mem[0x20c5])

	game.Apple = Cord{x, y}
}

func (game *State) dumpState() {
	log.Printf("lenght: %d\n", len(game.Snake))
	head := game.Snake[len(game.Snake)-1]
	log.Printf("head x:%d   y:%d\n", head.pos.x, head.pos.y)
}

func (game *State) Step() int {
	//game.dumpState()
	snakeHead := len(game.Snake) - 1
	snakeHeadCord := game.Snake[snakeHead]
	hDir := Hamilton298[snakeHeadCord.pos.y][snakeHeadCord.pos.x]

	game.locateApple()

	if len(game.Snake) >= 299 {
		game.Win = true
		game.Done = true
		return hDir
	}

	if vectorEqual(game.Apple, Cord{22, 0}) {
		Hamilton298[1][22] = 0
	} else {
		Hamilton298[1][22] = 3
	}

	if len(game.Snake) == 281 && snakeHeadCord.pos.x == 22 && snakeHeadCord.pos.y == 0 {
		Hamilton298[0][21] = 2
	}

	game.Dir = Directions[hDir]
	newHead := vectorAdd(game.Snake[snakeHead].pos, game.Dir)

	if newHead.x < 0 || newHead.y < 0 || newHead.x >= BOARD_WIDTH || newHead.y >= BOARD_HEIGHT {
		game.Lost = true
		game.Done = true
		return hDir
	}

	if vectorEqual(newHead, game.Apple) {
		game.Board[newHead.y][newHead.x] = true
		game.Snake = append(
			game.Snake,
			NewSegment(newHead),
		)
		game.locateApple()
		return hDir
	}

	tail := game.Snake[0].pos
	game.Board[tail.y][tail.x] = false

	if game.Board[newHead.y][newHead.x] == true {
		game.Lost = true
		game.Done = true
		return hDir
	}

	game.Board[newHead.y][newHead.x] = true

	for n := range snakeHead {
		game.Snake[n].pos = game.Snake[n+1].pos
	}

	game.Snake[snakeHead].pos = newHead

	return hDir
}
