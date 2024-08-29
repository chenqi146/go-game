package page

type Pageable struct {
	PageSize int64 `json:"pageSize"`
	PageNum  int64 `json:"pageNum"`
}

type Page struct {
	Total    int64       `json:"total"`
	Records  interface{} `json:"records"`
	Pageable Pageable    `json:"pageable"`
}

func NewPageable(pageSize int64, pageNum int64) *Pageable {
	return &Pageable{
		PageSize: pageSize,
		PageNum:  pageNum,
	}
}

func NewPage(total int64, records interface{}, pageable Pageable) *Page {
	return &Page{
		Total:    total,
		Records:  records,
		Pageable: pageable,
	}
}
