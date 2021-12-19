package shortten

import "context"

func Encode(ctx context.Context, plain string) (string, error) {
	return plain, nil
}

func Decode(ctx context.Context, enc string) (string, error) {
	return enc, nil
}
