package expenses

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/romain-h/gone-fishing/internal/cache"
	"github.com/romain-h/gone-fishing/internal/config"
	"github.com/romain-h/gone-fishing/internal/oauth"
)

func GetLocation(date time.Time) *time.Location {
	// 29/10/18 - Peru
	// 23/11/18 - Bolivia
	// 27/11/18 - Chile
	// 20/12/18 11:00 - Argentina
	// 17/02/19 18:00 - Santiago
	// other - Quito
	var loc *time.Location
	if date.Before(time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)) {
		loc, _ = time.LoadLocation("America/Lima")
	} else if date.Before(time.Date(2018, time.November, 23, 0, 0, 0, 0, time.UTC)) {
		loc, _ = time.LoadLocation("America/La_Paz")
	} else if date.Before(time.Date(2018, time.November, 27, 0, 0, 0, 0, time.UTC)) {
		loc, _ = time.LoadLocation("America/Santiago")
	} else if date.Before(time.Date(2018, time.December, 20, 11, 0, 0, 0, time.UTC)) {
		loc, _ = time.LoadLocation("America/Argentina/Jujuy")
	} else if date.Before(time.Date(2019, time.February, 17, 18, 0, 0, 0, time.UTC)) {
		loc, _ = time.LoadLocation("America/Santiago")
	} else {
		loc, _ = time.LoadLocation("America/Guayaquil")
	}
	return loc
}

func GetExps(cfg config.Config, cache cache.CacheManager, monzo oauth.AuthProvider, splitwise oauth.AuthProvider) []Expense {
	var allExps []Expense
	spEx := GetSplitwiseExpenses(cache, splitwise)
	mEx := GetMonzoTransactions(cfg, cache, monzo)

	isIgnored := regexp.MustCompile("#ignore")
	isNightsRegx := regexp.MustCompile("#nights|#spread")
	isCashRegx, _ := regexp.Compile("#cash")
	cashWithLocalInfo := regexp.MustCompile(`#cash (\w{3}) ([0-9]*\.?[0-9]+)`)

	var inputExps []Expense
	for _, v := range spEx.Expenses {
		// skip deleted expenses
		if !v.DeletedAt.IsZero() {
			continue
		}

		// Handle Quito Cats expenses
		// TODO refactor in a better way
		var famount float64
		if v.Gid == 11593186 {
			var count int
			var totalShared float64

			for _, u := range v.Users {
				// If it's me or Alison then we take the share into account
				if u.UserId == 841958 || u.UserId == 2788798 {
					count = count + 1
					s, _ := strconv.ParseFloat(u.Share, 64)
					totalShared = totalShared + s
				}
			}

			// we keep it only if we are both included into this expense
			if count == 2 {
				famount = totalShared * 0.76 // Convert USD to GBP
			} else {
				continue
			}
		}

		// Colombianos
		if v.Gid == 11862121 {
			var count int
			var totalShared float64

			for _, u := range v.Users {
				// If it's me or Alison then we take the share into account
				if u.UserId == 841958 || u.UserId == 2788798 {
					count = count + 1
					s, _ := strconv.ParseFloat(u.Share, 64)
					totalShared = totalShared + s
				}
			}

			// we keep it only if we are both included into this expense
			if count == 2 {
				famount = totalShared * 0.00024 // COP to GBP
			} else {
				continue
			}
		}

		// if not a Quito expense handled
		if famount == 0 {
			famount, _ = strconv.ParseFloat(v.Cost, 64)
		}

		// createdAt := v.CreatedAt.Truncate(24 * time.Hour)
		loc := GetLocation(v.CreatedAt)
		createdAt := v.CreatedAt.In(loc)

		var ex = Expense{
			Id:          strconv.Itoa(v.Id),
			CreatedAt:   createdAt,
			Amount:      famount,
			Description: v.Description,
			Notes:       v.Details,
		}
		inputExps = append(inputExps, ex)
	}

	for _, v := range mEx.Transactions {
		// Skip declined transactions
		if !v.IncludeInSpending {
			continue
		}
		famount := float64(v.Amount) / 100
		if famount > 0 {
			continue
		}
		famount = math.Abs(famount)
		flocamount := math.Abs(float64(v.LocalAmount))
		// zero-decimal currencies
		// Chilean pesos is zero-decimal so no need to divide by 100
		if v.LocalCurrency != "CLP" {
			flocamount = flocamount / 100
		}

		loc := GetLocation(v.CreatedAt)
		createdAt := v.CreatedAt.In(loc)

		var ex = Expense{
			Id:            v.Id,
			CreatedAt:     createdAt,
			Amount:        famount,
			LocalAmount:   flocamount,
			LocalCurrency: v.LocalCurrency,
			Description:   v.Merchant.Name,
			Notes:         v.Notes,
		}
		inputExps = append(inputExps, ex)
	}

	for _, exp := range inputExps {

		if isIgnored.MatchString(exp.Notes) {
			continue
		}

		if isNightsRegx.MatchString(exp.Notes) {
			nightsExpenses := HandleNights(exp)
			allExps = append(allExps, nightsExpenses...)
			continue
		}

		if isCashRegx.MatchString(exp.Notes) {
			// Splitwise cash entries follow #cash PEN 1499
			// We need to manually add LocalAmount and LocalCurrency
			if cashWithLocalInfo.MatchString(exp.Notes) {
				m := cashWithLocalInfo.FindStringSubmatch(exp.Notes)
				exp.LocalCurrency = m[1]
				exp.LocalAmount, _ = strconv.ParseFloat(m[2], 10)
			}

			cashExps := HandleCash(exp)
			allExps = append(allExps, cashExps...)

			continue
		}

		allExps = append(allExps, exp)
	}

	return allExps
}

func GetAllExpenses(cfg config.Config, cache cache.CacheManager, monzo oauth.AuthProvider, splitwise oauth.AuthProvider) []GroupedExpenses {
	allExps := GetExps(cfg, cache, monzo, splitwise)
	conv := make(map[time.Time]GroupedExpenses)
	for _, exp := range allExps {
		year, month, day := exp.CreatedAt.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		ge := conv[date]
		ge.Date = date
		ge.Expenses = append(ge.Expenses, exp)
		ge.Total = ge.Total + exp.Amount
		conv[date] = ge
	}

	var ret []GroupedExpenses
	for _, v := range conv {
		ret = append(ret, v)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Date.Before(ret[j].Date) })

	return ret
}

func GetStats(exps []GroupedExpenses, grandTotal bool) (float64, float64) {
	var data stats.Float64Data
	for _, exp := range exps {
		if grandTotal {
			if exp.Date.Before(time.Now()) {
				data = append(data, exp.Total)
			}
		} else {
			data = append(data, exp.Total)
		}
	}

	mean, _ := data.Mean()
	median, _ := data.Median()
	return mean, median
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}

func GetExpensesByWeek(exps []GroupedExpenses) []ByWeek {
	grouped := make(map[int]ByWeek)
	var all []ByWeek

	for _, ex := range exps {
		year, i := ex.Date.ISOWeek()
		n := (year - 2018) * 52
		wi := i + n - 38
		wex := grouped[wi]
		wex.Week = wi
		wex.Days = append(wex.Days, ex)
		grouped[wi] = wex
	}

	for _, week := range grouped {
		week.Mean, _ = GetStats(week.Days, false)
		sample := week.Days[0].Date
		_, isoWeek := sample.ISOWeek()
		week.StartDate = FirstDayOfISOWeek(sample.Year(), isoWeek, sample.Location())
		week.EndDate = week.StartDate.AddDate(0, 0, 6)
		all = append(all, week)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Week > all[j].Week })

	return all
}
