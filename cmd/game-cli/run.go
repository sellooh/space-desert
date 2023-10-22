package main

import (
	"log"
	"os"

	"github.com/sellooh/space-desert/internal/adapters/fs_resource_generator"
	"github.com/sellooh/space-desert/internal/core/services"
)

func run(args []string, stdout *os.File) error {
	calculateService := services.NewCalculateService()

	score := calculateService.Calculate(fs_resource_generator.NewFsResourceGenerator(args[1]))

	log.Println("score", score)
	return nil
}
