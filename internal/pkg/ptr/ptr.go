package ptr

func String(s string) *string {
	return &s
}

func Int32(v int32) *int32 {
	return &v
}
