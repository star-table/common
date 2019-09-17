package strs

func Len(str string) int{
	return len([]rune(str))
}