package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

type Mode int

const EASY Mode = 0
const HARD Mode = 1

func Fretboard(m Mode) {
	var score int
	var combo int
	var cur = time.Now()
	for {
		punishScore := int(time.Now().Sub(cur).Seconds())
		cur = time.Now()

		if m == EASY {
			score = score - easy(punishScore)
		} else {
			score = score - hard(punishScore)
		}

		fmt.Printf("Score: %d, Combo: %d\n", score, combo)
		fmt.Println()
		x := rand.Intn(6)
		y := rand.Intn(12)
		printFretboard(x, y)
		var res string
		fmt.Scanf("%s", &res)
		if inSlice(getNote(x, y), res) {
			combo++
			score += combo
			continue
		}
		fmt.Printf("Answwer: %v", getNote(x, y))
		combo = 0
	}
}

func inSlice(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func easy(seconds int) int {
	return seconds
}

func hard(seconds int) int {
	if seconds <= 0 {
		return 0
	}
	return seconds + hard(seconds-1)
}

func getNote(x, y int) []string {
	res := make([]string, 0)
	index := (offset[x] + y + 1) % 12
	if dict[index] != "" {
		res = append(res, dict[index])
		return res
	}
	res = append(res, dict[index+1]+"b")
	res = append(res, dict[index-1]+"#")
	return res
}

var head = map[int]string{
	0: "E",
	1: "B",
	2: "G",
	3: "D",
	4: "A",
	5: "E",
}

var offset = map[int]int{
	0: 4,
	1: 11,
	2: 7,
	3: 2,
	4: 9,
	5: 4,
}

var dict = map[int]string{
	0:  "C",
	1:  "",
	2:  "D",
	3:  "",
	4:  "E",
	5:  "F",
	6:  "",
	7:  "G",
	8:  "",
	9:  "A",
	10: "",
	11: "B",
}

func printFretboard(x, y int) {

	for i := 0; i < 6; i++ {
		fmt.Printf(head[i] + ":")
		for j := 0; j < 12; j++ {
			str := "---|"
			if x == i && y == j {
				str = "???|"
			}
			fmt.Printf(str)
		}
		fmt.Println()
	}
	fmt.Printf("  ")
	for i := 0; i < 12; i++ {
		fmt.Printf(" %02d ", i+1)
	}
	fmt.Println()

}
