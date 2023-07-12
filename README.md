# gomake-omar

Gomake is a command-line tool designed to check for cyclic dependencies in a Makefile and run a specific target along with its dependencies. This tool helps ensure the smooth execution of Makefile targets by preventing cyclic dependencies and facilitating the efficient building of complex projects.

## Features

- Detects cyclic dependencies: Gomake analyzes the provided Makefile and checks for any cyclic dependencies among the targets.

- Executes specific targets: With Gomake, you can specify a target to execute. Gomake will ensure that all the target's dependencies are executed first, following the correct order of execution specified in the Makefile.

- Customizable Makefile path: The `-f` flag allows you to specify the path to the Makefile you want to analyze and execute. If not provided, Gomake will assume that the Makefile is present in the same directory and named "Makefile."

## Installation

To install Gomake, follow these steps:

1. Clone the Gomake repository to your local machine:

   ```shell
   $ git clone https://github.com/codescalersinternships/gomake-omar
   $ cd gomake-omar
   ```
2. Build the Gomake binary using the Go compiler:
   ```shell
   $ go build -o gomake ./cmd/make.go
   ```

## Usage

To use Gomake, open your terminal or command prompt and run the following command:
   ```shell
   $ gomake -t <target> [-f <filepath>]
   ```

- `-t <target>`: Specifies the target you want to execute. This flag is mandatory and must be provided.

- `-f <filepath>`: Specifies the path to the Makefile you want to analyze and execute. This flag is optional. If not provided, Gomake assumes that the Makefile is located in the same directory and named "Makefile".

## Examples

Here are some examples demonstrating the usage of Gomake:

1. Run the target `build` with its dependencies:
```shell
$ gomake -t build
```
2. Run the target `test` with its dependencies, using a specific Makefile located at `/path/to/custom/Makefile`:
```shell
$ gomake -t test -f /path/to/custom/Makefile
```

