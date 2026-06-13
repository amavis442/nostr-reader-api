package db

// Return to client API
type Pagination struct {
	Cursor    uint64 `json:"cursor"`
	Direction string `json:"direction"`
	PerPage   uint   `json:"per_page,omitempty" query:"per_page"`
	Since     uint   `json:"since"`
}

func (p *Pagination) SetCursor(cursor uint64) {
	p.Cursor = cursor
}

func (p *Pagination) GetCursor() uint64 {
	return p.Cursor
}

func (p *Pagination) SetDirection(direction string) {
	p.Direction = direction
}

func (p *Pagination) GetDirection() string {
	return p.Direction
}

/*
 * How many records per page should be shown
 */
func (p *Pagination) SetPerPage(perPage uint) {
	p.PerPage = perPage
}

func (p *Pagination) GetPerPage() uint {
	if p.PerPage == 0 {
		p.PerPage = 10
	}
	return p.PerPage
}

func (p *Pagination) SetSince(since uint) {
	p.Since = since
}

func (p *Pagination) GetSince() uint {
	if p.Since == 0 {
		p.Since = 1
	}
	return p.Since
}
