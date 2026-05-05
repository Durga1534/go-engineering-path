package expense

import (
	"time"
)

func Summary(expenses []Expense, month int) float64 {
	var total float64
	currentYear := time.Now().Year()

	for _, e := range expenses {
		if month == 0 || (int(e.Date.Month()) == month && e.Date.Year() == currentYear) {
			total += e.Amount
		}
	}
	return total
}
