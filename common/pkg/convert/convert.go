package convert

type Int interface {
	int | int32
}

func ConvertIDs[T, K Int](inputs []T) []K {
	if len(inputs) == 0 {
		return nil
	}

	var res = make([]K, 0, len(inputs))
	for _, n := range inputs {
		res = append(res, K(n))
	}

	return res
}
