package gentable

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.ANSI)
}

var txs []tx

type tx struct {
	id       uuid.UUID
	pounds   int
	date     time.Time
	postCode string
}

func TestMain(m *testing.M) {
	f, err := os.Open("./testdata/uk-house-transactions.csv")
	if err != nil {
		panic(err.Error())
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err.Error())
	}
	txs = make([]tx, len(records))
	for i, fields := range records {
		id, err := uuid.Parse(fields[0])
		if err != nil {
			panic(err.Error())
		}

		pounds, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(err.Error())
		}

		date, err := time.Parse("2006-01-02 15:04", fields[2])
		if err != nil {
			panic(err.Error())
		}

		txs[i] = tx{
			id:       id,
			pounds:   pounds,
			date:     date,
			postCode: fields[3],
		}
	}

	os.Exit(m.Run())
}
