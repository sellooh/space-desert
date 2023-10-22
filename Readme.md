```sh
# run cli
$ go run cmd/game-cli/*.go

# sam build
$ sam build -t infrastructure/template.yaml

# sam local invoke
$ sam local invoke CalculateScoreFunction -e events/10k.json
```
