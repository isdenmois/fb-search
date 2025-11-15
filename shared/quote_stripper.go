package shared

type QuoteStripper struct{}

func (q QuoteStripper) Transform(dst, src []byte, atEOF bool) (int, int, error) {
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
func (q QuoteStripper) Reset() {}
