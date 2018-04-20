func sum(n int) int {
        if n == 0 {
                return 0
        }
        x := sum(n-1)
        return 2 + x
}

func main() {
        x := sum(5)
        printInt x
        return
}
