package convertors

func SimplePaginationToLimitOffset(page, pageSize int) (int, int) {
	offset := (page - 1) * pageSize
	return pageSize, offset
}
