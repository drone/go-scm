package fake

func paginated(page, size, items int) (start, end int) {
	start = (page - 1) * size
	if start > items {
		start = items
	}
	end = start + size
	if end > items {
		end = items
	}
	return
}
