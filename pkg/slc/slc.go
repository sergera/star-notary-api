package slc

func Map[IN any, OUT any](slc []IN, f func(IN) OUT) []OUT {
	mapped := make([]OUT, len(slc))

	for i, e := range slc {
		mapped[i] = f(e)
	}

	return mapped
}
