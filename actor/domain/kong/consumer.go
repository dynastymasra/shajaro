package kong

type (
	Consumer struct {
		ID       string `json:"id,omitempty"`
		Username string `json:"username"`
		CustomID string `json:"custom_id"`
	}

	Konger interface {
		CreateConsumer(Consumer) (*Consumer, int, error)
	}
)
