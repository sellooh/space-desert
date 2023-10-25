I used this project to learn golang concurrency.

As the project evolved I've added:
* Created an Hexagonal project structure, also known as Ports and Adapters
* Setup sam serverless deployment
* Performed lambda ideal sizing with [Lambda Power Tools](https://github.com/alexcasalboni/aws-lambda-power-tuning/tree/master)

```sh
# run cli
$ go run cmd/game-cli/*.go data/1k-automata.txt

# sam build
$ sam build -t infrastructure/template.yaml

# sam local invoke
$ sam local invoke CalculateScoreFunction -e events/10k.json
```
