package main

import (
	"strconv"
)
// Join accepts a join query of two relations, and returns the sum of
// relation0.col0 in the final result.
// Input arguments:
//   f0: file name of the given relation0
//   f1: file name of the given relation1
//   offset0: offsets of which columns the given relation0 should be joined
//   offset1: offsets of which columns the given relation1 should be joined
// Output arguments:
//   sum: sum of relation0.col0 in the final result
func Join(f0, f1 string, offset0, offset1 []int) (sum uint64) {
	return sum
}

func JoinLoop(f0, f1 string, offset0, offset1 []int) (sum uint64) {
	offLen := len(offset0)
	tbl0, tbl1 := readCSVFileIntoTbl(f0), readCSVFileIntoTbl(f1)
	for _, row0 := range tbl0 {
		for _, row1 := range tbl1 {
			cnt := 0
			for id, off := range offset0 {
				if row0[off] == row1[offset1[id]] {
					cnt ++
				}
			}
			if cnt == offLen {
				v, err := strconv.ParseUint(row0[0], 10, 64)
				if err != nil {
					panic("Join panic\n" + err.Error())
				}
				sum += v
			}
		}
	}
	return sum
}
