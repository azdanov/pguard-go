package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(0)

	var graceful bool
	flag.BoolVar(&graceful, "graceful", false, "gracefully terminate the process on timeout (sends SIGINT instead of SIGKILL)")

	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalf("Usage: %s [options] <timeout> <command> [args...]", os.Args[0])
	}

	timeout, err := time.ParseDuration(flag.Arg(0))
	if err != nil {
		log.Fatalf("invalid timeout: %v", err)
	}
	cmdName := flag.Arg(1)
	cmdArgs := flag.Args()[2:]

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		if err == context.DeadlineExceeded {
			return
		}
		log.Fatalf("failed to start command: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		defer close(done)
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Fatalf("command exited with error: %v", err)
		}
	case <-ch:
	case <-ctx.Done():
		var signal os.Signal
		if graceful {
			signal = os.Interrupt
		} else {
			signal = os.Kill
		}
		err := cmd.Process.Signal(signal)
		if err != nil {
			log.Fatalf("failed to send signal(%d): %v", signal, err)
		}
	}
}
