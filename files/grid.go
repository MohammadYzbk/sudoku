package files

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

const (
	puzzleLength = 81
	Digits       = "123456789"
	rows         = "ABCDEFGHI"
	cols         = Digits
)

type puzzle struct {
	Peers map[string][]string // For each square, provide all the squares at its row, column, and box.
	// Units are a row, a column, or box. Anything that must have values strictly 1-9.
	allUnits [][]string            // All units in the puzzle.
	Units    map[string][][]string // For each square, provide a list of every unit it belongs to.
	// Units can be a row, a column, or box. Anything that must have values strictly 1-9.
	Squares []string           // List of all the squares in the puzzle.
	Grid    *map[string]string // puzzle itself
}

func NewPuzzle(p map[string]string) puzzle {
	var grid = defaultDigitSet()
	for key, val := range p {
		grid[key] = val
	}

	return puzzle{
		Squares:  squares,
		Units:    units,
		allUnits: allUnits,
		Peers:    peers,
		Grid:     &grid,
	}
}

func GeneratePuzzle() (Grid, Grid) {
	puzzleSeed := map[string]string{}
	pool := make([]string, len(squares))
	copy(pool, squares)
	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	i := 0
	for ; i < 9; i++ {
		digit := string(Digits[i])
		puzzleSeed[pool[i]] = digit
	}

	generatedPuzzle := NewPuzzle(puzzleSeed)
	*generatedPuzzle.Grid = search(constrain(*generatedPuzzle.Grid))
	for i < puzzleLength {
		square := pool[i]
		newGrid := make(map[string]string, len(*generatedPuzzle.Grid))
		for key, value := range *generatedPuzzle.Grid {
			newGrid[key] = value
		}
		dtd := newGrid[square]
		delete(newGrid, square)
		newGrid = eliminate(newGrid, square, dtd)
		numberOfSolutions := msearchWithLimit(constrain(newGrid), 2)
		i++
		if numberOfSolutions == 1 {
			*generatedPuzzle.Grid = newGrid
		}
	}

	return *generatedPuzzle.Grid, search(constrain(*generatedPuzzle.Grid))
}

func Enumarate(grid Grid) int {
	return msearch(constrain(grid))
}

func (p puzzle) IsSolution(solution map[string]string) bool {
	if p.Grid == nil {
		log.Fatal("grid is nil")
	}
	containsAll := true
	for position, val := range solution {
		if !contains((*p.Grid)[position], val) {
			containsAll = false
			break
		}
	}

	isValid := true
	for _, unit := range p.allUnits {
		if !isValid {
			break
		}
		solutionForUnit := ""
		for _, sqr := range unit {
			solutionForUnit += solution[sqr]
		}
		for _, num := range Digits {
			str := string(num)
			if !contains(solutionForUnit, str) {
				isValid = false
				break
			}
		}
	}

	return len(solution) != 0 && containsAll && isValid
}

type Grid map[string]string

func (s Grid) ToString() string {
	return flattenPuzzle(s)
}

func (s Grid) ToVerboseString() string {
	return flattenCandidatePuzzle(constrain(s))
}

func (p puzzle) Solve() Grid {
	if p.Grid == nil {
		log.Fatal("Grid is nil")
	}
	return search(constrain(*p.Grid))
}

func (p puzzle) ToString() string {
	if p.Grid == nil {
		log.Fatal("Grid is nil")
	}
	return flattenPuzzle(*p.Grid)
}

func flattenPuzzle(grid map[string]string) string {
	str := ""
	for _, sqr := range GetSquares() {
		if len(grid[sqr]) != 1 {
			str += "-"
			continue
		}
		str += grid[sqr]
	}
	return str
}

func flattenCandidatePuzzle(grid map[string]string) string {
	str := ""
	for _, sqr := range GetSquares() {
		if len(grid[sqr]) != 1 {
			str += fmt.Sprintf("{%s}", grid[sqr])
			continue
		}
		str += grid[sqr]
	}
	return str
}

var squares = cross(rows, cols)
var allUnits, units = generateAllUnits()
var peers = generatePeers()

func GetSquares() []string {
	return cross(rows, cols)
}

func defaultDigitSet() map[string]string {
	digitSet := map[string]string{}
	for _, sqr := range squares {
		digitSet[sqr] = Digits
	}
	return digitSet
}

func cross(A string, B string) []string {
	product := []string{}
	for _, i := range A {
		for _, j := range B {
			product = append(product, fmt.Sprintf("%s%s", string(i), string(j)))
		}
	}
	return product
}

func generateAllUnits() ([][]string, map[string][][]string) {
	product := [][]string{}
	units := map[string][][]string{}
	rs := []string{"ABC", "DEF", "GHI"}
	cs := []string{"123", "456", "789"}
	for _, row := range rs {
		for _, col := range cs {
			crossProduct := cross(row, col)
			product = append(product, cross(row, col))
			for _, sqr := range crossProduct {
				units[sqr] = append(units[sqr], crossProduct)
			}
		}
	}

	for _, row := range rows {
		crossProduct := cross(string(row), cols)
		product = append(product, crossProduct)
		for _, sqr := range crossProduct {
			units[sqr] = append(units[sqr], crossProduct)
		}
	}

	for _, col := range cols {
		crossProduct := cross(rows, string(col))
		product = append(product, crossProduct)
		for _, sqr := range crossProduct {
			units[sqr] = append(units[sqr], crossProduct)
		}
	}

	return product, units
}

func generatePeers() map[string][]string {
	product := map[string][]string{}
	for _, s := range squares {
		memo := map[string]int{}
		memo[s] = 1
		for _, unit := range units[s] {
			for _, s2 := range unit {
				if memo[s2] != 1 {
					product[s] = append(product[s], s2)
					memo[s2] = 1
				}
			}
		}
	}
	return product
}

func contains(s string, t string) bool {
	if len(t) > 1 {
		return false
	}
	for _, single := range s {
		if string(single) == t {
			return true
		}
	}

	return false
}

func constrain(grid map[string]string) map[string]string {
	result := defaultDigitSet()
	for sqr := range grid {
		if len(grid[sqr]) == 1 {
			fill(result, sqr, grid[sqr])
		}
	}
	return result
}

func fill(grid map[string]string, sqr string, digit string) map[string]string {
	if grid[sqr] == digit {
		return grid
	}

	for _, cand := range grid[sqr] {
		if string(cand) != digit {
			if newGrid := eliminate(grid, sqr, string(cand)); len(newGrid) == 0 {
				return map[string]string{}
			}
		}
	}

	return grid
}

func eliminate(grid map[string]string, sqr string, digitToDie string) map[string]string {
	if !contains(grid[sqr], digitToDie) {
		return grid
	}
	grid[sqr] = strings.Replace(grid[sqr], digitToDie, "", 1)
	if grid[sqr] == "" {
		return map[string]string{}
	} else if len(grid[sqr]) == 1 {
		chosenDigit := grid[sqr]
		for _, peer := range peers[sqr] {
			if newGrid := eliminate(grid, peer, chosenDigit); len(newGrid) == 0 {
				return map[string]string{}
			}
		}
	}
	for _, u := range units[sqr] {
		dplaces := []string{}
		for _, peer := range u {
			if contains(grid[peer], digitToDie) {
				dplaces = append(dplaces, peer)
			}
		}
		if len(dplaces) == 0 {
			return map[string]string{}
		}

		if len(dplaces) == 1 {
			res := fill(grid, dplaces[0], digitToDie)
			if len(res) == 0 {
				return map[string]string{}
			}
		}
	}
	return grid
}

func search(grid map[string]string) map[string]string {
	if len(grid) == 0 {
		return map[string]string{}
	}
	minS := ""
	minVal := 10
	for _, sqr := range squares {
		if len(grid[sqr]) > 1 {
			if minVal > len(grid[sqr]) {
				minVal = len(grid[sqr])
				minS = sqr
			}
		}
	}

	if minS == "" {
		return grid
	}

	for _, digit := range grid[minS] {
		newGrid := make(map[string]string, len(grid))
		for key, value := range grid {
			newGrid[key] = value
		}
		solution := search(fill(newGrid, minS, string(digit)))
		if len(solution) != 0 {
			return solution
		}
	}
	return map[string]string{}
}

func msearch(grid map[string]string) int {
	if len(grid) == 0 {
		return 0
	}
	minS := ""
	minVal := 10
	for _, sqr := range squares {
		if len(grid[sqr]) > 1 {
			if minVal > len(grid[sqr]) {
				minVal = len(grid[sqr])
				minS = sqr
			}
		}
	}

	if minS == "" {
		return 1
	}

	counter := 0
	for _, digit := range grid[minS] {
		newGrid := make(map[string]string, len(grid))
		for key, value := range grid {
			newGrid[key] = value
		}
		solution := msearch(fill(newGrid, minS, string(digit)))
		counter += solution
	}

	return counter
}

func msearchWithLimit(grid map[string]string, limit int) int {
	if len(grid) == 0 {
		return 0
	}
	minS := ""
	minVal := 10
	for _, sqr := range squares {
		if len(grid[sqr]) > 1 {
			if minVal > len(grid[sqr]) {
				minVal = len(grid[sqr])
				minS = sqr
			}
		}
	}

	if minS == "" {
		return 1
	}

	counter := 0
	for _, digit := range grid[minS] {
		newGrid := make(map[string]string, len(grid))
		for key, value := range grid {
			newGrid[key] = value
		}
		solution := msearchWithLimit(fill(newGrid, minS, string(digit)), limit)
		counter += solution
		if counter >= limit {
			return counter
		}
	}

	return counter
}
