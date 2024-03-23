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

type Bill struct {
	Customer           string
	TotalAmount        float64
	TotalVolumeCredits float64
}

type Rate struct {
	Play          Play
	Amount        float64
	VolumnCredits float64
	Audience      int
}

func renderPlainText(invoice Invoice, plays Plays) string {
	bill := Bill{
		Customer:           invoice.Customer,
		TotalAmount:        totalAmountFor(invoice, plays),
		TotalVolumeCredits: totalVolumeCreditsFor(invoice.Performances, plays),
	}
	rates := []Rate{}

	for _, performance := range invoice.Performances {
		rate := Rate{
			Play:          playFor(plays, performance),
			Amount:        amountFor(performance, playFor(plays, performance)),
			VolumnCredits: volumeCreditsFor(plays, performance),
			Audience:      performance.Audience,
		}
		rates = append(rates, rate)

	}
	result := fmt.Sprintf("Statement for %s\n", bill.Customer)
	for _, rate := range rates {
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", rate.Play.Name, rate.Amount/100, rate.Audience)
	}
	result += fmt.Sprintf("Amount owed is $%.2f\n", bill.TotalAmount/100)
	result += fmt.Sprintf("you earned %.0f credits\n", bill.TotalVolumeCredits)
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
