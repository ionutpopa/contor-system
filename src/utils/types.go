package utils

// Structuri si tipuri pentru reprezentarea sistemului

// Custom type for state open and close
type StateType string
type TransformerType string

// Define open and close on StateType
const (
	StateOpen  StateType = "open"
	StateClose StateType = "close"
)

const (
	TransformerTypeMeasure TransformerType = "measure"
	TransformerTypePower   TransformerType = "power"
)

type Source struct {
	ID              string  `json:"id"`
	Power           float64 `json:"power"` // MW
	Voltage         float64 `json:"voltage"`
	ConnectedTo     string  `json:"connectedTo"`
	AdditionalPower float64 `json:"additionalPower"`
}

type Transformer struct {
	ID              string          `json:"id"`
	InputVoltage    float64         `json:"inputVoltage"`
	OutputVoltage   float64         `json:"outputVoltage"`
	ConnectedTo     string          `json:"connectedTo"`
	Type            TransformerType `json:"type"`
	Efficency       float64         `json:"efficency"`
	ApparentPower   float64         `json:"apparentPower"` // MW
	CooperLosses    float64         `json:"cooperLosses"`  // KW
	SteelLosses     float64         `json:"steelLosses"`   // KW
	PowerTransfered float64         `json:"powerTransfered"`
}

type Line struct {
	ID                  string  `json:"id"`
	Voltage             float64 `json:"voltage"`
	Length              int     `json:"length"` // km
	ConnectedTo         string  `json:"connectedTo"`
	Area                float64 `json:"area"`
	Currnet             float64 `json:"current"`
	Ro                  float64 `json:"ro"`
	Drs                 float64 `json:"Drs"`
	Dst                 float64 `json:"Dst"`
	Drt                 float64 `json:"Drt"`
	ConductorDiameter   float64 `json:"conductorDiameter"`
	R                   float64 `json:"r"`
	PowerTransfered     float64 `json:"powerTransfered"`
	ReactivePowerLosses float64 `json:"reactivePowerLosses"`
	ActivePowerLosses   float64 `json:"activePowerLosses"`
}

type Consumer struct {
	ID             string  `json:"id"`
	PowerNeeded    float64 `json:"powerNeeded"` // MW
	Voltage        float64 `json:"voltage"`
	ConnectedTo    string  `json:"connectedTo,omitempty"`
	RemainingPower float64 `json:"remainingPower"`
}

type Separator struct {
	ConnectsFrom string    `json:"connectsFrom"`
	ID           string    `json:"id"`
	State        StateType `json:"state"`
	ConnectedTo  string    `json:"connectedTo"`
}

type System struct {
	Source            Source        `json:"source"`
	Transformers      []Transformer `json:"transformers"`
	Lines             []Line        `json:"lines"`
	Consumers         []Consumer    `json:"consumers"`
	Separators        []Separator   `json:"separators"`
	AdditionalSources []Source      `json:"additionalSources"`
}

// This type struct also represents the parquet schema which is pretty cool
type LogEntry struct {
	Timestamp   string `parquet:"name=timestamp, type=BYTE_ARRAY, convertedtype=UTF8"`
	ComponentID string `parquet:"name=component_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Message     string `parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type ConnectedElement struct {
	ID      string
	Details interface{} // Store the full details of the element
}

type ConsumerPowerDetails struct {
	id             string
	remainingPower float64
}

type TransferedPowerDetails struct {
	powerTransfered float64
}

type LinePowerDetails struct {
	TransferedPowerDetails
	id string
}

type TransformerPowerDetails struct {
	TransferedPowerDetails
	id string
}
