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

const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"

const BadWinnerInputErrMsg = "Bad value received for winner input, correct format is '<player name> wins'"

type CLI struct {
	input  *bufio.Scanner
	output io.Writer
	game   Game
}

func (cli *CLI) readLine() string {
	cli.input.Scan()
	return cli.input.Text()
}

func extractWinner(text string) (string, error) {
	winner := strings.Replace(text, " wins", "", 1)
	if winner == text {
		return "", fmt.Errorf("invalid winner input %q", text)
	}
	return winner, nil
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.output, PlayerPrompt)
	numPlayers, err := strconv.Atoi(cli.readLine())
	if err != nil {
		fmt.Fprint(cli.output, BadPlayerInputErrMsg)
		return
	}
	cli.game.Start(numPlayers)
	userInput := cli.readLine()
	if userInput != "" {
		winner, err := extractWinner(userInput)
		if err != nil {
			fmt.Fprint(cli.output, BadWinnerInputErrMsg)
			return
		}
		cli.game.Finish(winner)
	}
}

func NewCLI(input io.Reader, output io.Writer, game Game) *CLI {
	return &CLI{
		input:  bufio.NewScanner(input),
		output: output,
		game:   game,
	}
}
