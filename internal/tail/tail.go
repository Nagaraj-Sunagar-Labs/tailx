// Package tail contains the core logic for the tailx utility.
package tail

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Options holds the runtime options for a tail invocation.
type Options struct {
	// Lines is the number of trailing lines to output (-n).
	Lines int
	// Follow enables follow mode (-f).
	Follow bool
	// Files are the input file paths to read.
	Files []string
}

// ErrNotImplemented is returned by Run until the core logic is implemented.
var ErrNotImplemented = errors.New("tail logic not implemented yet")

// followDuration is deliberately fixed: -f follows appended data for one
// minute without exposing a timer flag.
const followDuration = 60 * time.Second

const followPollInterval = 250 * time.Millisecond

// Run executes the tail behaviour described by opts, writing output to w.
func Run(w io.Writer, opts Options) error {
	lines := opts.Lines
	if lines == 0 {
		lines = 5
	}
	files := opts.Files
	if len(files) == 0 {
		return errors.New("no input file was provided")
	}

	path := files[0]
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %q: %w", path, err)
	}
	defer file.Close()

	totallines, err := countTotalNoLines(file)
	if err != nil {
		return fmt.Errorf("count lines in %q: %w", path, err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind %q: %w", path, err)
	}

	reader := bufio.NewReader(file)
	startcount := totallines - lines
	inlinecount := 0
	for {
		line, readErr := reader.ReadString('\n')
		if len(line) > 0 {
			inlinecount++
			if inlinecount > startcount {
				if _, err := io.WriteString(w, line); err != nil {
					return fmt.Errorf("write output: %w", err)
				}
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return fmt.Errorf("read %q: %w", path, readErr)
		}
	}

	if opts.Follow {
		if err := follow(w, reader, path); err != nil {
			return err
		}
	}

	return nil
}

// follow writes data appended to the open file until the fixed follow window
// expires. Reading a regular file at its current end returns io.EOF, so the
// reader is polled for newly appended data.
func follow(w io.Writer, reader *bufio.Reader, path string) error {
	return followFor(w, reader, path, followDuration, followPollInterval)
}

// followFor contains the follow loop with explicit durations so its timing
// behaviour can be tested without waiting for the production one-minute limit.
func followFor(w io.Writer, reader *bufio.Reader, path string, duration, pollInterval time.Duration) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	timeout := time.NewTimer(duration)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			return nil
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					if _, writeErr := io.WriteString(w, line); writeErr != nil {
						return fmt.Errorf("write followed output: %w", writeErr)
					}
				}

				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("follow %q: %w", path, err)
				}
			}
		}
	}
}

func countTotalNoLines(input *os.File) (int, error) {
	scanner := bufio.NewScanner(input)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
