package main

type quoteStripper struct{}

func (q quoteStripper) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	di, si := 0, 0
	for si < len(src) {
		if src[si] != '"' {
			dst[di] = src[si]
			di++
		}
		si++
	}
	return di, si, nil
}
func (q quoteStripper) Reset() {}
