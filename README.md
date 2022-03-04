# worker

This repository holds the code for computing statistics from the results of a Benchttp run.

It is deployed as a piece of software independant from the other parts making our system.

The main function defined by `worker` is deployed and run as a stateless invocation inside Google Cloud Functions and is triggered by Cloud events.

## Invocation flow

- a user executes the `runner` process on their machine via `benchttp run` command
- the results are posted to Benchttp server and a `Report` is created inside Cloud Firestore
- the create event is dispatched to Google Cloud Functions and `worker.Digest` runs
- the results are saved inside a SQL persistence layer

```txt
    +--------+
    |        |
    | runner |
    |        |
    +--------+
        |
        |
        v
    +--------+     /-----------\               +--------+     /-----------\
    |        |     |           |               |        |     |           |
    | server | --> | firestore | ---trigger--> | worker | --> | sql layer |
    |        |     |           |               |        |     |           |
    +--------+     \-----------/               +--------+     \-----------/

```

## Deployment

The infrastructure code defining the deployment of `worker` is located inside [benchttp/infra](https://github.com/benchttp/infra).

## Development

### Lint

Run the linter:

> We use [golangci-lint](https://golangci-lint.run/) in our CI.

```sh
make lint
# alias to:
# golangci-lint run
```

### Test

Run all tests:

```sh
make test
# alias to:
# go test -v -timeout 30s ./...
```

Run a specific test passing `t` to specify a test and `p` to specify a package (parameters are independent):

```sh
make test t=TestCompute p=pkg/golastic
# alias to:
# go test -v -timeout 30s -run TestCompute ./stats
```
