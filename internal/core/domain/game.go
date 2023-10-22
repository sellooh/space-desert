package domain

type Color string

const (
	DarkBlue    Color = "DB"
	LightBlue   Color = "LB"
	DarkGreen   Color = "DG"
	LightGreen  Color = "LG"
	DarkPurple  Color = "DP"
	LightPurple Color = "LP"
)

type ShiftDirection string

const (
	L  ShiftDirection = "L"
	R  ShiftDirection = "R"
	TL ShiftDirection = "TL"
	TR ShiftDirection = "TR"
	BL ShiftDirection = "BL"
	BR ShiftDirection = "BR"
)

type Transition struct {
	row    int8
	column int8
}

type TransitionMap map[ShiftDirection]Transition

var shiftMap = TransitionMap{
	L:  {row: 0, column: -1},
	R:  {row: 0, column: 1},
	TR: {row: 1, column: 1},
	TL: {row: 1, column: 0},
	BR: {row: -1, column: 0},
	BL: {row: -1, column: -1},
}

func MakeTransition(row, column int32, transition Transition) Position {
	return Position{Row: row + int32(transition.row), Column: column + int32(transition.column)}
}

// --------

type Position struct {
	Row    int32
	Column int32
}

type Resource struct {
	Position
	Color      Color
	Multiplier uint8
}

type Board []Resource

func GenerateNeighbours(resource Resource) [6]Position {
	return [6]Position{
		MakeTransition(resource.Row, resource.Column, shiftMap[L]),
		MakeTransition(resource.Row, resource.Column, shiftMap[TL]),
		MakeTransition(resource.Row, resource.Column, shiftMap[TR]),
		MakeTransition(resource.Row, resource.Column, shiftMap[R]),
		MakeTransition(resource.Row, resource.Column, shiftMap[BR]),
		MakeTransition(resource.Row, resource.Column, shiftMap[BL]),
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
