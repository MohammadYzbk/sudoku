package files

import (
	"fmt"
	"strings"
)

const (
	Digits = "123456789"
	rows   = "ABCDEFGHI"
	cols   = Digits
)

var squares = cross(rows, cols)
var allUnits, units = generateAllUnits()
var peers = generatePeers()

func GetSquares() []string {
	return cross(rows, cols)
}

// squares   = cross(rows, cols)
// Grid      = Dict[Square, DigitSet] # E.g. {'A9': '123', ...}
// all_boxes = [cross(rs, cs)  for rs in ('ABC','DEF','GHI') for cs in ('123','456','789')]
// all_units = [cross(rows, c) for c in cols] + [cross(r, cols) for r in rows] + all_boxes
// units     = {s: tuple(u for u in all_units if s in u) for s in squares}
// peers     = {s: set().union(*units[s]) - {s} for s in squares}

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

// func generateAllBoxes() [][]string {
// 	product := [][]string{}
// 	rs := []string{"ABC", "DEF", "GHI"}
// 	cs := []string{"123", "456", "789"}
// 	for _, row := range rs {
// 		for _, col := range cs {
// 			product = append(product, cross(row, col))
// 		}
// 	}
// 	return product
// }

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

func isSolution(solution, puzzle map[string]string) bool {
	containsAll := true
	for position, val := range solution {
		if !contains(puzzle[position], val) {
			containsAll = false
			break
		}
	}

	isValid := true
	for _, unit := range allUnits {
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

// def constrain(grid) -> Grid:
//     "Propagate constraints on a copy of grid to yield a new constrained Grid."
//     result: Grid = {s: digits for s in squares}
//     for s in grid:
//         if len(grid[s]) == 1:
//             fill(result, s,  grid[s])
//     return result

func Constrain(grid map[string]string) map[string]string {
	result := defaultDigitSet()
	for sqr := range grid {
		if len(grid[sqr]) == 1 {
			fill(result, sqr, grid[sqr])
		}
	}
	return result
}

// def fill(grid: Grid, s: Square, d: Digit) -> Optional[Grid]:
//     """Eliminate all the digits except d from grid[s]."""
//     if grid[s] == d or all(eliminate(grid, s, d2) for d2 in grid[s] if d2 != d):
//         return grid
//     else:
//         return None

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

// def eliminate(grid: Grid, s: Square, d: Digit) -> Optional[Grid]:
//     """Eliminate d from grid[s]; implement the two constraint propagation strategies."""
//     if d not in grid[s]:
//         return grid        ## Already eliminated
//     grid[s] = grid[s].replace(d, '')
//     if not grid[s]:
//         return None        ## None: no legal digit left
//     elif len(grid[s]) == 1:
//         # 1. If a square has only one possible digit, then eliminate that digit as a possibility for each of the square's peers.
//         d2 = grid[s]
//         if not all(eliminate(grid, s2, d2) for s2 in peers[s]):
//             return None    ## None: can't eliminate d2 from some square
//     for u in units[s]:
//         dplaces = [s for s in u if d in grid[s]]
//         # 2. If a unit has only one possible square that can hold a digit, then fill the square with the digit.
//         if not dplaces or (len(dplaces) == 1 and not fill(grid, dplaces[0], d)):
//             return None    ## None: no place in u for d
//     return grid

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

// def search(grid) -> Grid:
//     "Depth-first search with constraint propagation to find a solution."
//     if grid is None:
//         return None
//     s = min((s for s in squares if len(grid[s]) > 1),
//             default=None, key=lambda s: len(grid[s]))
//     if s is None: # No squares with multiple possibilities; the search has succeeded
//         return grid
//     for d in grid[s]:
//         solution = search(fill(grid.copy(), s, d))
//         if solution:
//             return solution
//     return None

func Search(grid map[string]string) map[string]string {
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
		solution := Search(fill(newGrid, minS, string(digit)))
		if len(solution) != 0 {
			return solution
		}
	}
	return map[string]string{}
}
