package main

import "fmt"

// 0 1 1 0 1
// 1 1 0 1 0
// 0 1 1 1 0
// 1 1 1 1 0
// 1 1 1 1 1
// 0 0 0 0 1

func main() {
	field := [][]int{
		{0, 1, 1, 0, 1},
		{1, 1, 0, 1, 0},
		{0, 1, 1, 1, 0},
		{1, 1, 1, 1, 0},
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 1},
	}

	maxSize := 0
	pos := make([]int, 2)

	// horizontally
	for i, row := range field {
		for j, cell := range row {

			// if the i or j are 0, that means we can't process the values on i-1=0-1=-1
			// in this case we will continue for the next iteration.
			if i == 0 || j == 0 {
				continue
			}

			// if the cell value is '0' then continue for the next cell.
			if cell == 0 {
				continue
			}

			// if the cell is not '0', then check for the left, right and diag values.
			top := field[i-1][j]
			left := field[i][j-1]
			diag := field[i-1][j-1]

			if top > 0 && left > 0 && diag > 0 {
				min := min(top, left, diag)
				field[i][j] = min + 1
				curr := field[i][j]
				if curr > maxSize {
					maxSize = curr
					pos[0] = i
					pos[1] = j
				}
			}
		}
	}

	fmt.Printf("Max Size : %d | i: %d | j: %d\n", maxSize, pos[0], pos[1])
	fmt.Printf("Area from location: (%d, %d) of Size: %d", pos[0]-maxSize+1, pos[1]-maxSize+1, maxSize)
}
