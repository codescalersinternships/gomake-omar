# gomake

Gomake is a command-line tool designed to check for cyclic dependencies in a Makefile and run a specific target along with its dependencies. This tool helps ensure the smooth execution of Makefile targets by preventing cyclic dependencies and facilitating the efficient building of complex projects.

## Features

- Detects cyclic dependencies: Gomake analyzes the provided Makefile and checks for any cyclic dependencies among the targets.

- Executes specific targets: With Gomake, you can specify a target to execute. Gomake will ensure that all the target's dependencies are executed first, following the correct order of execution specified in the Makefile.

- Customizable Makefile path: The `-f` flag allows you to specify the path to the Makefile you want to analyze and execute. If not provided, Gomake will assume that the Makefile is present in the same directory and named "Makefile."

## Installation

Before proceeding with the installation, ensure you have Go programming language (version 1.16 or higher) on your system, you can download and install Go from the [official Go website](https://golang.org)

### option1: Download Pre-built Binary
This option allows you to download and use the pre-built binary for your operating system, which simplifies the installation process.

1. Download the Binary, go to the GitHub [Releases page](hhttps://github.com/codescalersinternships/gomake-omar/releases) and download the binary that corresponds to your operating system and architecture. For example, for Linux 64-bit, you may find a file like `gomake-omar_Linux_arm64.tar.gz`.

2. Extract the Binary, after downloading the binary archive, extract its contents to a directory of your choice.
```bash
tar -xzf gomake-omar_Linux_arm64.tar.gz
```

3. Run the gomake app, see [Usage section](#usage)

### option2: Clone and Build from Source
This option is suitable if you want to build the Go Console App from the source code.

1. Clone the Gomake repository to your local machine:
```shell
$ git clone https://github.com/codescalersinternships/gomake-omar
$ cd gomake-omar
```

2. Build the Gomake binary using the Go compiler:
```shell
$ go build -o gomake ./cmd/make.go
```

3. Run the gomake app, see [Usage section](#usage)

## Usage

To use Gomake, open your terminal or command prompt and run the following command:
```shell
$ gomake -t <target> -f <filepath>
```

- `-t <target>`: Specifies the target you want to execute. This flag is mandatory and must be provided.

- `-f <filepath>`: Specifies the path to the Makefile you want to analyze and execute. This flag is optional. If not provided, Gomake assumes that the Makefile is located in the same directory and named "Makefile".

## Examples

Here are some examples demonstrating the usage of Gomake

```
# Default Makefile example for a Go project

all: build test

build:
	go build -o ./bin/myapp main.go

test:
	go test ./...

clean:
	rm -rf ./bin
```


1. Run the target `all` with its dependencies, assume makefile exists in the same directory with name 'Makefile':
```shell
$ gomake -t all
```
2. Run the target `build` with its dependencies, using a specific Makefile located at `/path/to/custom/Makefile`:
```shell
$ gomake -t build -f /path/to/custom/Makefile
```
## Test

To run the automated tests for this project, follow these steps:

1. Install the necessary dependencies by running `go get -d ./...`.
2. Run the tests by running `go test ./...`.
3. If all tests pass, the output should indicate that the tests have passed. If any tests fail, the output will provide information on which tests failed.
