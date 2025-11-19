# Sudoku Map Generator & Solver
This repository is a throwaway implementation of a solver and generator of Sudoku puzzles.
This implementation is based on the work of Peter Norvig and uses his representation of Sudoku puzzle

## Puzzle Representation
Each puzzle represented by a 9x9 matrix with rows represented through letters (ABCDEFGHI) and columns represented through numbers (123456789).
A square on the puzzle is indexed by a concatenation of its row and its column (eg. A1, A2, B1, B2 ...).

The value of each square will be the possible candidates or digits (123456789) that can take that square.

### Technical Representation
In code, the puzzle above is translated into a hash map where both keys and values are strings (eg. "A1":"123", B2:"1" ...). 
```go
grid := map[string]string{
	"A1":"12345",
	"A2":"1",
	"A3":"23",
}
```

This representation is useful for detecting wrong paths while solving, as a square in the map with an empty string `""` implies the puzzle has become unsolvable and as such we must backtrack, more on this when we discuss solutions.

## Solution Detection

### Units
Units are the lists of squares from every row, column and box. All units share the same property of needing to be filled by every value from 1 to 9.
Examples of units:
- {A1, A2, A3, A4, A5, A6, A7, A8, A9}
- {A1, B1, C1, D1, E1, F1, G1, H1, I1}
- {A1, A2, A3, B1, B2, B3, C1, C2, C3}
- ...

The usefulness of units is that if every unit of the puzzle is full and valid. Then the puzzle is solved.

## Solving:

### Peers
The peers of some square A1, are defined as squares that lie on the same row, column, or box as A1, as such, those squares constrain the possible values of A1.

This solution uses constraint satisfaction to eliminate candidates for squares and a DFS algo to choose different candidates for squares with multiple candidates.
Solving a Sudoku is then constraining the puzzle to remove as many candidates as possible and filling in squares with 1 possible candidate left, then picking a possible candidate for a square and constraining again.
This process continues until we find a solution, or we explore every possible path to find that the puzzle is unsolvable.

### Constraining
When eliminating every impossible candidate from the puzzle, we are actively iterating over every square in the puzzle to find ones with 1 possible candidates. 
We will then fill that square with said candidate then iterate over the peers of that square and eliminate that square in each of their candidate list. 
This process repeats until we run out of squares with 1 possible candidate and sometimes that could mean that we already solved puzzle but that's usually restricted to easy puzzle (check out the `constrainOnly.txt` puzzle for an example of one).

```go
func constrain(grid map[string]string) map[string]string {
	// create new puzzle with each square filled with every digit as a candidate
	result := defaultDigitSet()
	// for every square with 1 candidate in the original puzzle, fill that square with said value in the new puzzle  
	for sqr := range grid {
		if len(grid[sqr]) == 1 {
			fill(result, sqr, grid[sqr])
		}
	}
	return result
}
```

```go
func fill(grid map[string]string, sqr string, digit string) map[string]string {
	// if the only candidate left is that exact value, we return instantly.
	if grid[sqr] == digit {
		return grid
	}

	// we eliminate every other digit from the candidates from the square
	for _, cand := range grid[sqr] {
		if string(cand) != digit {
			if newGrid := eliminate(grid, sqr, string(cand)); len(newGrid) == 0 {
				// if the elimination resulted in an unsolvable puzzle, we return an empty grid as a sign of failure
				return map[string]string{}
			}
		}
	}

	return grid
}
```

```go
func eliminate(grid map[string]string, sqr string, digitToDie string) map[string]string {
	// if the digit is already eliminated, return success
	if !contains(grid[sqr], digitToDie) {
		return grid
	}
	// remove digit from candidates
	grid[sqr] = strings.Replace(grid[sqr], digitToDie, "", 1)
	// return failure if we are out of candidates
	if grid[sqr] == "" {
		return map[string]string{} 
	// if only one candidate remains, remove that digit from the candidate list of the peers
	} else if len(grid[sqr]) == 1 {
		chosenDigit := grid[sqr]
		for _, peer := range peers[sqr] {
			if newGrid := eliminate(grid, peer, chosenDigit); len(newGrid) == 0 {
				return map[string]string{}
			}
		}
	}
	// checking if puzzle is still solvable
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
```

WIP