package expenses

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	"github.com/romain-h/gone-fishing/internal/cache"
	"github.com/romain-h/gone-fishing/internal/config"
	"github.com/romain-h/gone-fishing/internal/oauth"
)

func GetMonzoTransactions(cfg config.Config, cache cache.CacheManager, monzo oauth.AuthProvider) Transactions {
	var trans Transactions
	res, _ := cache.GetByte("expenses_monzo")

	if len(res) == 0 {
		FetchMonzoTransactions(cfg, cache, monzo)
	}
	res, _ = cache.GetByte("expenses_monzo")
	decBuf := bytes.NewBuffer(res)
	if err := gob.NewDecoder(decBuf).Decode(&trans); err != nil {
		log.Fatal(err)
	}

	return trans
}

func FetchMonzoTransactions(cfg config.Config, cache cache.CacheManager, monzo oauth.AuthProvider) {
	var trans Transactions
	resp, err := monzo.GetClient().Get(fmt.Sprintf("https://api.monzo.com/transactions?expand[]=merchant&account_id=%s&since=%s", cfg.Monzo.AccountID, cfg.StartDate))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&trans); err != nil {
		log.Fatal(err)
	}

	encBuf := new(bytes.Buffer)
	err = gob.NewEncoder(encBuf).Encode(trans)
	if err != nil {
		log.Fatal(err)
	}
	value := encBuf.Bytes()
	cache.SetByte("expenses_monzo", value)
}

func GetSplitwiseExpenses(cache cache.CacheManager, splitwise oauth.AuthProvider) SplitwiseExpenses {
	var exps SplitwiseExpenses
	res, _ := cache.GetByte("expenses_splitwise")

	if len(res) == 0 {
		FetchSplitwiseExpenses(cache, splitwise)
	}
	res, _ = cache.GetByte("expenses_splitwise")
	decBuf := bytes.NewBuffer(res)
	if err := gob.NewDecoder(decBuf).Decode(&exps); err != nil {
		log.Fatal(err)
	}

	return exps
}

func FetchSplitwiseExpenses(cache cache.CacheManager, splitwise oauth.AuthProvider) {
	var exps SplitwiseExpenses
	var expsQuito SplitwiseExpenses
	var expsColombianos SplitwiseExpenses
	resp, err := splitwise.GetClient().Get("https://www.splitwise.com/api/v3.0/get_expenses?group_id=7786790&limit=0")
	respQuito, _ := splitwise.GetClient().Get("https://www.splitwise.com/api/v3.0/get_expenses?group_id=11593186&limit=0")
	respColombianos, _ := splitwise.GetClient().Get("https://www.splitwise.com/api/v3.0/get_expenses?group_id=11862121&limit=0")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	defer respQuito.Body.Close()
	defer respColombianos.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&exps); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(respQuito.Body).Decode(&expsQuito); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(respColombianos.Body).Decode(&expsColombianos); err != nil {
		log.Fatal(err)
	}
	exps.Expenses = append(exps.Expenses, expsQuito.Expenses...)
	exps.Expenses = append(exps.Expenses, expsColombianos.Expenses...)

	encBuf := new(bytes.Buffer)
	err = gob.NewEncoder(encBuf).Encode(exps)
	if err != nil {
		log.Fatal(err)
	}
	value := encBuf.Bytes()
	cache.SetByte("expenses_splitwise", value)
}
