func main() {
	x := [10]int{}
	for i:=0; i<10; i++ {
		x[i] = i
	}
	var a int 
	for i = 0; i<10; i++ {
		if i == 5 {
			break
		}
		a = x[i]
		printInt a 
	}
	return 
}