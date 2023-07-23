package overWs

type Pack struct {
	Operation int                    `json:"operation"`
	RawBody   map[string]interface{} `json:"body"`
}
