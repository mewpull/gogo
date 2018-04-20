func sqrt(x int) int {
	y := x*x
	return y
}

func main() {
	var n int
	scanInt n
	var x int = sqrt(n)
	printInt x
	return
}