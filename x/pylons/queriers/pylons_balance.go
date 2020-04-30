package queriers

const (
	KeyPylonsBalance = "balance"
)

type QueryResBalance struct {
	Balance int64 `json:"balance"`
}
