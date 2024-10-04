package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	poker "github.com/tulkdan/go-web/src"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type scheduleAlert struct {
	at     time.Duration
	amount int
}

func (s scheduleAlert) String() string {
	return fmt.Sprintf("%d chip at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduleAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduleAlert{at, amount})
}

type GameSpy struct {
	StartCalled  bool
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int, to io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nChris wins\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Chris" {
			t.Errorf("got %q, want %q", game.FinishedWith, "Chris")
		}
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nCleo wins\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		if game.FinishedWith != "Cleo" {
			t.Errorf("got %q, want %q", game.FinishedWith, "Chris")
		}
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}

func assertScheduledAlert(t testing.TB, got, want scheduleAlert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}

	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()

	want := strings.Join(messages, "")
	got := stdout.String()

	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
