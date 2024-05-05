package apiv1

func ToRef[T any](v T) *T {
	return &v
}
