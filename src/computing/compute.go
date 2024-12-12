package computing

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"contor-system/src/utils"
)

var logs []LogEntry

type LogEntry = utils.LogEntry

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

func reactivePowerLineLoss(i float64, x float64) float64 {
	var lineLoss float64 = math.Pow(i, 2) * x
	return lineLoss
}

/*
Functie ce calculeaza rezistenta de pe linie: R = ro*l/A ohm/km
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
func sinfi(cosfi float64) float64 {
	var sinfi = math.Sqrt(1 - math.Pow(cosfi, 2))
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
Functie pentru distanta medie geometrica
*/
func geometricDistance(Drs float64, Dst float64, Drt float64) float64 {
	var Dm = math.Pow(Drs*Dst*Drt, 1/3)
	return Dm
}

/*
Functie pentru raza echivalenta
*/
func equivalentRadius(r float64) float64 {
	var re = math.Pow(math.E, -(1/4)) * r
	return re
}

/*
Inductanta pe unitate de lungime
*/
func inductionOnLength(Dm float64, re float64) float64 {
	var L = 2 * math.Pow(10, -7) * math.Log(Dm/re)
	return L
}

/*
Reactanta totala
*/
func totalReactation(w float64, l float64) float64 {
	var XL = w * l
	return XL
}

/*
Capacitatea liniei
Dm - Distanta medie geometrica intre conductori (m)
r - raza conductorului (m)
*/
func lineCapacity(Dm float64, r float64) float64 {
	epsilon := 8.854e-12 // Permittivity of free space (permitivitatea electrică a vidului)
	er := 1.0            // Permitivitatea relativă a aerului
	numerator := 2 * math.Pi * epsilon * er
	denumerator := math.Log(Dm / r)
	var C = numerator / denumerator
	return C
}

func findConnectedTo(data utils.System, target string) []utils.ConnectedElement {
	var results []utils.ConnectedElement

	// Check if the source matches
	if data.Source.ConnectedTo == target {
		results = append(results, utils.ConnectedElement{
			ID:      data.Source.ID,
			Details: data.Source,
		})
	}

	// Check transformers
	for _, transformer := range data.Transformers {
		if transformer.ConnectedTo == target {
			results = append(results, utils.ConnectedElement{
				ID:      transformer.ID,
				Details: transformer,
			})
		}
	}

	// Check lines
	for _, line := range data.Lines {
		if line.ConnectedTo == target {
			results = append(results, utils.ConnectedElement{
				ID:      line.ID,
				Details: line,
			})
		}
	}

	// Check consumers
	for _, consumer := range data.Consumers {
		if consumer.ConnectedTo == target {
			results = append(results, utils.ConnectedElement{
				ID:      consumer.ID,
				Details: consumer,
			})
		}
	}

	// Check separators
	for _, separator := range data.Separators {
		if separator.ConnectedTo == target {
			results = append(results, utils.ConnectedElement{
				ID:      separator.ID,
				Details: separator,
			})
		}
	}

	return results
}

// PowerMap creates a map where keys are the field names of the Data struct and values are maps with indices as keys and objects containing "power" = 0.
func PowerMap(data utils.System) map[string]map[string]map[string]float64 {
	result := make(map[string]map[string]map[string]float64)

	// Use reflection to iterate over the fields of the Data struct
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		fieldName := field.Name
		fieldValue := dataValue.Field(i)

		// If the field is a slice, we want to map each index to a map
		if fieldValue.Kind() == reflect.Slice {
			fieldLen := fieldValue.Len()
			objects := make(map[string]map[string]float64)
			for j := 0; j < fieldLen; j++ {
				objects[fmt.Sprintf("%d", j)] = map[string]float64{"power": 0}
			}
			result[fieldName] = objects
		} else {
			// If it's not a slice, create a single object with "power" = 0
			result[fieldName] = map[string]map[string]float64{
				"0": {"power": 0},
			}
		}
	}

	return result
}

var systemAfterConfig utils.System

func isActive(nID string, config *utils.System) bool {
	var isActive = true
	for _, separator := range config.Separators {
		if separator.ConnectedTo == nID && separator.State == utils.StateOpen {
			// Separatorul este deschis si este conectat in elementul 'n', calculul de puteri va fi oprit aici
			isActive = false
		}
	}
	return isActive
}

func calculatePowerFlow(id string, inputPower float64, config *utils.System, nodes map[string]interface{}, visited map[string]bool, consumersWithoutPower *[]utils.ConsumerPowerDetails) {
	const LowVoltage = 20

	// Check if the node has already been visited
	if visited[id] {
		return
	}
	visited[id] = true

	node, exists := nodes[id]
	if !exists {
		fmt.Printf("Node %s not found\n", id)
		return
	}

	fmt.Printf("Traversing node: %s\n", id)

	switch n := node.(type) {
	case utils.Source:
		fmt.Printf("Source %s receiving power: %.2f\n", n.ID, inputPower)
		// Update the source's power in the config
		if n.ID == config.Source.ID {
			config.Source.Power = inputPower
		} else {
			for i, source := range config.AdditionalSources {
				if source.ID == n.ID {
					config.AdditionalSources[i].AdditionalPower = inputPower
					break
				}
			}
		}
		calculatePowerFlow(n.ConnectedTo, inputPower, config, nodes, visited, consumersWithoutPower)

	case utils.Transformer:
		var isActive = isActive(n.ID, config)

		if !isActive {
			calculatePowerFlow(n.ConnectedTo, 0, config, nodes, visited, consumersWithoutPower)
			return
		}

		var totalCooperAndSteelLosses = (n.SteelLosses / 1000) + (n.CooperLosses / 1000)

		var outputPower float64

		// Print the losses and output power
		var transformerMessage = fmt.Sprintf("Transformer %s transferring power: %.2f -> %.2f (losses: %.2f)\n", n.ID, math.Abs(inputPower), math.Abs(outputPower), totalCooperAndSteelLosses)

		transformerInfoLog := LogEntry{
			Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
			ComponentID: n.ID,
			Message:     transformerMessage,
		}

		logs = append(logs, transformerInfoLog)

		if n.InputVoltage == LowVoltage {
			outputPower = inputPower - totalCooperAndSteelLosses
		} else {
			outputPower = inputPower
		}

		// fmt.Printf("Transformer %s transferring power: %.2f -> %.2f\n", n.ID, inputPower, outputPower)
		// Update the transformer's power in the config
		for i, transformer := range config.Transformers {
			if transformer.ID == n.ID {
				config.Transformers[i].PowerTransfered = outputPower
				break
			}
		}
		calculatePowerFlow(n.ConnectedTo, outputPower, config, nodes, visited, consumersWithoutPower)

	case utils.Line:
		var isActive = isActive(n.ID, config)

		if !isActive {
			calculatePowerFlow(n.ConnectedTo, 0, config, nodes, visited, consumersWithoutPower)
			return
		}

		// Print the power being transferred through the line
		fmt.Printf("Line %s transferring power: %.2f\n", n.ID, inputPower)

		// Update the line's power in the config
		for i, line := range config.Lines {
			if line.ID == n.ID {
				var ro = line.Ro
				var l = int(line.Length)
				var A = line.Area
				var lineResistence = lineResistence(ro, l, A)
				var Dm = geometricDistance(line.Drs, line.Dst, line.Drt)
				var re = equivalentRadius(line.R)
				var L = inductionOnLength(Dm, re)
				var XL = totalReactation(math.Pi*50*L, float64(line.Length)) // Reactanta liniei: ohm/km
				// var C = lineCapacity(Dm, line.R)
				var lineCurrent = (inputPower * 1000000) / (line.Voltage * 1000)
				var apparentPower = (line.Voltage * 1000) * lineCurrent
				// distante intre faze: Dab = 4m, Dbc = 4m, Dac = 4m
				// diametrul conductorului: d = 2cm (r = 0.01m)
				var activePowerLoseesPerLine = wattToMegawatt(powerLineLoss(lineCurrent, lineResistence))
				var reactivePowerLoseesPerLine = wattToMegawatt(reactivePowerLineLoss(lineCurrent, XL))
				var cosFi = (inputPower * 1000000) / apparentPower
				var sinFi = sinfi(cosFi)
				var Q = (line.Voltage * 1000) * lineCurrent * sinFi

				fmt.Printf("Reactive power on line: %.2f \n", Q)

				var lineInfoMessage = fmt.Sprintf("Line %s (%d km) has voltage %.2f kV , Active power losses per line %.3f, Reactive power losses per line %.3f \n", line.ID, line.Length, line.Voltage, activePowerLoseesPerLine, reactivePowerLoseesPerLine)

				lineInfoLog := LogEntry{
					Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
					ComponentID: line.ID,
					Message:     lineInfoMessage,
				}

				logs = append(logs, lineInfoLog)

				config.Lines[i].PowerTransfered = inputPower
				break
			}
		}
		calculatePowerFlow(n.ConnectedTo, inputPower, config, nodes, visited, consumersWithoutPower)

	case utils.Separator:
		if n.State == utils.StateOpen {
			fmt.Printf("Separator %s is open, stopping power flow.\n", n.ID)
			calculatePowerFlow(n.ConnectedTo, 0, config, nodes, visited, consumersWithoutPower)
			return
		}
		for i, separator := range config.Separators {
			if separator.ID == n.ID {
				var separatorMessage = fmt.Sprintf("Separator %s is in %s state \n", separator.ID, separator.State)

				separatorLog := LogEntry{
					Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
					ComponentID: separator.ID,
					Message:     separatorMessage,
				}

				logs = append(logs, separatorLog)

				config.Separators[i].State = n.State
				break
			}
		}
		calculatePowerFlow(n.ConnectedTo, inputPower, config, nodes, visited, consumersWithoutPower)

	case utils.Consumer:
		var isActive = isActive(n.ID, config)

		remainingPower := inputPower - n.PowerNeeded

		if inputPower < n.PowerNeeded {
			*consumersWithoutPower = append(*consumersWithoutPower, utils.ConsumerPowerDetails{
				ID:                   n.ID,
				RemainingPowerNeeded: remainingPower,
			})
		}

		if !isActive {
			calculatePowerFlow(n.ConnectedTo, 0, config, nodes, visited, consumersWithoutPower)
			return
		}

		fmt.Printf("Consumer %s received power: %.2f, remaining: %.2f\n", n.ID, inputPower, remainingPower)
		// Update the consumer's remaining power in the config
		for i, consumer := range config.Consumers {
			if consumer.ID == n.ID {
				config.Consumers[i].RemainingPower = remainingPower
				var consumerMessage = fmt.Sprintf("Consumer %s draws %.2f MW at %.2f kV\n", consumer.ID, consumer.PowerNeeded, consumer.Voltage)

				consumerLog := LogEntry{
					Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
					ComponentID: consumer.ID,
					Message:     consumerMessage,
				}

				logs = append(logs, consumerLog)

				break
			}
		}
		calculatePowerFlow(n.ConnectedTo, remainingPower, config, nodes, visited, consumersWithoutPower)

	default:
		fmt.Printf("Unhandled node type for ID: %s\n", id)
	}
}

// Funcția principală pentru calcul
func ComputeSystem(system utils.System) []LogEntry {
	fmt.Println("Calculating power flow for the system...")

	// Verifică sursa inițială
	sourcePower := system.Source.Power
	sourceVoltage := system.Source.Voltage

	var sourceMessage = fmt.Sprintf("Source %s supplying %.2f MW at %.2f kV\n", system.Source.ID, sourcePower, sourceVoltage)

	// fmt.Println(sourceMessage)

	sourceLog := LogEntry{
		Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
		ComponentID: system.Source.ID,
		Message:     sourceMessage,
	}

	logs = append(logs, sourceLog)

	// Map nodes for quick lookup
	nodes := map[string]interface{}{}
	visited := map[string]bool{} // Track visited nodes

	// Populate nodes
	nodes[system.Source.ID] = system.Source

	for _, as := range system.AdditionalSources {
		nodes[as.ID] = as
	}
	for _, t := range system.Transformers {
		nodes[t.ID] = t
	}
	for _, l := range system.Lines {
		nodes[l.ID] = l
	}
	for _, s := range system.Separators {
		nodes[s.ID] = s
	}
	for _, c := range system.Consumers {
		nodes[c.ID] = c
	}

	consumersWithoutPower := []utils.ConsumerPowerDetails{}

	// Start traversal from the source
	calculatePowerFlow(system.Source.ID, system.Source.Power, &system, nodes, visited, &consumersWithoutPower)

	// Reverse calculation from additional sources
	for _, source := range system.AdditionalSources {
		var isSeparatorClose = false
		for _, separator := range system.Separators {
			if separator.ConnectedTo == source.ID && separator.State == utils.StateClose {
				isSeparatorClose = true
			}
		}
		if isSeparatorClose {
			calculatePowerFlow(source.ID, source.Power, &system, nodes, visited, &consumersWithoutPower)
		}
	}

	if len(consumersWithoutPower) > 0 {
		fmt.Println("Consumers without power:")
		for _, consumerID := range consumersWithoutPower {
			var emptyConsumerMessage = fmt.Sprintf("Consumer %s needs more power: %.2f MW\n", consumerID.ID, math.Abs(consumerID.RemainingPowerNeeded))

			emptyConsumerMessageLog := LogEntry{
				Timestamp:   time.Now().Format("2006/01/02-15:04:05"),
				ComponentID: consumerID.ID,
				Message:     emptyConsumerMessage,
			}

			logs = append(logs, emptyConsumerMessageLog)
			fmt.Printf("- Consumer ID: %sn", consumerID.ID)
			fmt.Printf("- Consumer Power Needed: %fn", consumerID.RemainingPowerNeeded)
		}
	} else {
		fmt.Println("All consumers are powered.")
	}

	return logs
}
