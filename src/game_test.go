package poker_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	poker "github.com/tulkdan/go-web/src"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedule alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5, io.Discard)

		cases := []scheduleAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.amount, c.at), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatal("expected a blind alert to be scheduled")
				}

				alert := blindAlerter.alerts[i]

				assertScheduledAlert(t, alert, c)
			})
		}
	})

	t.Run("schedule alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7, io.Discard)

		cases := []scheduleAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]

				assertScheduledAlert(t, got, want)
			})
		}
	})
}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}
