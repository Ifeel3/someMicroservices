package structs

type ItemStruct struct {
	Id    int32  `json:"id"`
	Info  string `json:"info"`
	Price int64  `json:"price"`
	Owner string `json:"owner"`
}
