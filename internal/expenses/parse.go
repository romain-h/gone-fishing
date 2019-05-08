package expenses

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var dateSubReg = `([0-9]{2}/[0-9]{2}/[0-9]{2})`
var dateLayout = "02/01/06"

func HandleNights(exp Expense) []Expense {
	var nightExpenses []Expense
	nightsRegx := regexp.MustCompile(`^#(nights|spread) ` + dateSubReg + `( to ` + dateSubReg + `)?`)

	note := []byte(exp.Notes)

	allIndexes := nightsRegx.FindAllSubmatchIndex(note, -1)
	for _, loc := range allIndexes {
		from, _ := time.Parse(dateLayout, string(note[loc[4]:loc[5]]))
		unit := 1

		if loc[8] != -1 {
			to, _ := time.Parse(dateLayout, string(note[loc[8]:loc[9]]))
			unit = int(to.Sub(from).Hours()/24) + 1
		}

		for i := 0; i < unit; i++ {
			hoursToAdd, _ := time.ParseDuration(fmt.Sprint(i*24, "h"))
			newExp := exp
			newExp.CreatedAt = from.Add(hoursToAdd)
			newExp.Amount = exp.Amount / float64(unit)
			newExp.LocalAmount = exp.LocalAmount / float64(unit)
			nightExpenses = append(nightExpenses, newExp)
		}
	}

	return nightExpenses
}

func HandleCash(exp Expense) []Expense {
	var cashExpenses []Expense
	note := []byte(exp.Notes)

	exchangeRate := exp.LocalAmount / exp.Amount

	cashExpenseRegx, _ := regexp.Compile(`(?m)^` + dateSubReg + `( to ` + dateSubReg + `)? ([0-9]*\.?[0-9]+) (.*)$`)

	for _, loc := range cashExpenseRegx.FindAllSubmatchIndex(note, -1) {
		unit := 1
		created, _ := time.Parse(dateLayout, string(note[loc[2]:loc[3]]))
		localAmount, _ := strconv.ParseFloat(string(note[loc[8]:loc[9]]), 10)

		if loc[6] != -1 {
			to, _ := time.Parse(dateLayout, string(note[loc[6]:loc[7]]))
			unit = int(to.Sub(created).Hours()/24) + 1
		}

		for i := 0; i < unit; i++ {
			hoursToAdd, _ := time.ParseDuration(fmt.Sprint(i*24, "h"))
			newExp := Expense{
				Id:            fmt.Sprint(exp.Id, "-cash-", i),
				CreatedAt:     created.Add(hoursToAdd),
				Amount:        (localAmount / exchangeRate) / float64(unit),
				LocalAmount:   localAmount / float64(unit),
				LocalCurrency: exp.LocalCurrency,
				Description:   string(note[loc[10]:loc[11]]),
			}
			cashExpenses = append(cashExpenses, newExp)
		}
	}

	return cashExpenses
}
