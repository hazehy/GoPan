package logic

import "gopan/gopan/define"

const adminListMaxPageSize = 100

// normalizePageAndSize applies consistent defaults for admin list endpoints.
func normalizePageAndSize(page, size int) (normalizedPage, normalizedSize, offset int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = define.PageSize
	}
	if size > adminListMaxPageSize {
		size = adminListMaxPageSize
	}
	return page, size, (page - 1) * size
}
