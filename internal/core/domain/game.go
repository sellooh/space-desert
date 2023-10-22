package domain

import "errors"

// -------- errors --------

var (
	ErrInvalidColor      = errors.New("invalid color")
	ErrInvalidMultiplier = errors.New("invalid multiplier")
)

// -------- constants --------

type Color string

const (
	DarkBlue    Color = "DB"
	LightBlue   Color = "LB"
	DarkGreen   Color = "DG"
	LightGreen  Color = "LG"
	DarkPurple  Color = "DP"
	LightPurple Color = "LP"
)

type shiftDirection string

const (
	L  shiftDirection = "L"
	R  shiftDirection = "R"
	TL shiftDirection = "TL"
	TR shiftDirection = "TR"
	BL shiftDirection = "BL"
	BR shiftDirection = "BR"
)

type direction struct {
	row    int8
	column int8
}

type transition map[shiftDirection]direction

var shiftMap = transition{
	L:  direction{row: 0, column: -1},
	R:  direction{row: 0, column: 1},
	TR: direction{row: 1, column: 1},
	TL: direction{row: 1, column: 0},
	BR: direction{row: -1, column: 0},
	BL: direction{row: -1, column: -1},
}

// -------- private functions --------

func makeTransition(row, column int32, transition direction) Position {
	return NewPosition(row+int32(transition.row), column+int32(transition.column))
}

func isValidColor(color Color) bool {
	return color == DarkBlue || color == LightBlue || color == DarkGreen || color == LightGreen || color == DarkPurple || color == LightPurple
}

func isValidMultiplier(multiplier uint8) bool {
	return multiplier >= 1 && multiplier <= 6
}

// -------- public functions --------

type Position struct {
	row    int32
	column int32
}

func (p Position) GetRow() int32 {
	return p.row
}

func (p Position) GetColumn() int32 {
	return p.column
}

func NewPosition(row int32, column int32) Position {
	return Position{row: row, column: column}
}

type Resource struct {
	Position
	Color      Color
	Multiplier uint8
}

func NewResource(x int32, y int32, color Color, multiplier uint8) (Resource, error) {
	if !isValidColor(color) {
		return Resource{}, ErrInvalidColor
	}
	if !isValidMultiplier(multiplier) {
		return Resource{}, ErrInvalidMultiplier
	}
	return Resource{
		Position:   NewPosition(x, y),
		Color:      color,
		Multiplier: multiplier,
	}, nil
}

type Board []Resource

func GenerateNeighbours(resource Resource) [6]Position {
	return [6]Position{
		makeTransition(resource.row, resource.column, shiftMap[L]),
		makeTransition(resource.row, resource.column, shiftMap[TL]),
		makeTransition(resource.row, resource.column, shiftMap[TR]),
		makeTransition(resource.row, resource.column, shiftMap[R]),
		makeTransition(resource.row, resource.column, shiftMap[BR]),
		makeTransition(resource.row, resource.column, shiftMap[BL]),
	}
}

type GroupScoreResult struct {
	Category       string
	Score          uint32
	PositionMatrix [][]Position
}

type BoardResult struct {
	Score             uint32
	GroupScoreResults []GroupScoreResult
}
