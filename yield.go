package util

func YieldArr(data []interface{}) func() (interface{}, bool) {
	i := 0
	return func() (interface{}, bool) {
		if i < len(data) {
			item := data[i]
			i++
			return item, true
		}
		return nil, false
	}
}
