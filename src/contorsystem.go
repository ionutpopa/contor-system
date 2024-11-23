package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
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
	Separators        []Separator   `json:"separator"`
	AdditionalSources []Source      `json:"additionalSources"`
}

type LogEntry struct {
	Timestamp   string `parquet:"name=timestamp, type=BYTE_ARRAY"`
	ComponentID string `parquet:"name=component_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Message     string `parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8"`
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
	// print(system)
}

// Funcția principală pentru calcul
func calculateSystem(system System) []LogEntry {
	var logs []LogEntry

	fmt.Println("Calculating power flow for the system...")

	// Verifică sursa inițială
	sourcePower := system.Source.Power
	sourceVoltage := system.Source.Voltage
	var sourceMessage = fmt.Sprintf("Source %s supplying %.2f MW at %.2f kV\n", system.Source.ID, sourcePower, sourceVoltage)

	sourceLog := LogEntry{
		Timestamp:   time.Now().String(),
		ComponentID: system.Source.ID,
		Message:     sourceMessage,
	}

	fmt.Println(sourceMessage)

	logs = append(logs, sourceLog)

	// Traversează transformatoarele și liniile
	for _, transformer := range system.Transformers {
		var transformerMessage = fmt.Sprintf("Transformer %s steps %.2f kV to %.2f kV\n", transformer.ID, transformer.InputVoltage, transformer.OutputVoltage)
		fmt.Println(transformerMessage)

		transformerLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: transformer.ID,
			Message:     transformerMessage,
		}

		logs = append(logs, transformerLog)
	}

	for _, line := range system.Lines {
		var ro = line.Ro
		var l = int(line.Length)
		var A = line.Area
		var current = line.Currnet
		var lineResistence = lineResistence(ro, l, A)
		var powerLoseesPerLine = wattToMegawatt(powerLineLoss(current, lineResistence))

		var lineInfoMessage = fmt.Sprintf("Line %s (%d km) has voltage %.2f kV\n", line.ID, line.Length, line.Voltage)
		var linePowerLosses = fmt.Sprintf("Power losses per line: %.3f \n", powerLoseesPerLine)

		fmt.Println(lineInfoMessage)
		fmt.Println(linePowerLosses)

		lineInfoLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: line.ID,
			Message:     lineInfoMessage,
		}

		lineLossesLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: line.ID,
			Message:     linePowerLosses,
		}

		logs = append(logs, lineInfoLog)
		logs = append(logs, lineLossesLog)
	}

	// Calculează consumatorii
	for _, consumer := range system.Consumers {
		var consumerMessage = fmt.Sprintf("Consumer %s draws %.2f MW at %.2f kV\n", consumer.ID, consumer.Power, consumer.Voltage)
		fmt.Println(consumerMessage)

		consumerLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: consumer.ID,
			Message:     consumerMessage,
		}

		logs = append(logs, consumerLog)
	}

	// Verifică separatorul și sursa adițională
	for _, separator := range system.Separators {
		var separatorMessage = fmt.Sprintf("Separator %s is in %s state", separator.ID, separator.State)

		fmt.Println(separatorMessage)

		separatorLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: separator.ID,
			Message:     separatorMessage,
		}

		logs = append(logs, separatorLog)

		if separator.State == StateClose {
			for _, additionalSource := range system.AdditionalSources {
				var additionalSourceMessage = fmt.Sprintf("Additional source %s supplying %.2f MW at %.2f kV\n", additionalSource.ID, additionalSource.Power, additionalSource.Voltage)
				fmt.Println(additionalSourceMessage)

				additionalSourceLog := LogEntry{
					Timestamp:   time.Now().String(),
					ComponentID: additionalSource.ID,
					Message:     additionalSourceMessage,
				}

				logs = append(logs, additionalSourceLog)
			}
		}
	}

	return logs
}

// ensureDirectory ensures the directory exists, creating it if necessary.
func ensureDirectory(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// Funcția principală
func main() {
	// Deschide fișierul config.json
	file, configErr := os.Open("./config.json")

	if configErr != nil {
		log.Fatalf("Failed to open config.json: %v", configErr)
	}

	defer file.Close()

	// Decodează JSON-ul
	var system System

	decoder := json.NewDecoder(file)

	if decodeError := decoder.Decode(&system); decodeError != nil {
		log.Fatalf("Failed to decode JSON: %v", decodeError)
	}

	// Open a Parquet file named with the current date\
	fileName := fmt.Sprintf("logs/%s", time.Now().Format("2006-01-02")+".parquet")

	fileNameError := ensureDirectory("logs")

	if fileNameError != nil {
		log.Fatalf("Failed to create logs directory: %v", fileNameError)
	}

	fw, parquetWritterError := NewLocalFileWriter(fileName)
	if parquetWritterError != nil {
		log.Println("Can't create local file", parquetWritterError)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(LogEntry), 4)

	// pw.RowGroupSize = 128 * 1024 * 1024 //128M
	// pw.RowGroupSize = 1 * 1024 //1k
	// pw.PageSize = 1 * 1024     //1K
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	defer pw.WriteStop()
	defer fw.Close()

	// Create a ticker for logging every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Log system data every second
	for range ticker.C {
		// Collect logs from the system
		// Rulează calculul pe baza configurării încărcate
		logs := calculateSystem(system)

		fmt.Println(logs)

		// Write logs to the Parquet file
		for _, logEntry := range logs {
			if err := pw.Write(logEntry); err != nil {
				log.Printf("Failed to write log entry: %v", err)
			}
			// Force flush to immediately write the log to disk
			if err := pw.Flush(true); err != nil {
				log.Printf("Failed to flush parquet writer: %v", err)
			}
		}

		fmt.Println("Logged data to Parquet file")
	}

	// zones(system)
}
