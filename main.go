package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sudoku/files"
)

func main() {
	f, err := os.Open("puzzle.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := bufio.NewReader(f)
	data, err := file.Peek(81)
	if err != nil {
		panic(err)
	}
	fmt.Printf("81 bytes: %s\n", string(data))

	parsedPuzzle := ParsePuzzle(data)

	fullPuzzle := files.NewPuzzle(parsedPuzzle)
	solution := fullPuzzle.Solve()

	OutputPuzzle([]byte(fullPuzzle.ToString()))
	OutputPuzzle([]byte(solution.ToString()))

	solved := fullPuzzle.IsSolution(solution)
	Answer := "NO"
	if solved {
		Answer = "YES"
	}
	fmt.Println("Solved? ", Answer)
	fmt.Println("Number of solutions: ", files.Enumarate(parsedPuzzle))

	newPuzzle, newSolution := files.GeneratePuzzle()
	fmt.Println(newPuzzle.ToString())
	OutputPuzzle([]byte(newPuzzle.ToString()))
	fmt.Println(newPuzzle.ToVerboseString())
	OutputPuzzle([]byte(newSolution.ToString()))

}

func OutputPuzzle(puzzle []byte) {
	var output []byte
	counter := 1
	separator := "-------+-------+-------"
	for _, character := range puzzle {
		if (counter-1)%9 == 0 {
			output = append(output, ' ')
		}
		output = append(output, character, ' ')
		switch {
		case counter%81 == 0:
			output = append(output, '\n')
		case counter%27 == 0:
			output = append(output, '\n')
			output = append(output, []byte(separator)...)
			output = append(output, '\n')
		case counter%9 == 0:
			output = append(output, '\n')
		case counter%3 == 0:
			output = append(output, '|', ' ')
		}
		counter++
	}

	buf := bytes.NewBuffer(output)
	fmt.Println(buf.String())
}

func ParsePuzzle(puzzle []byte) map[string]string {
	if len(puzzle) != 81 {
		log.Fatal("puzzle too small or too big")
	}
	grid := map[string]string{}
	sqrs := files.GetSquares()
	for i := 0; i < 81; i++ {
		val := string(puzzle[i])
		if val == "-" {
			val = files.Digits
		}
		grid[sqrs[i]] = val
	}
	return grid
}
