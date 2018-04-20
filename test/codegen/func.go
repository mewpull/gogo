func temp(a, c int) (int, int) {
	return a, c
}

// func noRet() {
// 	a := 1
// 	printInt a
// 	return
// }
//
func main() {
	x, y := temp(1, 2)
	printInt y
	// noRet()
	return
}
