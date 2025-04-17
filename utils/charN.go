package utils

func CharN(n int, char string) string {
	// Make sure we dont use 0
	if n == 0 {
		n = 1
	}
	out := ""
	for i := 0; i < n; i++ {
		out += char
	}
	return string(out)
}
