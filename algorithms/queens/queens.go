package queens

//   0 1 2
// 0     x
// 1   x
// 2 x

// backtracking problem
// loop through the cols each time the rows
// will happen on each dfs
// keep a tally of queens that are placed in
// rows, columns, asc diag and desc diag
// if placing a queen on the next spot would
// add another queen to the one already in
// a specifc index of the above tallies then
// we can't place it
// asc calc row+col
// desc row-col
// func solveNQueens(n int) [][]string {
// 	rv := make([][]string, 0, n)
// 	curr := make([][]string, 0, n)
//
// 	rows := map[int]bool{}
// 	cols := map[int]bool{}
// 	diagAsc := map[int]bool{}
// 	diagDesc := map[int]bool{}
//
// 	var dfs func(row int)
// 	dfs = func(row int) {
// 		if row == n {
// 			tmp := make([]string, 0, n)
// 			for _, v := range curr {
// 				c := append([]string{}, v...)
// 				tmp = append(tmp, strings.Join(c, ""))
// 			}
// 			rv = append(rv, tmp)
// 			return
// 		}
//
// 		// make new row
// 		r := make([]string, 0, n)
// 		for i := range n {
// 			r[i] = "."
// 		}
//
// 		for col := range n {
//
// 		}
// 	}
//
// 	for i := range n {
// 		dfs(i)
// 	}
//
// 	return rv
// }
