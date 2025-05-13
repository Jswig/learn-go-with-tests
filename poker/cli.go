package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	input *bufio.Scanner
}

func (cli *CLI) PlayPoker() {
	cli.input.Scan()
	name := strings.Replace(cli.input.Text(), " wins", "", 1)
	cli.store.RecordWin(name)
}

func NewCLI(store PlayerStore, input io.Reader) *CLI {
	return &CLI{
		store: store,
		input: bufio.NewScanner(input),
	}
}
