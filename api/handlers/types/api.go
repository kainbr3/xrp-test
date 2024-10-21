package types

type HealthStatus struct {
	App    string `json:"app"`
	Status string `json:"status"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type Result struct {
	Result string `json:"result"`
}

type ExecuteOperationResult struct {
	OperationId string `json:"operation_id" example:"6709778cf22e601d8921bd1a"`
	Error       error  `json:"error" example:""`
}
