package queriers

// query endpoints supported by the nameservice Querier
const (
	KeyAddrFromPubKey = "addr_from_pub_key"
)

// AddrResponse holds the bech32 encoded address
type AddrResponse struct {
	Bech32Addr string
}
