package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Game interface {
	Start(numPlayers int)
	Finish(winner string)
}

type TexasHoldEm struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewTexasHoldEm(store PlayerStore, alerter BlindAlerter) *TexasHoldEm {
	return &TexasHoldEm{
		store:   store,
		alerter: alerter,
	}
}

func (game *TexasHoldEm) Start(numPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	var blindTime time.Duration = 0
	multiplier := time.Duration(5 + numPlayers)
	for _, blind := range blinds {
		game.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + multiplier*time.Minute
	}
}

func (game *TexasHoldEm) Finish(winner string) {
	game.store.RecordWin(winner)
}

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	input  *bufio.Scanner
	output io.Writer
	game   Game
}

func (cli *CLI) readLine() string {
	cli.input.Scan()
	return cli.input.Text()
}

func extractWinner(text string) string {
	return strings.Replace(text, " wins", "", 1)
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.output, PlayerPrompt)
	numPlayers, _ := strconv.Atoi(cli.readLine())
	cli.game.Start(numPlayers)
	userInput := cli.readLine()
	cli.game.Finish(extractWinner(userInput))
}

func NewCLI(input io.Reader, output io.Writer, game Game) *CLI {
	return &CLI{
		input:  bufio.NewScanner(input),
		output: output,
		game:   game,
	}
}
