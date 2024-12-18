package main

import "fmt"

func main() {
	board := getDefaultBoard()
	printBoard(board)
	fmt.Println(evaluateState(board))
	result := move(board, AbsoluteMove{Coordinate{1, 2}, Coordinate{3, 2}})
	fmt.Println(evaluateState(board))
	fmt.Println(result)
	printBoard(board)
}

// func gameLoop() {
// 	board := getDefaultBoard()
// 	currentPlayer := false // false == white, true == black
// 	for {
// 		if(!hasAllowedMoves(board, true)) {
// 			break
// 		}
// 	}
// 	if currentPlayer {
// 		winner = "White"

// 	}
// 	fmt.Println(winner + " won!")
// }

// TODO: make main loop that takes input and figures out who won
