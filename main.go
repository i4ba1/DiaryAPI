package main

/**
go build && ./DiaryAPI
 */
func main() {
	a := App{}
	a.Initialize()
	a.Run(":8789")
}

