# pguard-go

A simple, cross-platform Go utility to execute commands with a specified timeout.
Similar to [timeout](https://man7.org/linux/man-pages/man1/timeout.1.html) command.

## Installation

Ensure you have [Go](https://go.dev/) installed. Clone the repository and build:

```bash
git clone https://github.com/azdanov/pguard-go.git
cd pguard-go
go build -o pguard-go .
```

## Usage

```bash
./pguard-go [options] <timeout> <command> [args...]
```

### Arguments

- `timeout`: Duration to wait before terminating the command (e.g., `5s`, `1m30s`).
- `command`: The executable to run.
- `args...`: Arguments to pass to the executable.

### Options

- `-graceful`: Gracefully terminate the process on timeout (sends SIGINT instead of SIGKILL).

### Examples

**Run a command with a 5-second timeout (killed forcefully if it exceeds):**

```bash
./pguard-go 5s sleep 10
```

**Run a ping and gracefully terminate it after 1 minute:**

```bash
./pguard-go -graceful 1m ping google.com
```

**Guard a command reading from standard input:**

```bash
cat large-file.txt | ./pguard-go 10s grep 'search term'
```
