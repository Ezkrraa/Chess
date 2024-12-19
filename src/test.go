package main

import "fmt"

func runTests() bool {

	fmt.Println("Testing testCapture...")
	result1 := testCapture()
	if result1 {
		fmt.Println("testCapture succeeded")
	} else {
		fmt.Println("testCapture failed")
	}
	fmt.Println("testing testInvalidMove...")
	result2 := testInvalidMove()
	if result2 {
		fmt.Println("testInvalidMove succeeded")
	} else {
		fmt.Println("testInvalidMove failed")
	}
	fmt.Println("testing testStartDoubleMove...")
	result3 := testStartDoubleMove()
	if result3 {
		fmt.Println("testStartDoubleMove succeeded")
	} else {
		fmt.Println("testStartDoubleMove failed")
	}
	fmt.Printf("%d%% of tests succeeded", getTestPassedPercentage([]bool{result1, result2, result3}))
	return result1 && result2
}

func getTestPassedPercentage(tests []bool) int {
	totalSucces := float32(0)
	for i := range tests {
		totalSucces += boolToFloat(tests[i])
	}
	relativeAmount := totalSucces / float32(len(tests))
	return int(relativeAmount * 100)
}

func boolToFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func testCapture() bool {
	board := getDefaultBoard()
	board.gameboard = [8][8]Field{
		{{Rook, true}, {Bishop, true}, {Knight, true}, {Queen, true}, {King, true}, {Knight, true}, {Bishop, true}, {Rook, true}},
		{{Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}},
		{},
		{{None, true}, {Pawn, true}},
		{{Pawn, false}},
		{},
		{{Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}},
		{{Rook, false}, {Bishop, false}, {Knight, false}, {Queen, false}, {King, false}, {Knight, false}, {Bishop, false}, {Rook, false}},
	}
	move := AbsoluteMove{Coordinate{3, 0}, Coordinate{4, 1}}
	result := board.attemptMove(move, true)
	if result {
		fmt.Println("Move validation went wrong")
		return false
	}
	return true
}

func testInvalidMove() bool {
	board := getDefaultBoard()
	move := AbsoluteMove{Coordinate{1, 0}, Coordinate{2, 7}} // move WAY sideways and a step forwards
	return !board.attemptMove(move, true)
}

// this one still fails
func testStartDoubleMove() bool {
	board := getDefaultBoard()
	move := AbsoluteMove{Coordinate{1, 0}, Coordinate{3, 0}} // double move forwards
	return board.attemptMove(move, true)
}
