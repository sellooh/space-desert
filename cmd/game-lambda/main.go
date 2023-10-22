package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sellooh/space-desert/internal/adapters/fs_resource_generator"
	"github.com/sellooh/space-desert/internal/core/services"
)

type LambdaEvent struct {
	File string `json:"file"`
}

func handler(ctx context.Context, event LambdaEvent) (uint32, error) {
	calculateService := services.NewCalculateService()
	score := calculateService.Calculate(fs_resource_generator.NewFsResourceGenerator(event.File))

	log.Println("file", event.File)
	log.Println("score", score)
	return score, nil
}

func main() {
	lambda.Start(handler)
}
