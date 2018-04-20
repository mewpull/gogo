type Rect struct {
	length int
	breadth int
}

type Circle struct {
	radius int
}

func main() {
	var length, breadth int = 5, 3
	x := Rect{length, breadth}
	y := Circle{5}
	var area_of_rect int = x.length * x.breadth
	var area_of_circle int = 3 * y.radius * y.radius
	printInt area_of_circle
	printInt area_of_rect
	return
}
