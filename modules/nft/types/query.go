package types

type QueryPagesParams struct {
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

// QueryCollectionParams defines the params for queries:
type QueryCollectionParams struct {
	Denom string `json:"Denom"`
	QueryPagesParams
}
