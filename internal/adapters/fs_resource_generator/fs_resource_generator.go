package fs_resource_generator

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/sellooh/space-desert/internal/core/domain"
)

type FsResoureGenerator struct {
	Filename string
}

func NewFsResourceGenerator(filename string) *FsResoureGenerator {
	return &FsResoureGenerator{
		filename,
	}
}

func (fs *FsResoureGenerator) Generate(errorChannel chan<- error) <-chan domain.Resource {
	resourceChan := make(chan domain.Resource)

	go func() {
		defer close(resourceChan)
		file, err := os.Open(fs.Filename)
		if err != nil {
			errorChannel <- err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			argSlice := strings.Split(line, ",")
			row, err := strconv.Atoi(argSlice[0])
			if err != nil {
				errorChannel <- err
				continue
			}
			column, err := strconv.Atoi(argSlice[1])
			if err != nil {
				errorChannel <- err
				continue
			}
			color := argSlice[2]
			multiplier, err := strconv.Atoi(argSlice[3])
			if err != nil {
				errorChannel <- err
				continue
			}
			resource, resourceError := domain.NewResource(int32(row), int32(column), domain.Color(color), uint8(multiplier))
			if resourceError != nil {
				errorChannel <- err
				continue
			}
			resourceChan <- resource
		}

		if err := scanner.Err(); err != nil {
			errorChannel <- err
		}
	}()

	return resourceChan
}
