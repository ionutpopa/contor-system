package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

// Structuri si tipuri pentru reprezentarea sistemului

// Custom type for state open and close
type StateType string

// Define open and close on StateType
const (
	StateOpen  StateType = "open"
	StateClose StateType = "close"
)

type Source struct {
	ID          string  `json:"id"`
	Power       float64 `json:"power"` // MW
	Voltage     float64 `json:"voltage"`
	ConnectedTo string  `json:"connectedTo"`
}

type Transformer struct {
	ID            string  `json:"id"`
	InputVoltage  float64 `json:"inputVoltage"`
	OutputVoltage float64 `json:"outputVoltage"`
	ConnectedTo   string  `json:"connectedTo"`
}

type Line struct {
	ID          string  `json:"id"`
	Voltage     float64 `json:"voltage"`
	Length      int     `json:"length"` // km
	ConnectedTo string  `json:"connectedTo"`
	Area        float64 `json:"area"`
	Currnet     float64 `json:"current"`
	Ro          float64 `json:"ro"`
}

type Consumer struct {
	ID          string  `json:"id"`
	Power       float64 `json:"power"` // MW
	Voltage     float64 `json:"voltage"`
	ConnectedTo string  `json:"connectedTo,omitempty"`
}

type Separator struct {
	ID          string    `json:"id"`
	State       StateType `json:"state"`
	ConnectedTo string    `json:"connectedTo"`
}

type System struct {
	Source            Source        `json:"source"`
	Transformers      []Transformer `json:"transformers"`
	Lines             []Line        `json:"lines"`
	Consumers         []Consumer    `json:"consumers"`
	Separator         Separator     `json:"separator"`
	AdditionalSources []Source      `json:"additionalSources"`
}

// Funcția de conversie volți în kilovolți
func fromKiloVoltToVolt(voltage float64) float64 {
	return voltage * 1000
}

/*
Functie pentru calcul pierderi pe linie Ppierderi: I^2 * R
- Ppierderi: W
- I: A
- R: ohm
*/
func powerLineLoss(i float64, r float64) float64 {
	var lineLoss float64 = math.Pow(i, 2) * r
	return lineLoss
}

/*
Functie ce calculeaza rezistenta de pe linie: R = ro*l/A
- ro: rezistivitatea materialului conductorului: ohmi metru, de ex 0.0178 pentru cupru la 20°C
- l: lungimea liniei electrice: m
- A: aria sectiunii transversale a conductorului: m^2
*/
func lineResistence(ro float64, l int, A float64) float64 {
	var lineResistence = ro * (float64(l) / A)
	return lineResistence
}

/*
Functia pentru determinarea curentului electric: I = Pconsum/U*cosfi
- pCon: puterea consumata: W
- u: tensiunea nominala a liniei: V
- cosfi: factorul de putere
*/
func currentIntensity(pCons float64, u float64, cosfi float64) float64 {
	var current = pCons / (u * cosfi)
	return current
}

/*
Functia pentru calculul pierderilor de energie: Epierderi = Ppierderi * t: kWh
pLosses: pierderi de putere: kW
t: timp: ore
*/
func powerLosses(pLosses float64, t float64) float64 {
	var eLosses = pLosses * t
	return eLosses
}

/*
Functie calcul pierderi totale transformator: Ptotal = Pmiez + Pjoule
*/
func transformerLosses(coreLosses float64, jouleLosses float64) float64 {
	var totalLosses = coreLosses + jouleLosses
	return totalLosses
}

/*
Functie pentru calculul puterii reactive
*/
func reactivePowerTotal(u float64, i float64, sinfi float64) float64 {
	var q = u * i * sinfi
	return q
}

/*
Functie pentru calculul puterii aparente: Stotal = srqt(Ptotal^2 + Qtotal^2)
O alta formula ce nu e inclusa in functie mai este: S = P / cosfi
*/
func apparentPower(p float64, q float64) float64 {
	var s = math.Sqrt(math.Pow(p, 2) + math.Pow(q, 2))
	return s
}

/*
Functia pentru calculul puterii active
*/
func activePower(s float64, cosfi float64) float64 {
	var p = s * cosfi
	return p
}

/*
Functia pentru calculul puterii reactive
*/
func reactivePower(p float64, tanfi float64) float64 {
	var q = p * tanfi
	return q
}

/*
Functie pentru calculul cosfi: cosfi = P / sqrt(P^2 + Q^2)
*/
func cosfi(p float64, q float64) float64 {
	var cosfi = p / (math.Sqrt(math.Pow(p, 2) + math.Pow(q, 2)))
	return cosfi
}

/*
Functie pentru calculul sinfi
*/
func sinfi(q float64, s float64) float64 {
	var sinfi = q / s
	return sinfi
}

/*
Functie pentru calculul tanfi
*/
func tanfi(cosfi float64) float64 {
	var tanfi = math.Sqrt(1-cosfi) / cosfi
	return tanfi
}

/*
Functie ce converteste de la W la MW
*/
func wattToMegawatt(value float64) float64 {
	return value / 1000000
}

/*
Functia ce calculeaza puterea ce iese din transformatoare
*/
func transformerLossesBasedOnEfficency(powerInput float64, efficency float64) float64 {
	var transformerLossesOutput = powerInput * efficency
	return transformerLossesOutput
}

/*
Functie ce imparte sistemul pe zone
*/
func zones(system System) {
	print(system)
}

// Funcția principală pentru calcul
func calculateSystem(system System) {
	fmt.Println("Calculating power flow for the system...")

	// Verifică sursa inițială
	sourcePower := system.Source.Power
	sourceVoltage := system.Source.Voltage
	fmt.Printf("Source %s supplying %.2f MW at %.2f kV\n", system.Source.ID, sourcePower, sourceVoltage)

	// Traversează transformatoarele și liniile
	for _, transformer := range system.Transformers {
		fmt.Printf("Transformer %s steps %.2f kV to %.2f kV\n", transformer.ID, transformer.InputVoltage, transformer.OutputVoltage)
	}

	for _, line := range system.Lines {
		var ro = line.Ro
		var l = int(line.Length)
		var A = line.Area
		var current = line.Currnet
		var lineResistence = lineResistence(ro, l, A)
		var powerLoseesPerLine = wattToMegawatt(powerLineLoss(current, lineResistence))

		fmt.Printf("Line %s (%d km) has voltage %.2f kV\n", line.ID, line.Length, line.Voltage)
		fmt.Printf("Line resisstance: %.3f, Power losses per line: %.3f \n", lineResistence, powerLoseesPerLine)
	}

	// Calculează consumatorii
	for _, consumer := range system.Consumers {
		fmt.Printf("Consumer %s draws %.2f MW at %.2f kV\n", consumer.ID, consumer.Power, consumer.Voltage)
	}

	// Verifică separatorul și sursa adițională
	if system.Separator.State == StateClose {
		for _, additionalSource := range system.AdditionalSources {
			fmt.Printf("Additional source %s supplying %.2f MW at %.2f kV\n", additionalSource.ID, additionalSource.Power, additionalSource.Voltage)
		}
	}
}

// Funcția principală
func main() {
	// Deschide fișierul config.json
	file, err := os.Open("../config.json")

	if err != nil {
		log.Fatalf("Failed to open config.json: %v", err)
	}

	defer file.Close()

	// Decodează JSON-ul
	var system System

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&system); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	// linie de 10 km lungime, un conductor de aluminiu cu secțiunea de 𝐴=50𝑚𝑚^2 și un curent 𝐼=100

	// Rulează calculul pe baza configurării încărcate
	// calculateSystem(system)

	zones(system)
}
