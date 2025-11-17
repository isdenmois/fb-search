package domain

type Book struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	File    string  `json:"file"`
	Authors *string `json:"authors"`
	Series  *string `json:"series"`
	Serno   *string `json:"serno"`
	Lang    *string `json:"lang"`
	Size    *uint   `json:"size"`
}
