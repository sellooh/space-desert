package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	mapset "github.com/deckarep/golang-set/v2"
)

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

func makeTransition(row, column int32, transition Transition) Position {
	return Position{row: row + int32(transition.row), column: column + int32(transition.column)}
}

// --------

type Position struct {
	row    int32
	column int32
}

type Resource struct {
	Position
	color      Color
	multiplier uint8
}

type Board []Resource

func generateNeighbours(resource Resource) [6]Position {
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
	category       string
	score          uint32
	positionMatrix [][]Position
}

func calculateGroupScore(category string, In chan Resource, Out chan GroupScoreResult) {
	// find groups

	// during IN stream
	// - generate neighbors
	// - if new resource matches the HashMap
	// -- group match list with resource
	// implement for X (X) X && X X (X) --- when resource is in the middle or ends

	type resourceInGroup struct {
		Resource
		group *mapset.Set[*resourceInGroup]
	}
	type groupSet = mapset.Set[*resourceInGroup]

	groups := mapset.NewSet[groupSet]()

	candidates := make(map[Position]*resourceInGroup)

	total := uint32(0)
	for resource := range In {

		group := mapset.NewSet[*resourceInGroup]()
		groupPointer := &group
		resourceInGroup := resourceInGroup{
			Resource: resource,
			// group:    groupPointer,
		}
		candidates[Position{resourceInGroup.row, resourceInGroup.column}] = &resourceInGroup

		// map for neighbours groups
		neighbourGroups := mapset.NewSet[*groupSet]()

		// generate neighbours
		neighbours := generateNeighbours(resource)

		for _, neighbour := range neighbours {
			if candidate, found := candidates[Position{neighbour.row, neighbour.column}]; found {
				// merges the resources in a group
				*groupPointer = group.Union(*candidate.group)
				groups.Remove(*candidate.group)
				neighbourGroups.Add(candidate.group)
			}
		}
		neighbourGroupsSlice := neighbourGroups.ToSlice()
		singleNeighbourGroup, _ := neighbourGroups.Pop()

		if len(neighbourGroupsSlice) == 1 {
			// fmt.Println("linear add (ngs = 1)")
			(*singleNeighbourGroup).Add(&resourceInGroup)
			resourceInGroup.group = singleNeighbourGroup
			groups.Add(*singleNeighbourGroup)
		} else if len(neighbourGroupsSlice) == 0 {
			// fmt.Println("new group (ngs = 0)")
			resourceInGroup.group = groupPointer
			group.Add(&resourceInGroup)
			groups.Add(*groupPointer)
		} else if len(neighbourGroupsSlice) > 1 {
			// fmt.Println("groups merge (ngs > 1)")
			group.Add(&resourceInGroup)
			groups.Add(*groupPointer)
			for resource := range group.Iter() {
				resource.group = groupPointer
			}
		}

	}
	groupPositions := [][]Position{}
	for group := range groups.Iter() {

		groupLength := len(group.ToSlice())
		if groupLength >= 3 {
			positions := []Position{}
			for resource := range group.Iter() {
				total += uint32(resource.multiplier)
				positions = append(positions, (*resource).Position)
			}
			groupPositions = append(groupPositions, positions)
		}
	}
	Out <- GroupScoreResult{category, total, groupPositions}
}

type BoardResult struct {
	score             uint32
	groupScoreResults []GroupScoreResult
}

func calculateScoreFromBoard(board Board) <-chan BoardResult {
	resourceChan := make(chan Resource, len(board))

	resultChannel := calculateScoreFromResourceChannel(resourceChan)
	go func() {
		defer close(resourceChan)
		for _, resource := range board {
			resourceChan <- resource
		}
	}()

	return resultChannel
}

func calculateScoreFromResourceChannel(resourceChan <-chan Resource) <-chan BoardResult {
	out := make(chan BoardResult)
	colorChanMap := map[Color]struct {
		In  chan Resource
		Out chan GroupScoreResult
	}{
		LightBlue:   {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		DarkBlue:    {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		LightGreen:  {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		DarkGreen:   {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		LightPurple: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		DarkPurple:  {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
	}
	multiplierChanMap := map[uint8]struct {
		In  chan Resource
		Out chan GroupScoreResult
	}{
		1: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		2: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		3: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		4: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		5: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
		6: {In: make(chan Resource), Out: make(chan GroupScoreResult, 1)},
	}
	var inChannels [12]chan Resource
	var outChannels [12]chan GroupScoreResult
	index := 0

	for category, channels := range colorChanMap {
		inChannels[index] = channels.In
		outChannels[index] = channels.Out
		go calculateGroupScore(string(category), channels.In, channels.Out)
		index++
	}
	for multiplier, channels := range multiplierChanMap {
		inChannels[index] = channels.In
		outChannels[index] = channels.Out
		go calculateGroupScore(strconv.Itoa(int(multiplier)), channels.In, channels.Out)
		index++
	}

	go func() {
		defer close(out)
		for resource := range resourceChan {
			colorChanMap[resource.color].In <- resource
			multiplierChanMap[resource.multiplier].In <- resource
		}

		// No more resources to calculate
		for _, inChan := range inChannels {
			close(inChan)
		}

		totalScore := uint32(0)
		resultGroups := []GroupScoreResult{}
		for _, outChannel := range outChannels {
			if groupScoreResult := <-outChannel; groupScoreResult.score > 0 {
				totalScore += groupScoreResult.score
				resultGroups = append(resultGroups, groupScoreResult)
			}
		}
		out <- BoardResult{totalScore, resultGroups}
	}()

	return out
}

func calculate(filePath string) uint32 {
	resourceChan := make(chan Resource)

	resultChannel := calculateScoreFromResourceChannel(resourceChan)
	go func() {
		defer close(resourceChan)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			// fmt.Println(line)
			argSlice := strings.Split(line, ",")
			row, _ := strconv.Atoi(argSlice[0])
			column, _ := strconv.Atoi(argSlice[1])
			color := argSlice[2]
			multiplier, _ := strconv.Atoi(argSlice[3])
			resource := Resource{
				Position:   Position{row: int32(row), column: int32(column)},
				color:      Color(color),
				multiplier: uint8(multiplier),
			}
			resourceChan <- resource
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	result := <-resultChannel
	return result.score
}

type LambdaEvent struct {
	File string `json:"file"`
}

func handler(ctx context.Context, event LambdaEvent) (uint32, error) {
	score := calculate(event.File)

	log.Println("file", event.File)
	log.Println("score", score)
	return score, nil
}

func main() {
	lambda.Start(handler)
}
