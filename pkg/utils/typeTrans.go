package utils

func ToPtrSlice[T any](data []T) []*T {
	return mapTrans[T, *T](data, func(idx int, v T) *T {
		return &v
	})
}

func ToPtr[T any](t T) *T {
	return &t
}

func mapTrans[T any, S any](vs []T, fc func(idx int, v T) S) []S {
	res := make([]S, len(vs))
	for i, v := range vs {
		res[i] = fc(i, v)
	}
	return res
}
