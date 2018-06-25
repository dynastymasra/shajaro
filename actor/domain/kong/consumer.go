package kong

type (
	Kong struct {
		ID       string `json:"id,omitempty"`
		Username string `json:"username"`
		CustomID string `json:"custom_id"`
	}

	Consumer interface {
		CreateConsumer(Kong) (*Kong, int, error)
		DeleteConsumer(string) (int, error)
	}
)
