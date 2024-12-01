package computing

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"contor-system/src/utils"
)

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

func getPowerOfNode(connectedTo []utils.ConnectedElement, system utils.System) float64 {
	for _, connection := range connectedTo {
		var element = connection

		if element.ID == "source1" {
			// Suntem in primul element, deci o sursa, deci puterea va fi cea instalata
			return element.Details.(utils.Source).Power
		}

		// Verificam daca tipul potrivit lui Details este transformer
		if transformer, ok := connection.Details.(utils.Transformer); ok {
			// Suntem in unul din transformatoare
			fmt.Println(transformer)
		}
	}

	return 0.0
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

// Pierderile din trafo se monitorizeaza doar alea din amonte, (110), deci pierderile in fier influenteaza doar partea de 110kv, nu 20kv, 5A in secundar

// Funcția principală pentru calcul
func ComputeSystem(system utils.System) []LogEntry {
	var logs []LogEntry
	// var totalActivePowerLosses20KV float64
	// var totalReactivePowerLosses20KV float64
	// var totalActivePowerLosses110KV float64
	// var totalReactivePowerLosses110KV float64
	// var totalActivePowerLosses220KV float64
	// var totalReactivePowerLosses220KV float64
	// var totalActivePowerLosses400KV float64
	// var totalReactivePowerLosses400KV float64

	// Create the power map
	powerMap := PowerMap(system)

	fmt.Println("Calculating power flow for the system...")

	// Verifică sursa inițială
	sourcePower := system.Source.Power
	sourceVoltage := system.Source.Voltage

	// Modify the power value of the "Source" key
	// Access the slice associated with "Source" and modify the power field
	if sourceFromMap, ok := powerMap["Source"]; ok && len(sourceFromMap) > 0 {
		// First position in Source map will always be the first source of the config file
		sourceFromMap["0"]["power"] = sourcePower
	}

	var sourceMessage = fmt.Sprintf("Source %s supplying %.2f MW at %.2f kV\n", system.Source.ID, sourcePower, sourceVoltage)

	fmt.Println(sourceMessage)

	sourceLog := LogEntry{
		Timestamp:   time.Now().String(),
		ComponentID: system.Source.ID,
		Message:     sourceMessage,
	}

	logs = append(logs, sourceLog)

	// Pierderile in fier si cupru nu se aplica la tensiunea de 20kV
	// Pierderile in fier si cupru se aplica in schimb la inalta tensiune, cum ar fi 110, 220, 400 kV
	var transformerLosses float64

	// Traversează transformatoarele și liniile
	for transformerIndex, transformer := range system.Transformers {
		var transformerMessage = fmt.Sprintf("Transformer %s steps %.2f kV to %.2f kV, type: %s \n", transformer.ID, transformer.InputVoltage, transformer.OutputVoltage, transformer.Type)
		var totalCooperAndSteelLosses = (transformer.SteelLosses / 1000) + (transformer.CooperLosses / 1000)
		fmt.Println(transformerMessage)

		if transformerFromPowerMap, ok := powerMap["Transformers"]; ok && len(transformerFromPowerMap) > 0 {
			var transformerIndexStr = strconv.Itoa(transformerIndex)
			fmt.Println("power that gets thorugh value of "+transformerIndexStr+" transformer", transformerFromPowerMap[transformerIndexStr]["power"])
			// If the tranformer index is 0 then the first power to substrant the transformer losses from is the Source power
			if transformerIndex == 0 {
				transformerFromPowerMap[transformerIndexStr]["power"] = powerMap["Source"]["0"]["power"] - totalCooperAndSteelLosses
			} else {
				var lastConnection = findConnectedTo(system, transformer.ID)
				if len(lastConnection) > 0 {
					var lastTransformerIndexStr = strconv.Itoa(transformerIndex - 1)
					fmt.Println("Ultima conexiune a lui", transformer.ID, lastConnection[0].ID, lastConnection[0].Details)
					transformerFromPowerMap[transformerIndexStr]["power"] = transformerFromPowerMap[lastTransformerIndexStr]["power"] - totalCooperAndSteelLosses
				}
			}
		}

		// Pierderile in cupru si fier vor fi masurate in kW

		transformerLog := LogEntry{
			Timestamp:   time.Now().String(),
			ComponentID: transformer.ID,
			Message:     transformerMessage,
		}

		if transformer.Type == utils.TransformerTypeMeasure {
			// measure
		}

		if transformer.Type == utils.TransformerTypePower {
			// More then 20KV -> 110, 220, 400
			if transformer.OutputVoltage == 110 {
				var connectedTo = findConnectedTo(system, transformer.ID)

				for _, connection := range connectedTo {
					if connection.ID == "source1" {

					}
				}

				fmt.Println(transformerLosses)
				// // Aflam ce putere a mai ramas in nod si scadem pierderile din fier si cupru
				// var powerInput float64 = getPowerOfNode(connectedTo, system) - totalCooperAndSteelLosses

				// transformerLosses += transformerLossesBasedOnEfficency(powerInput, transformer.Efficency)
			}
		}

		logs = append(logs, transformerLog)
	}

	var activePowerLoseesPerLine float64
	var reactivePowerLoseesPerLine float64

	for _, line := range system.Lines {
		var ro = line.Ro
		var l = int(line.Length)
		var A = line.Area
		var current = line.Currnet
		var lineResistence = lineResistence(ro, l, A)
		activePowerLoseesPerLine += wattToMegawatt(powerLineLoss(current, lineResistence))

		var Dm = geometricDistance(line.Drs, line.Dst, line.Drt)
		var re = equivalentRadius(line.R)
		var L = inductionOnLength(Dm, re)
		var XL = totalReactation(math.Pi*50*L, float64(line.Length)) // Reactanta liniei

		// distante intre faze: Dab = 4m, Dbc = 4m, Dac = 4m
		// diametrul conductorului: d = 2cm (r = 0.01m)
		reactivePowerLoseesPerLine += wattToMegawatt(reactivePowerLineLoss(current, XL))

		var lineInfoMessage = fmt.Sprintf("Line %s (%d km) has voltage %.2f kV\n", line.ID, line.Length, line.Voltage)
		var linePowerLosses = fmt.Sprintf("Active power losses per line: %.3f, Reactive power losses per line %.3f \n", activePowerLoseesPerLine, reactivePowerLoseesPerLine)

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
		var consumerMessage = fmt.Sprintf("Consumer %s draws %.2f MW at %.2f kV\n", consumer.ID, consumer.PowerNeeded, consumer.Voltage)
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

		if separator.State == utils.StateClose {
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

	fmt.Println(powerMap)
	fmt.Println("---")

	return logs
}
