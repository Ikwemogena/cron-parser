# Deliveroo Cron Parser

## Overview
This is a Go-based cron expression parser deliveroo solution that parses standard cron syntax and outputs an expanded list of values for each field (minute, hour, day of month, month, and day of week). The program ensures valid input and provides detailed output.

## Features
- Parses standard cron expressions
- Handles wildcards (`*`), ranges (`-`), step values (`*/N`), and lists (`,`)
- Outputs formatted and structured results

## Installation & Build
Ensure you have Go installed, then clone the repository and build the project:

```sh
# Clone the repository
git clone https://github.com/Ikwemogena/cron-parser.git
cd cron-parser

# Build the program
go build -o deliveroo-cron-parser
```

## Usage
Run the compiled binary with a cron expression as an argument:

```sh
./deliveroo-cron-parser "*/15 0 1,15 * 1-5 /usr/bin/find"
```

## Example Output
```
minute         0 15 30 45
hour           0
day of month   1 15
month          1 2 3 4 5 6 7 8 9 10 11 12
day of week    1 2 3 4 5
command        /usr/bin/find
```

## Running Tests
To ensure correctness, run the included unit tests:

```sh
go test ./../deliveroo-cron-parser
```

## Author
Challenge by [Ikwemogena](https://github.com/Ikwemogena).


