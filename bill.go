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

func statement(invoice Invoice, plays Plays) string {
	totalAmount := 0.0
	volumeCredits := 0.0
	result := fmt.Sprintf("Statement for %s\n", invoice.Customer)

	for _, performance := range invoice.Performances {
		play := plays[performance.PlayID]
		thisAmount := amountFor(performance, play)
		// add volume credits
		volumeCredits += math.Max(float64(performance.Audience-30), 0)
		// add extra credit for every ten comedy attendees
		if "comedy" == playType(play) {
			volumeCredits += math.Floor(float64(performance.Audience / 5))
		}

		// print line for this order
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", playName(play), thisAmount/100, performance.Audience)
		totalAmount += thisAmount
	}
	result += fmt.Sprintf("Amount owed is $%.2f\n", totalAmount/100)
	result += fmt.Sprintf("you earned %.0f credits\n", volumeCredits)
	return result
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
