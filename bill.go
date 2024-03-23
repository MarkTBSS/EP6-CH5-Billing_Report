package main

import (
	"fmt"
	"math"
)

type Plays map[string]Play

type Play struct {
	Name string
	Type string
}

type Performance struct {
	PlayID   string `json:"playID"`
	Audience int    `json:"audience"`
}

type Invoice struct {
	Customer     string        `json:"customer"`
	Performances []Performance `json:"performances"`
}

func playName(play Play) string {
	return play.Name
}

func playType(play Play) string {
	return play.Type
}

func playFor(plays Plays, performance Performance) Play {
	return plays[performance.PlayID]
}

func amountFor(performance Performance, play Play) float64 {
	result := 0.0
	switch playType(play) {
	case "tragedy":
		result = 40000
		if performance.Audience > 30 {
			result += 1000 * (float64(performance.Audience - 30))
		}
	case "comedy":
		result = 30000
		if performance.Audience > 20 {
			result += 10000 + 500*(float64(performance.Audience-20))
		}
		result += 300 * float64(performance.Audience)
	default:
		panic(fmt.Sprintf("unknow type: %s", play.Type))
	}
	return result
}

func volumeCreditsFor(plays Plays, performance Performance) float64 {
	result := 0.0
	result += math.Max(float64(performance.Audience-30), 0)
	// add extra credit for every ten comedy attendees
	if "comedy" == playType(playFor(plays, performance)) {
		result += math.Floor(float64(performance.Audience / 5))
	}
	return result
}

func totalAmountFor(invoice Invoice, plays Plays) float64 {
	result := 0.0
	for _, performance := range invoice.Performances {
		play := playFor(plays, performance)
		result += amountFor(performance, play)
	}
	return result
}

func totalVolumeCreditsFor(performances []Performance, plays Plays) float64 {
	result := 0.0
	for _, performance := range performances {
		result += volumeCreditsFor(plays, performance)
	}
	return result
}

func renderPlainText(invoice Invoice, plays Plays) string {
	result := fmt.Sprintf("Statement for %s\n", invoice.Customer)
	for _, performance := range invoice.Performances {
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", playName(playFor(plays, performance)), amountFor(performance, playFor(plays, performance))/100, performance.Audience)
	}
	result += fmt.Sprintf("Amount owed is $%.2f\n", totalAmountFor(invoice, plays)/100)
	result += fmt.Sprintf("you earned %.0f credits\n", totalVolumeCreditsFor(invoice.Performances, plays))
	return result
}

func statement(invoice Invoice, plays Plays) string {
	return renderPlainText(invoice, plays)
}

func main() {
	inv := Invoice{
		Customer: "Bigco",
		Performances: []Performance{
			{PlayID: "hamlet", Audience: 55},
			{PlayID: "as-like", Audience: 35},
			{PlayID: "othello", Audience: 40},
		}}
	plays := map[string]Play{
		"hamlet":  {Name: "Hamlet", Type: "tragedy"},
		"as-like": {Name: "As You Like It", Type: "comedy"},
		"othello": {Name: "Othello", Type: "tragedy"},
	}

	bill := statement(inv, plays)
	fmt.Println(bill)
}
