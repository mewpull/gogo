func main() {
	var sum int = 0
	var y, n int
	scanInt y
	scanInt n
	for i := 0; i<n; i++ {
		sum = sum + y
	}
	printInt sum
	return
}