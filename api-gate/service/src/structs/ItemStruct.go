package structs

type ItemStruct struct {
	Item  string `json:"item"`
	Info  string `json:"info"`
	Price int64  `json:"price"`
	Owner string `json:"Owner"`
}
