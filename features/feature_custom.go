package features

func NewBytesListFeature(bytes ...[]byte) *Feature {
	return &Feature{
		Kind: &Feature_BytesList{BytesList: &BytesList{
			Value: bytes,
		}},
	}
}

func NewInt64ListFeature(ints ...int64) *Feature {
	return &Feature{
		Kind: &Feature_Int64List{Int64List: &Int64List{
			Value: ints,
		}},
	}
}

func NewFloat32ListFeature(floats ...float32) *Feature {
	return &Feature{
		Kind: &Feature_FloatList{FloatList: &FloatList{
			Value: floats,
		}},
	}
}

func LabelFeature(name string) *Feature {
	return &Feature{
		Kind: &Feature_BytesList{BytesList: &BytesList{
			Value: [][]byte{[]byte(name)},
		}},
	}
}

func ImageSizeFeature(x, y int64) *Feature {
	return NewInt64ListFeature(x, y)
}
