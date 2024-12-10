package main

func getEncoderFunction(flag byte) func() {
	switch flag {
	case 0:
		return StringEncoding
	case 1:
		return ListEndcoding
	case 2:
		return SetEncoding
	case 4:
		return HashEncoding
	}
	return nil
}

func LengthEncoding() {
}

func StringEncoding() {
}
func ListEndcoding() {}
func SetEncoding() {
}
func HashEncoding() {}
