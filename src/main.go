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
	clearTerminal()
	fmt.Println("Stopping...")
}

func gameLoop() {
	state := getDefaultBoard()
	currentPlayer := true // false == white, true == black
	for {
		if !state.hasAllowedMoves(currentPlayer) {
			break
		}

		startMoveCoordinate, result1 := selectToMove(state, currentPlayer)
		if result1 == quit {
			return
		}
		for {
			result2 := selectDestination(state, currentPlayer, startMoveCoordinate)
			if result2 == quit {
				return
			} else if result2 == undo {
				startMoveCoordinate, result1 = selectToMove(state, currentPlayer)
				if result1 == quit {
					return
				}
			} else {
				break
			}
		}

		currentPlayer = !currentPlayer
	}
	fmt.Println(" won!")
	termbox.PollEvent()
}

func selectDestination(state *gameState, currentPlayer bool, startCoord Coordinate) uiResult {
	selectedCoord := startCoord
	for {
		clearTerminal()
		state.printInfo()
		state.printBoard([]Coordinate{selectedCoord}, filterOnBoardMoves(getMoves(state.gameboard[startCoord.y][startCoord.x].pieceType, startCoord, currentPlayer)))
		a := termbox.PollEvent().Key
		switch a {
		case termbox.KeyEsc:
			return quit
		case termbox.KeyBackspace:
			return undo
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
				return success
			}
		}
		selectedCoord = clampCoord(selectedCoord)
	}
}

// return a position of a piece to move
func selectToMove(state *gameState, currentPlayer bool) (Coordinate, uiResult) {
	selectedCoord := Coordinate{0, 0}
	for {
		clearTerminal()
		state.printInfo()
		state.printBoard([]Coordinate{selectedCoord}, []Coordinate{})
		a := termbox.PollEvent().Key
		switch a {
		case termbox.KeyEsc:
			return Coordinate{-1, -1}, quit
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
				return selectedCoord, success
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
