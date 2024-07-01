## Benchmarking same operations with and without channels

## ☔️ Code Coverage
**main**: &nbsp; [![coverage report](https://github.com/syedvasil/concurrent/blob/main/coverage.svg)](https://github.com/syedvasil/concurrent/blob/main/coverage.out)

## 🔧 Installation

After cloning the repository, navigate into the directory and make sure you are in **main** branch.

- Run `go mod tidy`, to download dependencies.

- Run `go run ./` to run this code on your local.
- Optionally, you check benchmark results by running `go test -bench=.` 
- Note: suggested to change `GOMAXPROCS` and run  multiple times



- Run `go test -coverprofile=coverage.out ./...` to create coverage report.
- Run `go tool cover -html="coverage.out"` to view the report on web browser.

In case Local go setup is not an option, other alternate is Docker

`docker build --progress=plain --no-cache -t concurrent .`
`docker run -p 8080:8080 concurrent`
