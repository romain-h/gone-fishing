package expenses

import (
	"fmt"
	"time"
)

var JointAccountId = "acc_00009Yx2VDcMTatdAVTu0P"
var SudamericaGid = "7786790"
var QuitoCatsGid = "11593186"
var ColombianosGid = "11862121"

type SplitwiseExpense struct {
	Id          int       `json:"id"`
	Gid         int       `json:"group_id"`
	CreatedAt   time.Time `json:"date"`
	Cost        string    `json:"cost"`
	Description string    `json:"description"`
	Details     string    `json:"details"`
	DeletedAt   time.Time `json:"deleted_at"`
	Users       []struct {
		UserId int    `json:"user_id"`
		Share  string `json:"owed_share"`
	} `json:"users"`
}
type SplitwiseExpenses struct {
	Expenses []SplitwiseExpense `json:"expenses"`
}

type Transaction struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created"`
	Amount    int       `json:"amount"`
	Merchant  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"merchant"`
	Description       string `json:"Description"`
	Notes             string `json:"notes"`
	LocalAmount       int    `json:"local_amount"`
	LocalCurrency     string `json:"local_currency"`
	IncludeInSpending bool   `json:"include_in_spending"`
	Metadata          struct {
		MastercardApprovalType string `json:"mastercard_approval_type,omitempty"`
	} `json:"metadata"`
	Settled string `json:"settled,omitempty"`
}
type Transactions struct {
	Transactions []Transaction
}

type Expense struct {
	Id            string    `json:"id"`
	CreatedAt     time.Time `json:"created"`
	Amount        float64   `json:"amount"`
	LocalAmount   float64   `json:"local_amount"`
	LocalCurrency string    `json:"local_currency"`
	Description   string    `json:"description"`
	Notes         string    `json:"notes"`
}

func (e Expense) TestString() string {
	return fmt.Sprintf("Expense: { %s %s %.2f %.2f %s %s %s }", e.Id, e.CreatedAt, e.Amount, e.LocalAmount, e.LocalCurrency, e.Description, e.Notes)
}

type GroupedExpenses struct {
	Date       time.Time
	Expenses   []Expense
	Total      float64
	TotalLocal float64
}

type ByWeek struct {
	Week      int
	StartDate time.Time
	EndDate   time.Time
	Days      []GroupedExpenses
	Mean      float64
}
