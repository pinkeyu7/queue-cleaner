package queue_api

type Headers struct {
	ControlMsg string `json:"controlMsg"`
}

type Properties struct {
	Headers Headers `json:"headers"`
}

type CloseQueueBody struct {
	Properties      Properties `json:"properties"`
	RoutingKey      string     `json:"routing_key"`
	DeliveryMode    string     `json:"delivery_mode"`
	Payload         string     `json:"payload"`
	PayloadEncoding string     `json:"payload_encoding"`
}
