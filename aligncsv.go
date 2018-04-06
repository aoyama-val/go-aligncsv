package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const progName = "aligncsv"

var csvText = ""
var maxLens = make([]int, 0, 100)
var widthCache = make(map[string]int)

func main() {
	readStdin()
	calcMaxWidth()
	align()
}

// 標準入力を全部読み込む
func readStdin() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	csvText = string(b)
}

// いったん全行ループして各カラムの幅を計算する
func calcMaxWidth() {
	reader := csv.NewReader(strings.NewReader(csvText))
	reader.Comma = ','
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない！
	lineno := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		lineno += 1

		diff := len(record) - len(maxLens)
		for i := 0; i < diff; i++ {
			maxLens = append(maxLens, 0)
		}

		for i, field := range record {
			w := getStringDisplayWidth(field)
			key := fmt.Sprintf("%d:%d", lineno, i)
			widthCache[key] = w
			if w > maxLens[i] {
				maxLens[i] = w
			}
		}
	}
}

// もう一度全行ループして各フィールドを整形して表示する
func align() {
	reader := csv.NewReader(strings.NewReader(csvText))
	reader.Comma = ','
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない！
	lineno := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		lineno += 1

		for i, field := range record {
			key := fmt.Sprintf("%d:%d", lineno, i)
			w := widthCache[key]
			diff := maxLens[i] - w
			if i != 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%s%s", field, strings.Repeat(" ", diff))
		}
		fmt.Print("\n")
	}
}

// 文字列の表示幅を返す
func getStringDisplayWidth(str string) int {
	len := 0
	for _, runeValue := range str {
		if runeValue > 0xff {
			len += 2
		} else {
			len += 1
		}
	}
	return len
}
