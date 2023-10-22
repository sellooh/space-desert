package services

import (
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sellooh/space-desert/internal/core/domain"
	"github.com/sellooh/space-desert/internal/core/ports"
)

func calculateGroupScore(category string, In chan domain.Resource, Out chan domain.GroupScoreResult) {
	// find groups

	// during IN stream
	// - generate neighbors
	// - if new resource matches the HashMap
	// -- group match list with resource
	// implement for X (X) X && X X (X) --- when resource is in the middle or ends

	type resourceInGroup struct {
		domain.Resource
		group *mapset.Set[*resourceInGroup]
	}
	type groupSet = mapset.Set[*resourceInGroup]

	groups := mapset.NewSet[groupSet]()

	candidates := make(map[domain.Position]*resourceInGroup)

	total := uint32(0)
	for resource := range In {

		group := mapset.NewSet[*resourceInGroup]()
		groupPointer := &group
		resourceInGroup := resourceInGroup{
			Resource: resource,
		}
		candidates[domain.NewPosition(resourceInGroup.GetRow(), resourceInGroup.GetColumn())] = &resourceInGroup

		// map for neighbours groups
		neighbourGroups := mapset.NewSet[*groupSet]()

		// generate neighbours
		neighbours := domain.GenerateNeighbours(resource)

		for _, neighbour := range neighbours {
			if candidate, found := candidates[domain.NewPosition(neighbour.GetRow(), neighbour.GetColumn())]; found {
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
	groupPositions := [][]domain.Position{}
	for group := range groups.Iter() {

		groupLength := len(group.ToSlice())
		if groupLength >= 3 {
			positions := []domain.Position{}
			for resource := range group.Iter() {
				total += uint32(resource.Multiplier)
				positions = append(positions, (*resource).Position)
			}
			groupPositions = append(groupPositions, positions)
		}
	}
	Out <- domain.GroupScoreResult{category, total, groupPositions} // TODO: maybe this type should not be in the domain?
}

func calculateScoreFromBoard(board domain.Board) <-chan domain.BoardResult {
	resourceChan := make(chan domain.Resource, len(board))

	resultChannel := calculateScoreFromResourceChannel(resourceChan)
	go func() {
		defer close(resourceChan)
		for _, resource := range board {
			resourceChan <- resource
		}
	}()

	return resultChannel
}

func calculateScoreFromResourceChannel(resourceChan <-chan domain.Resource) <-chan domain.BoardResult {
	out := make(chan domain.BoardResult)
	colorChanMap := map[domain.Color]struct {
		In  chan domain.Resource
		Out chan domain.GroupScoreResult
	}{
		domain.LightBlue:   {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		domain.DarkBlue:    {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		domain.LightGreen:  {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		domain.DarkGreen:   {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		domain.LightPurple: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		domain.DarkPurple:  {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
	}
	multiplierChanMap := map[uint8]struct {
		In  chan domain.Resource
		Out chan domain.GroupScoreResult
	}{
		1: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		2: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		3: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		4: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		5: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
		6: {In: make(chan domain.Resource), Out: make(chan domain.GroupScoreResult, 1)},
	}
	var inChannels [12]chan domain.Resource
	var outChannels [12]chan domain.GroupScoreResult
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
			colorChanMap[resource.Color].In <- resource
			multiplierChanMap[resource.Multiplier].In <- resource
		}

		// No more resources to calculate
		for _, inChan := range inChannels {
			close(inChan)
		}

		totalScore := uint32(0)
		resultGroups := []domain.GroupScoreResult{}
		for _, outChannel := range outChannels {
			if groupScoreResult := <-outChannel; groupScoreResult.Score > 0 {
				totalScore += groupScoreResult.Score
				resultGroups = append(resultGroups, groupScoreResult)
			}
		}
		out <- domain.BoardResult{totalScore, resultGroups}
	}()

	return out
}

type CalculateService struct {
}

func NewCalculateService() CalculateService {
	return CalculateService{}
}

func (s CalculateService) Calculate(generator ports.ResourceGenerator) (uint32, error) {
	errorChannel := make(chan error, 1)
	resourceChannel := generator.Generate(errorChannel)
	resultChannel := calculateScoreFromResourceChannel(resourceChannel)
	select {
	case err := <-errorChannel:
		return 0, err
	case result := <-resultChannel:
		return result.Score, nil
	}
}
