package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	store   PlayerStore
	input   *bufio.Scanner
	output  io.Writer
	alerter BlindAlerter
}

func (cli *CLI) readLine() string {
	cli.input.Scan()
	return cli.input.Text()
}

func (cli *CLI) scheduleBlindAlerts(numPlayers int) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	var blindTime time.Duration = 0
	multiplier := time.Duration(5 + numPlayers)
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + multiplier*time.Minute
	}
}

func extractWinner(text string) string {
	return strings.Replace(text, " wins", "", 1)
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.output, PlayerPrompt)
	numPlayers, _ := strconv.Atoi(cli.readLine())
	cli.scheduleBlindAlerts(numPlayers)
	userInput := cli.readLine()
	cli.store.RecordWin(extractWinner(userInput))
}

func NewCLI(store PlayerStore, input io.Reader, output io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		input:   bufio.NewScanner(input),
		output:  output,
		alerter: alerter,
	}
}
