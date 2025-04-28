package db

// Page 分页
type Page struct {
	disable    bool     //开启关闭分页功能，关闭分页时则查询全部数据
	orders     []string //排序规则，如果不设置默认按id倒排
	PageNumber int      `json:"page_number" biggerEq:"0" nilable:"" verf:""`             //页码，默认1，可外部传参
	PageSize   int      `json:"page_size" biggerEq:"0" nilable:"" verf:"" lowerEq:"500"` //分页大小，默认10，可外部传参
}

// DisablePage 禁用分页
func (p *Page) DisablePage() *Page {
	p.disable = true
	return p
}

// EnablePage 开启分页，默认就是开启的
func (p *Page) EnablePage() *Page {
	p.disable = false
	return p
}

// GetPage 获取页码
func (p *Page) GetPage() (offset int, limit int, disable bool) {
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	offset = (p.PageNumber - 1) * p.PageSize
	limit = p.PageSize
	disable = p.disable
	return
}

// SetOrders 设置排序字段
func (p *Page) SetOrders(orders ...string) *Page {
	p.orders = orders
	return p
}

// GetOrders 获取排序字段
func (p *Page) GetOrders() []string {
	if len(p.orders) == 0 {
		return []string{"id desc"}
	}
	return p.orders
}

// SplitIndex 计算分页的起始位置
func SplitIndex(pageNumber, pageSize, total int) (canSplit bool, start int, end int) {
	start = (pageNumber - 1) * pageSize //起始下标，包含
	end = start + pageSize              //结束下标，不包含
	if start < total {
		if end > total {
			end = total
		}
		canSplit = true
	}
	return
}
