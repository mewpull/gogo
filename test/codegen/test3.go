var x int = -1

func print() {
	printInt x
	return
}

func main() {
	var x int = 0
	{
		var x int = 1
		{
			var x int = 2
			{
				var x int = 3
				{
					var x int = 4
					{
						var x int = 5
						{
							var x int = 6
							printInt x 
						}
						printInt x
					}
					printInt x
				}
				printInt x
			}
			printInt x
		}
		printInt x
	}
	printInt x
	return
}