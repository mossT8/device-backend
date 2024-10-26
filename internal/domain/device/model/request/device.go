package request

type Device struct {
	SerialNumber string                 `json:"serialNumber"`
	Name         string                 `json:"name"`
	ModelId      int64                  `json:"modelId"`
	ModelConfig  map[string]interface{} `json:"modelConfig"`
}
