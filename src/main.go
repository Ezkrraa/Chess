package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nsf/termbox-go"
)

func main() {
	termbox.Init() // should be done 'before any other functions are run'
	defer termbox.Close()

	if len(os.Args) > 1 && os.Args[1] == "test" {
		runTests()
		return
	}
	// TODO: build main menu with either local multiplayer or CPU (TODO: build CPU first)
	gameLoop()
}

func gameLoop() {
	state := getDefaultBoard()
	currentPlayer := true // false == white, true == black
	for {
		if !state.hasAllowedMoves(currentPlayer) {
			break
		}
		startMoveCoordinate := selectToMove(state, currentPlayer)
		if startMoveCoordinate.x == -1 {
			fmt.Println("Stopping...")
			return
		}
		madeMove := selectDestination(state, currentPlayer, startMoveCoordinate)
		if !madeMove {
			fmt.Println("Stopping...")
			return
		}
		currentPlayer = !currentPlayer
	}
	fmt.Println(" won!")
}

func selectDestination(state *gameState, currentPlayer bool, startCoord Coordinate) bool {
	selectedCoord := startCoord
	for {
		clearTerminal()
		if len(state.moves)%2 == 0 {
			fmt.Println("White's turn")
		} else {
			fmt.Println("Black's turn")
		}
		printBoard(state, []Coordinate{selectedCoord}, filterOnBoardMoves(getMoves(state.gameboard[startCoord.y][startCoord.x].pieceType, startCoord, currentPlayer)))
		a := termbox.PollEvent().Key
		switch a {
		case termbox.KeyEsc:
			return false
		case termbox.KeyArrowUp:
			selectedCoord.y += 1
		case termbox.KeyArrowDown:
			selectedCoord.y -= 1
		case termbox.KeyArrowRight:
			selectedCoord.x += 1
		case termbox.KeyArrowLeft:
			selectedCoord.x -= 1
		case termbox.KeyEnter:
			if state.attemptMove(AbsoluteMove{startCoord, selectedCoord}, currentPlayer) {
				return true
			}
		}
		selectedCoord = clampCoord(selectedCoord)
	}
}

// return a position of a piece to move
func selectToMove(state *gameState, currentPlayer bool) Coordinate {
	selectedCoord := Coordinate{0, 0}
	for {
		clearTerminal()
		if len(state.moves)%2 == 0 {
			fmt.Println("White's turn")
		} else {
			fmt.Println("Black's turn")
		}
		// fmt.Printf("Selected: \n%s\n", state.showInfo(selectedCoord))
		// fmt.Printf("IsMine: %t\n", state.isMine(selectedCoord, currentPlayer))
		printBoard(state, []Coordinate{selectedCoord}, []Coordinate{})
		a := termbox.PollEvent().Key
		switch a {
		case termbox.KeyEsc:
			return Coordinate{-1, -1}
		case termbox.KeyArrowUp:
			selectedCoord.y += 1
		case termbox.KeyArrowDown:
			selectedCoord.y -= 1
		case termbox.KeyArrowRight:
			selectedCoord.x += 1
		case termbox.KeyArrowLeft:
			selectedCoord.x -= 1
		case termbox.KeyEnter:
			if state.isMine(selectedCoord, currentPlayer) {
				return selectedCoord
			} else {
				state.showInfo(selectedCoord)
				// panic(state.gameboard[selectedCoord.y][selectedCoord.x].isWhite)
			}
		}
		selectedCoord = clampCoord(selectedCoord)
	}
}

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// TODO: make main loop that takes input and figures out who won
