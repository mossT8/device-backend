package request

type Sensor struct {
	Name           string                 `json:"name"`
	Code           string                 `json:"code"`
	UnitId         int64                  `json:"unitId"`
	DefaultConfig  map[string]interface{} `json:"defaultConfig"`
	ConfigRequried map[string]interface{} `json:"configRequried"`
}
