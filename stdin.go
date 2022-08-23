package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

//─────────────┤ StdinPipeOrOne ├─────────────

func StdinPipeOrOne(stdinPrompt string) ([]byte, bool) {
	// user enters a line of input at command line
	// or the input is read through stdin piping
	var piping bool
	var buf []byte
	o, _ := os.Stdin.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
		piping = false
	} else {
		piping = true
	}
	if piping {
		var data bytes.Buffer
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Fprintln(&data, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, false
		}
		buf = data.Bytes()
	} else { // get input from command prompt
		fmt.Print(stdinPrompt)
		inb := make([]byte, 256)
		n, _ := os.Stdin.Read(inb)
		if n > 1 {
			buf = bytes.Trim(inb, "\000\t \n")
		} else {
			return nil, false
		}
	}

	return buf, piping
}

//─────────────┤ StdinPipeOrMany ├─────────────

func StdinPipeOrMany(stdinPrompt ...string) ([][]byte, bool) {
	// user enters one or more lines of input at command line
	// or the input is read through stdin piping
	var piping bool
	var buf []byte
	var result [][]byte
	bufsz := 256

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
		piping = false
	} else {
		piping = true
	}
	if piping {
		var data bytes.Buffer
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Fprintln(&data, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, false
		}
		result = append(result, data.Bytes())
	} else { // get input from command prompt
		for _, pr := range stdinPrompt {
			fmt.Print(pr)
			inb := make([]byte, bufsz)
			n, _ := os.Stdin.Read(inb)
			if n > 1 {
				buf = bytes.Trim(inb, "\000\t \n")
				result = append(result, buf)
			} else {
				return nil, false
			}
		}
	}
	return result, piping
}
