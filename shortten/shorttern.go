package shortten

import "context"

func Encode(ctx context.Context, plain string) (string, error) {
	//TODO:
	return reverseString(plain), nil
}

func Decode(ctx context.Context, enc string) (string, error) {
	//TODO:
	return reverseString(enc), nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
