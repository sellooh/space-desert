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

# sam deploy (simplified)
$ sam deploy --config-file infrastructure/samconfig.toml

# sam deploy (with parameters)
$ sam deploy --stack-name space-desert -t infrastructure/template.yaml --parameter-overrides MountDataLayer=false

# sam invoke
$ sam remote invoke LAMBDA_ARN --event-file events/10k.json
```

##### sam local invoke

sam local doesn't evaluate cloudformation conditionals. Avoid using It [issue](https://github.com/aws/aws-sam-cli/issues/194).

