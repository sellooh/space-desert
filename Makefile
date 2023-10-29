build-CalculateScoreFunction:
	GOARCH=arm64 GOOS=linux go build -o ./bootstrap cmd/game-lambda/*.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.