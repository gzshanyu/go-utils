package utils

type Paginator struct {
	Paged     int32
	Total     int32
	PageCount int32
	PageSize  int32
	PrevPage  int32
	LastPage  int32
}

func NewPaginator(paged int32, pageSize ...int32) *Paginator {
	var psize int32
	if pageSize != nil {
		if pageSize[0] > 0 {
			psize = pageSize[0]
		}
	}
	if paged < 1 {
		paged = 1
	}

	if psize < 1 {
		psize = 20
	}

	entity := &Paginator{
		PageSize: psize,
		Paged:    paged,
	}
	return entity
}

func (a *Paginator) Offset() int32 {
	offset := (a.Paged - 1) * a.PageSize
	return offset
}

func (a *Paginator) Limit() int32 {
	return a.PageSize
}

func (a *Paginator) ToPager(pbPager interface{}) *interface{} {
	a.PageCount = (a.Total + a.PageSize - 1) / a.PageSize
	a.LastPage = a.Paged + 1
	a.PrevPage = a.Paged - 1
	if a.LastPage > a.PageCount {
		a.LastPage = a.PageCount
	}
	if a.PrevPage < 1 {
		a.PrevPage = 1
	}

	StructCopy(pbPager, a)
	return &pbPager
}

func Top(top int32) int32 {
	if top > 0 {
		return top
	}
	return 20
}
