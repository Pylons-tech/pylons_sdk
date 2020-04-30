package queriers

// query endpoints supported by the nameservice Querier
const (
	KeyAddrFromPubKey = "addr_from_pub_key"
)

// AddrResp holds the bech32 encoded address
type AddrResp struct {
	Bech32Addr string
}
