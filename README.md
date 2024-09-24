# Autoeval

`autoeval` is a portable command line tool to perform black-box testing of any other command line tool.
The test scenarios are defined in `lua` scripts, which are executed by `autoeval` in a controlled environment.

`autoeval` is designed around some simple principles:

- **Portable**: it should work on any common platform
- **Single binary**: download and run, no dependencies
- **Versatile**: The use of `lua` scripts allows for a wide range of test scenarios

_This is a work in progress and not yet ready for production use. Expect **heavy** changes to the API._

## Installation

Download the latest release from the releases page and put it in your path or in the same directory as the executable you want to test.

Rename the executable to `autoeval` or `autoeval.exe` if you prefer.

## Usage

`autoeval` is a command line tool that takes a single argument: the path to the lua script that defines the test scenarios.
The test file extension should be `.test` or `.lua`.

```shell
autoeval test.test
```

## Defining test scenarios

The test scenarios are defined in a `lua` script that is executed by `autoeval`.
The available functions are:

- `launch(command [, args])`: Launches the command with the given arguments.
- `expect(output)`: Expects the output of the last command to be equal to the given string.
- `write(input)`: Writes the given string to the standard input of the running command.
- `stop()`: Stops the running command.

A single test file may contain multiple test scenarios.
Just call `launch` and `stop` in sequence to define the test scenario.


## Features

- [x] Launch any command
- [x] Expect exact output
- [x] Write to stdin
- [x] Stop the command
- [ ] Expect exit code
- [ ] Expect case insensitive output
- [ ] Expect output to match a regular expression
- [ ] Read files (currently hacked by `launch("cat", "file.txt")`)
- [ ] More versatile _current dir_ scenarios.

