package iteration

func Repeat(character string, times int) (repeated string) {
	repeated = ""
	for range times {
		repeated += character
	}
	return
}
