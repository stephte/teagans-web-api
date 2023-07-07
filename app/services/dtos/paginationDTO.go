package dtos

import (
	"math"
)

type PaginationDTO struct {
	sort			string
	page			int
	limit			int

	totalRows		int64
	totalPages		int

	firstPage		string
	lastPage		string
	previousPage	string
	nextPage		string

	rows			interface{}
}

func(this PaginationDTO) GetOffset() int {
	return (this.page-1) * this.limit
}

func(this *PaginationDTO) Init(sort string, page int, limit int) {
	this.SetSort(sort)
	this.SetPage(page)
	this.SetLimit(limit)
}

// call this once TotalRows is set
func(this *PaginationDTO) GenAndSetData(path string) {
	this.SetTotalPages() // do this first
	this.SetFirstPage(path)
	this.SetLastPage(path)
	this.SetNextPage(path)
	this.SetPreviousPage(path)
}

// ----- Setters -----

func(this *PaginationDTO) SetPage(page int) {
	if page < 1 {
		page = 1
	}

	this.page = page
}

func(this PaginationDTO) GetPage() int {
	return this.page
}

func(this *PaginationDTO) SetSort(sort string) {
	// validate sort string? (tested for SQL injection... is safe w/o validation, but will return error when trying to execute order query if field doesn't exist)
	if sort == "" {
		sort = "created_at desc"
	}

	this.sort = sort
}

func(this PaginationDTO) GetSort() string {
	return this.sort
}

func(this *PaginationDTO) SetLimit(limit int) {
	if limit < 1 {
		limit = 10 // default to 10
	} else if limit > 100 { // max limit is 100 records per page
		limit = 100
	}

	this.limit = limit
}

func(this PaginationDTO) GetLimit() int {
	return this.limit
}


// requires TotalRows and limit to be set
func(this *PaginationDTO) SetTotalPages() {
	this.totalPages = int(math.Ceil(float64(this.totalRows)/float64(this.limit)))
}

func(this *PaginationDTO) SetFirstPage(path string) {
	this.firstPage = genQueryStr(path, this.limit, 1, this.sort)
}

// requires TotalPages to be set
func(this *PaginationDTO) SetLastPage(path string) {
	this.lastPage = genQueryStr(path, this.limit, this.totalPages, this.sort)
}

// requires TotalPages to be set
func(this *PaginationDTO) SetNextPage(path string) {
	nextPage := this.page + 1
	if nextPage > this.totalPages {
		return
	}

	this.nextPage = genQueryStr(path, this.limit, nextPage, this.sort)
}

func(this *PaginationDTO) SetPreviousPage(path string) {
	prevPage := this.page - 1
	if prevPage < 1 {
		return
	}

	this.previousPage = genQueryStr(path, this.limit, prevPage, this.sort)
}

func(this *PaginationDTO) SetTotalRows(rowCount int64) {
	this.totalRows = rowCount
}

func(this *PaginationDTO) SetRows(rows interface{}) {
	this.rows = rows
}


// ----- For returning pagination to the UI -----


type PageResponseDTO struct {
	Sort						string				`json:"sort"`
	Page						int						`json:"page"`
	Limit						int						`json:"limit"`
	TotalRows				int64					`json:"totalRows"`
	TotalPages			int						`json:"totalPages"`
	FirstPage				string				`json:"firstPage"`
	LastPage				string				`json:"lastPage"`
	PreviousPage		string				`json:"previousPage"`
	NextPage				string				`json:"nextPage"`
	Rows						interface{}		`json:"rows"`
}

// custom method to allow private fields to send with json
func(this PaginationDTO) GetPageResponse() PageResponseDTO {
	return PageResponseDTO {
		Sort:						this.sort,
		Page:						this.page,
		Limit:					this.limit,
		TotalRows:			this.totalRows,
		TotalPages:			this.totalPages,
		FirstPage:			this.firstPage,
		LastPage:				this.lastPage,
		PreviousPage:		this.previousPage,
		NextPage:				this.nextPage,
		Rows:						this.rows,
	}
}
