package main

func main() {
	bookIDs := make([]int64, 0, 2000)
	wordCounts := make([]int64, 0, 2000)
	bookIntros := make([][]float32, 0, 2000)
	for i := 0; i < 2000; i++ {
		bookIDs = append(bookIDs, int64(i))
		wordCounts = append(wordCounts, int64(i+10000))
		v := make([]float32, 0, 2)
		for j := 0; j < 2; j++ {
			v = append(v, rand.Float32())
		}
		bookIntros = append(bookIntros, v)
	}
}
