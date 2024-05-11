package pkg

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type Mode int

const EASY Mode = 0
const HARD Mode = 1

func Fretboard(m Mode, head bool, level int) {
	var score int
	var combo int

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		h, _ := os.UserHomeDir()
		f, err := os.OpenFile(h+"/.fretboard.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			os.Stdout.Write([]byte(err.Error()))
			os.Exit(0)
		}

		defer f.Close()
		s := fmt.Sprintf("score: %v\n", score)
		f.Write([]byte(s))
		os.Exit(0)
	}()
	var cur = time.Now()
	var cut float64 = 0
	for {
		punishScore := int(time.Now().Sub(cur).Seconds() - cut)
		cut = 0
		cur = time.Now()

		if m == EASY {
			score = score - easy(punishScore)
		} else {
			score = score - hard(punishScore)
		}

		fmt.Printf("Score: %d, Combo: %d\n", score, combo)
		fmt.Println()
		x := rand.Intn(6)
		y := rand.Intn(getLevel(level))
		printFretboard(x, y, head)
		var res string
		fmt.Scanf("%s", &res)
		if inSlice(getNote(x, y), res) {
			combo++
			score += combo
			continue
		}
		fmt.Printf("Answwer: %v", getNote(x, y))
		before := time.Now()
		fmt.Scanf("\n")
		after := time.Now()
		cut = after.Sub(before).Seconds()
		combo = 0
	}
}

// 1 3
// 2 5
// 3 10
// 4 12
// 5 24
func getLevel(l int) int {
	return level[l]
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

var level = map[int]int{
	1: 3,
	2: 5,
	3: 10,
	4: 12,
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

func printFretboard(x, y int, h bool) {

	for i := 0; i < 6; i++ {
		if h {
			fmt.Printf(head[i] + ":")
		} else {
			fmt.Printf(fmt.Sprintf("%d", i+1) + ":")

		}
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
