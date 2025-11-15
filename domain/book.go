package domain

type Book struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	Authors *string `json:"authors"`
	Series  *string `json:"series"`
	Serno   *string `json:"serno"`
	File    *string `json:"file"`
	Path    *string `json:"path"`
	Lang    *string `json:"lang"`
	Size    *uint   `json:"size"`
}
