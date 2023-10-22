package fs_resource_generator

import (
	"bufio"
	"log"
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

func (fs *FsResoureGenerator) Generate() <-chan domain.Resource {
	resourceChan := make(chan domain.Resource)

	go func() {
		defer close(resourceChan)
		file, err := os.Open(fs.Filename)
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
			argSlice := strings.Split(line, ",")
			row, _ := strconv.Atoi(argSlice[0])
			column, _ := strconv.Atoi(argSlice[1])
			color := argSlice[2]
			multiplier, _ := strconv.Atoi(argSlice[3])
			resource := domain.Resource{
				Position:   domain.Position{Row: int32(row), Column: int32(column)},
				Color:      domain.Color(color),
				Multiplier: uint8(multiplier),
			}
			resourceChan <- resource
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return resourceChan
}
