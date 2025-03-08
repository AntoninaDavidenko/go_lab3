package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type CalculationResult struct {
	FirstCheckProfit  float64
	SecondCheckProfit float64
}

func normalDistribution(p, pC, sigma float64) float64 {
	return (1 / (sigma * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow((p-pC)/sigma, 2)/2)
}

func integrateNormalDistribution(a, b, mean, stdDev float64, steps int) float64 {
	stepSize := (b - a) / float64(steps)
	area := 0.0
	for i := 0; i < steps; i++ {
		x1 := a + float64(i)*stepSize
		x2 := a + float64(i+1)*stepSize
		y1 := normalDistribution(x1, mean, stdDev)
		y2 := normalDistribution(x2, mean, stdDev)
		area += (y1 + y2) / 2 * stepSize
	}
	return area
}

func calculate(pC, sigma1, sigma2, costB float64) CalculationResult {
	lowerBound := pC - 0.25
	upperBound := pC + 0.25
	steps := 10000

	probability := math.Round((integrateNormalDistribution(lowerBound, upperBound, pC, sigma1, steps))*100) / 100
	power1 := pC * 24 * probability
	profit := power1 * costB

	power2 := pC * 24 * (1 - probability)
	penalty := power2 * costB
	firstCheckProfit := profit - penalty

	probability2 := math.Round((integrateNormalDistribution(lowerBound, upperBound, pC, sigma2, steps))*100) / 100
	power3 := pC * 24 * probability2
	profit2 := power3 * costB

	power4 := pC * 24 * (1 - probability2)
	penalty2 := power4 * costB
	secondCheckProfit := profit2 - penalty2

	return CalculationResult{firstCheckProfit, secondCheckProfit}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	if r.Method == http.MethodPost {
		r.ParseForm()
		pC, _ := strconv.ParseFloat(r.FormValue("pC"), 64)
		sigma1, _ := strconv.ParseFloat(r.FormValue("sigma1"), 64)
		sigma2, _ := strconv.ParseFloat(r.FormValue("sigma2"), 64)
		costB, _ := strconv.ParseFloat(r.FormValue("costB"), 64)
		result := calculate(pC, sigma1, sigma2, costB)
		tmpl.Execute(w, result)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
