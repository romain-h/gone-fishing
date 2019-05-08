package expenses_test

import (
	"strings"
	"testing"
	"time"

	. "github.com/romain-h/gone-fishing/internal/expenses"
)

func expsString(exp []Expense) string {
	cases := []string{}
	for _, e := range exp {
		cases = append(cases, e.TestString())
	}
	return strings.Join(cases, "\n")
}

func TestHandleNights(t *testing.T) {
	var empty []Expense

	cases := []struct {
		exp  Expense
		exps []Expense
	}{
		{
			Expense{
				Id:            "2",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        97.19,
				LocalAmount:   418.00,
				LocalCurrency: "PEN",
				Notes:         `Cia`,
			},
			empty,
		},
		{
			Expense{
				Id:            "2",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        97.19,
				LocalAmount:   418.00,
				LocalCurrency: "PEN",
				Notes:         `#nights 18/09/18`,
			},
			[]Expense{
				Expense{
					Id:            "2",
					CreatedAt:     time.Date(2018, 9, 18, 0, 0, 0, 0, time.UTC),
					Amount:        97.19,
					LocalAmount:   418.00,
					LocalCurrency: "PEN",
					Notes:         `#nights 18/09/18`,
				},
			},
		},
		{
			Expense{
				Id:            "3",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        97.19,
				LocalAmount:   418.00,
				LocalCurrency: "PEN",
				Notes:         `#nights 18/09/18 to 20/09/18`,
			},
			[]Expense{
				Expense{
					Id:            "3",
					CreatedAt:     time.Date(2018, 9, 18, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#nights 18/09/18 to 20/09/18`,
				},
				Expense{
					Id:            "3",
					CreatedAt:     time.Date(2018, 9, 19, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#nights 18/09/18 to 20/09/18`,
				},
				Expense{
					Id:            "3",
					CreatedAt:     time.Date(2018, 9, 20, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#nights 18/09/18 to 20/09/18`,
				},
			},
		},
		{
			Expense{
				Id:            "4",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        97.19,
				LocalAmount:   418.00,
				LocalCurrency: "PEN",
				Notes:         `#spread 18/09/18 to 20/09/18`,
			},
			[]Expense{
				Expense{
					Id:            "4",
					CreatedAt:     time.Date(2018, 9, 18, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#spread 18/09/18 to 20/09/18`,
				},
				Expense{
					Id:            "4",
					CreatedAt:     time.Date(2018, 9, 19, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#spread 18/09/18 to 20/09/18`,
				},
				Expense{
					Id:            "4",
					CreatedAt:     time.Date(2018, 9, 20, 0, 0, 0, 0, time.UTC),
					Amount:        97.19 / 3,
					LocalAmount:   418.00 / 3,
					LocalCurrency: "PEN",
					Notes:         `#spread 18/09/18 to 20/09/18`,
				},
			},
		},
	}

	for _, c := range cases {
		res := HandleNights(c.exp)
		if expsString(res) != expsString(c.exps) {
			t.Fatalf("no match:\n\n%s\n\n%s", expsString(res), expsString(c.exps))
		}
	}
}

func TestHandleCash(t *testing.T) {
	var empty []Expense

	cases := []struct {
		exp  Expense
		exps []Expense
	}{
		{
			Expense{
				Id:            "1",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        128.29,
				LocalAmount:   1938.20,
				LocalCurrency: "PEN",
				Description:   "ATM",
			},
			empty,
		},
		{
			Expense{
				Id:            "1",
				CreatedAt:     time.Date(2018, 9, 13, 12, 30, 0, 0, time.UTC),
				Amount:        128.29,
				LocalAmount:   1938.20,
				LocalCurrency: "PEN",
				Description:   "ATM",
				Notes: `#cash
28/09/18 10 Lima metro bus
29/09/18 65 airport taxi
24/10/18 to 25/10/18 140 Colca canyon entries
			`,
			},
			[]Expense{
				Expense{
					Id:            "1-cash-0",
					CreatedAt:     time.Date(2018, 9, 28, 0, 0, 0, 0, time.UTC),
					Amount:        0.66,
					LocalAmount:   10,
					LocalCurrency: "PEN",
					Description:   "Lima metro bus",
				},
				Expense{
					Id:            "1-cash-0",
					CreatedAt:     time.Date(2018, 9, 29, 0, 0, 0, 0, time.UTC),
					Amount:        4.30,
					LocalAmount:   65,
					LocalCurrency: "PEN",
					Description:   "airport taxi",
				},
				Expense{
					Id:            "1-cash-0",
					CreatedAt:     time.Date(2018, 10, 24, 0, 0, 0, 0, time.UTC),
					Amount:        4.63,
					LocalAmount:   70,
					LocalCurrency: "PEN",
					Description:   "Colca canyon entries",
				},
				Expense{
					Id:            "1-cash-1",
					CreatedAt:     time.Date(2018, 10, 25, 0, 0, 0, 0, time.UTC),
					Amount:        4.63,
					LocalAmount:   70,
					LocalCurrency: "PEN",
					Description:   "Colca canyon entries",
				},
			},
		},
	}

	for _, c := range cases {
		res := HandleCash(c.exp)
		if expsString(res) != expsString(c.exps) {
			t.Fatalf("no match:\n\n%s\n\n%s", expsString(res), expsString(c.exps))
		}
	}
}
