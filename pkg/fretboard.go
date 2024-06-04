package pkg

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Mode int

const EASY Mode = 0
const HARD Mode = 1

func Fretboard(m Mode, head bool, yLeft, yRight, cnt int) {
	var score int
	var combo int

	var failedBlock map[location]int = make(map[location]int)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		printResult(failedBlock)
		log(score, head, yLeft, yRight, cnt, m)
	}()
	var cur = time.Now()
	var cut float64 = 0
	var prex int
	var prey int
	for i := 0; i < cnt; i++ {
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
		y := rand.Intn(yRight-yLeft) + yLeft
		for prex == x && prey == y {
			x = rand.Intn(6)
			y = rand.Intn(yRight-yLeft) + yLeft
		}
		prex = x
		prey = y
		printFretboard(x, y, head)
		var res string
		for {
			fmt.Scanf("%s", &res)
			if res != "" {
				break
			}
		}
		res = convert(res)
		if inSlice(getNote(x, y), res) {
			combo = calculateCombo(combo, 1)
			score += combo
			continue
		}
		combo = calculateCombo(combo, -1)
		score += combo

		if failedBlock[location{x: x, y: y}] == 0 {
			failedBlock[location{x: x, y: y}] = 1
		} else {
			failedBlock[location{x: x, y: y}]++
		}

		fmt.Printf("Answwer: %v", getNote(x, y))
		before := time.Now()
		fmt.Scanf("\n")
		after := time.Now()
		cut = after.Sub(before).Seconds()
	}
	fmt.Println(getResult(score, head, yLeft, yRight, cnt, m))
	log(score, head, yLeft, yRight, cnt, m)
}

func calculateCombo(combo int, incr int) int {
	if combo*incr >= 0 {
		combo += incr
		return combo
	}

	combo = incr
	return combo
}

func printResult(res map[location]int) {
	fmt.Println("错误信息")
	for i := 0; i < 6; i++ {
		fmt.Printf(head[i])
		for j := 1; j < 13; j++ {
			var str string
			if res[location{x: i, y: j}] == 0 {
				str = "---|"
			} else {
				str = fmt.Sprintf("%03d|", res[location{x: i, y: j}])
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

func convert(s string) string {
	s = strings.Title(s)
	s = strings.ReplaceAll(s, "1", "b")
	s = strings.ReplaceAll(s, "2", "#")
	return s
}

func log(score int, head bool, yLeft, yRight, cnt int, mode Mode) {
	h, _ := os.UserHomeDir()
	f, err := os.OpenFile(h+"/.fretboard.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		os.Stdout.Write([]byte(err.Error()))
		os.Exit(0)
	}

	defer f.Close()
	f.Write([]byte(getResult(score, head, yLeft, yRight, cnt, mode)))
	os.Exit(0)

}

func getResult(score int, head bool, yLeft, yRight, cnt int, mode Mode) string {
	modeStr := "easy"
	if mode == HARD {
		modeStr = "hard"
	}

	headStr := "序号"
	if head {
		headStr = "音名"
	}
	levelStr := fmt.Sprintf("%v %v", yLeft, yRight)

	res := fmt.Sprintf("mode: %v\tscore: %v\thead: %v\t品数: %v\t测试次数: %v\n", modeStr, score, headStr, levelStr, cnt)
	return res

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
	index := (offset[x] + y) % 12
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
	3: 7,
	4: 12,
}

var head = map[int]string{
	0: "E:",
	1: "B:",
	2: "G:",
	3: "D:",
	4: "A:",
	5: "E:",
}

// 00 01 02 03 04 05 06 07 08 09 10 11 12
//
//	c c#  d d#  e  f f#  g g#  a a#  b b#
var offset = map[int]int{
	0: 28,
	1: 23,
	2: 19,
	3: 14,
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

var interval = map[int]string{
	0:  "P1",
	1:  "m2",
	2:  "M2",
	3:  "m3",
	4:  "M3",
	5:  "P4",
	6:  "A4",
	7:  "P5",
	8:  "m6",
	9:  "M6",
	10: "m7",
	11: "M7",
	12: "P8",
}

func printFretboard(x, y int, h bool) {

	for i := 0; i < 6; i++ {
		if h {
			fmt.Printf(head[i])
		} else {
			fmt.Printf(fmt.Sprintf("%d", i+1) + ":")

		}
		for j := 1; j < 13; j++ {
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

func PrintInterval() {
	for i := 0; i < 6; i++ {
		for j := 0; j < 13; j++ {
			printFretboardWithInterval(i, j)
		}
	}
}

func PrintSpecificInterval(x, y int, single bool) {
	if single {
		printFretboardWithSingleInterval(x, y)
		return
	}
	printFretboardWithInterval(x, y)
}

type location struct {
	x int
	y int
}

func printDict(d map[int]location) {
	for k, l := range d {
		fmt.Printf("%d ", k)
		fmt.Printf("%d ", l.x)
		fmt.Printf("%d ", l.y)
		fmt.Println()
	}

}

func printFretboardWithSingleInterval(x, y int) {
	res := make(map[int]location)
	for i := 0; i < 6; i++ {
		for j := 0; j < 13; j++ {
			interval := getSignedNotesDistance(x, y, i, j)
			if abs(interval) > 12 {
				continue
			}
			if _, exist := res[interval]; !exist {
				res[interval] = location{
					x: i,
					y: j,
				}
				continue
			}
			oldDistance := getFretsDistance(x, y, res[interval].x, res[interval].y)
			newDistance := getFretsDistance(x, y, i, j)
			if oldDistance > newDistance {
				res[interval] = location{
					x: i,
					y: j,
				}
			}
		}
	}

	revert := make(map[location]int)
	for k, l := range res {
		revert[l] = k
	}

	for i := 0; i < 6; i++ {
		for j := 0; j < 13; j++ {
			if l, ok := revert[location{x: i, y: j}]; ok {
				if l == 0 {
					fmt.Printf(" **:")
					continue
				}
				fmt.Printf(" %s:", interval[abs(l)])
				continue
			}
			fmt.Print("___:")

		}
		fmt.Println()
	}
	for i := 0; i < 13; i++ {
		fmt.Printf(" %02d ", i)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()

}

func printFretboardWithInterval(x, y int) {
	for i := 0; i < 6; i++ {
		res := abs(getSignedNotesDistance(x, y, i, 0))
		if res <= 12 {
			if res == 0 {
				fmt.Printf("**:")
			} else {
				fmt.Printf("%s:", interval[res])
			}
		} else {
			fmt.Print("__:")
		}

		for j := 1; j < 13; j++ {
			res := abs(getSignedNotesDistance(x, y, i, j))
			if res <= 12 {
				if res == 0 {
					fmt.Printf(" **:")
				} else {
					fmt.Printf(" %s:", interval[res])
				}
			} else {
				fmt.Print("___:")
			}
		}
		fmt.Println()
	}
	fmt.Printf("   ")
	for i := 0; i < 12; i++ {
		fmt.Printf(" %02d ", i+1)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func getSignedNotesDistance(x, y, x1, y1 int) int {
	a := offset[x] + y
	b := offset[x1] + y1
	return a - b
}

// getFretsDistance function  
func getFretsDistance(x, y, x1, y1 int) int {
	return abs(x-x1) + abs(y-y1)
}

func abs(v int) int {
	if v > 0 {
		return v
	}
	return v * -1
}
