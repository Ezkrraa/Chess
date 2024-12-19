package main

import (
	"fmt"
	"math"
	"slices"
)

type color string

const (
	colorRed      = "\033[0;7m"
	colorYellow   = "\033[0;33m"
	colorNormal   = "\033[0m"
	colorSelected = "\033[0;5;31m"
)

type gameState struct {
	gameboard [8][8]Field
	moves     []AbsoluteMove
}

func getDefaultBoard() *gameState {
	return &gameState{
		[8][8]Field{
			{{Rook, true}, {Knight, true}, {Bishop, true}, {Queen, true}, {King, true}, {Bishop, true}, {Knight, true}, {Rook, true}},
			{{Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}},
			{},
			{},
			{},
			{},
			{{Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}},
			{{Rook, false}, {Knight, false}, {Bishop, false}, {Queen, false}, {King, false}, {Bishop, false}, {Knight, false}, {Rook, false}},
		},
		[]AbsoluteMove{},
	}
}

func printBoard(state *gameState, selected []Coordinate, highlighted []Coordinate) {
	for y := int8(7); y > -1; y-- { // print highest indeces first since we draw top-down
		for x := int8(0); x < 8; x++ {
			if slices.ContainsFunc(selected, func(coord Coordinate) bool { return coord == Coordinate{y, x} }) { // color selected fields red
				fmt.Print(colorSelected)
			} else if slices.ContainsFunc(highlighted, func(coord Coordinate) bool { return coord == Coordinate{y, x} }) { // color highlighted fields yellow
				fmt.Print(colorYellow)
			} else {
				fmt.Print(colorNormal)
			}
			if state.gameboard[y][x].isWhite {
				switch state.gameboard[y][x].pieceType {
				case None: // draw checkerboard
					if (y+x)%2 == 0 {
						fmt.Print("▒")
					} else {
						fmt.Print("█")
					}
				case Pawn:
					fmt.Print("♙")
				case Knight:
					fmt.Print("♘")
				case Bishop:
					fmt.Print("♗")
				case Rook:
					fmt.Print("♖")
				case Queen:
					fmt.Print("♕")
				case King:
					fmt.Print("♔")
				}
			} else {
				switch state.gameboard[y][x].pieceType {
				case None:
					if (y+x)%2 == 0 {
						fmt.Print("▒")
					} else {
						fmt.Print("█")
					}
				case Pawn:
					fmt.Print("♟")
				case Knight:
					fmt.Print("♞")
				case Bishop:
					fmt.Print("♝")
				case Rook:
					fmt.Print("♜")
				case Queen:
					fmt.Print("♛")
				case King:
					fmt.Print("♚")
				}
			}
			fmt.Print(colorNormal + " ")
		}
		fmt.Print("\n")
	}
}

func (state *gameState) copyWithMove(move AbsoluteMove) (newState *gameState) {
	newState = getDefaultBoard()
	newState.moves = state.moves
	newState.gameboard = state.gameboard
	newState.gameboard[move.end.y][move.end.x] = newState.gameboard[move.start.y][move.start.x]
	newState.gameboard[move.start.y][move.start.x] = Field{None, true}
	newState.moves = append(newState.moves, move)
	return newState
}

// if not mate, relative difference of white material and black material (between -1-1)
func (state *gameState) evaluateState() float32 {
	whitePoints := float32(0)
	blackPoints := float32(0)
	for y := range state.gameboard {
		for x := range state.gameboard[y] {
			switch state.gameboard[y][x].pieceType {
			case None:
				continue
			case King:
				if state.isInCheck(Coordinate{int8(y), int8(x)}) {
					return -1
				}
			default:
				if state.gameboard[y][x].isWhite {
					whitePoints += getValue(state.gameboard[y][x].pieceType)
				} else {
					blackPoints += getValue(state.gameboard[y][x].pieceType)
				}
			}
		}
	}
	return clamp((whitePoints-blackPoints)*0.5, -.99, .99)
}

func (state *gameState) attemptMove(move AbsoluteMove, player bool) bool {
	if !state.isLegalMove(move, player) {
		return false
	}
	state.gameboard[move.end.y][move.end.x] = state.gameboard[move.start.y][move.start.x]
	state.gameboard[move.start.y][move.start.x] = Field{None, true}
	state.moves = append(state.moves, move)
	return true
}

// checks if the direct line between two points is clear
func (state *gameState) isBlockedMove(move AbsoluteMove) bool {
	pointsBetween := getPointsAlongMove(move)
	for i := range *pointsBetween {
		point := (*pointsBetween)[i]
		if state.gameboard[point.y][point.x].pieceType != None {
			return true
		}
	}
	return false
}

func (state *gameState) isMine(coord Coordinate, color bool) bool {
	return state.gameboard[coord.y][coord.x].isWhite == color
}

// checks:
// if there is a piece to move
// whether the piece can move like that (ie. rooks can't go diagonally)
// whether the destination isn't filled by a friendly piece or king
// whether a move is within bounds (coordinates within the 8x8 board)
// whether this move results in a check
func (state *gameState) isLegalMove(move AbsoluteMove, allowKingHits bool) bool {
	start := state.gameboard[move.start.y][move.start.x]
	end := state.gameboard[move.end.y][move.end.x]
	if end.pieceType == King && !allowKingHits {
		return false // cannot kill kings
	} else if end.pieceType != None && start.isWhite == end.isWhite {
		return false // cannot kill piece of same color
	} else if start.pieceType == None {
		return false // nonexistant piece cannot move
	}
	if start.pieceType != Pawn && start.pieceType != Knight && state.isBlockedMove(move) { // skip for pawns since they get a special case check
		return false
	}
	if start.pieceType != Pawn {
		allowedEndCoords := filterOnBoardMoves(getMoves((*state).gameboard[move.start.y][move.start.x].pieceType, move.start, (*state).gameboard[move.start.y][move.start.x].isWhite))
		isAllowedStep := false
		for i := range allowedEndCoords {
			if allowedEndCoords[i] == move.end {
				isAllowedStep = true
			}
		}
		if !isAllowedStep {
			return false
		}
	} else { // handle pawn movement
		if !isValidPawnMove(state, move) {
			return false
		}
	}
	// straight ahead moves, x=0,y=1|2
	// diagonal moves, absx=1,absy=1
	// en passant moves, absx=1,y=1
	// if absXDiff == 0 && // if absolute x diff is null AND absolute y diff
	// if this move puts you in check, NOT allowed
	for i := range int8(8) {
		for j := range int8(8) {
			if state.gameboard[j][i].pieceType == King && state.gameboard[j][i].isWhite == start.isWhite {
				afterMoveState := state.copyWithMove(move)
				// panic("6")
				return !afterMoveState.isInCheck(Coordinate{y: j, x: i})
			}
		}
	}
	panic("No king found on the board")
}

func isValidPawnMove(state *gameState, move AbsoluteMove) bool {
	yDiff := move.end.y - move.start.y
	absXDiff := math.Abs(float64(move.end.x - move.start.x))
	absYDiff := math.Abs(float64(move.end.y - move.start.y))
	startPiece := state.gameboard[move.start.y][move.start.x]
	endPiece := state.gameboard[move.end.y][move.end.x]
	if (startPiece.isWhite && yDiff < 0) || (!startPiece.isWhite && yDiff > 0) {
		// if moving backwards, return false
		return false
	} else if absYDiff == 1 && absXDiff == 1 {
		// handle diagonal hits
		return endPiece.pieceType != None && endPiece.isWhite != startPiece.isWhite
	} else if absXDiff == 1 && absYDiff == 1 && endPiece.pieceType == None {
		// handle forward moves
		return true
	} else if absXDiff == 0 && absYDiff == 1 && endPiece.pieceType == None {
		return true
	} else if absXDiff == 0 && absYDiff == 2 && endPiece.pieceType == None {
		// handle double step forward moves
		midCoordY := (move.start.y + move.end.y) / 2
		// check if not hopping over something
		return state.gameboard[midCoordY][move.start.x].pieceType == None && endPiece.pieceType == None &&
			// if trying to do a beginning jump when already moved, return false
			((startPiece.isWhite && move.start.y == 1) || (!startPiece.isWhite && move.start.y == 6))
	} else {
		return false
	}
}

// gets all the points along the line between two points
// this only works for straight and diagonal lines, which is good enough
func getPointsAlongMove(move AbsoluteMove) (coords *[]Coordinate) {
	coords = &[]Coordinate{}
	sY, sX := int8(1), int8(1)
	if move.start.y > move.end.y {
		sY = int8(-1)
	} else if move.start.y == move.end.y {
		sY = int8(0)
	}
	if move.start.x > move.end.x {
		sX = int8(-1)
	} else if move.start.x == move.end.x {
		sX = int8(0)
	}

	current := move.start
	for {
		current.y += sY
		current.x += sX
		if current == move.end {
			return coords
		}
		*coords = append(*coords, current)
	}
}

func (state *gameState) isInCheck(kingPos Coordinate) bool {
	king := state.gameboard[kingPos.y][kingPos.x]
	if king.pieceType != King {
		panic("Tried to check for mate with invalid piece")
	}
	moves := filterOnBoardMoves(append(append(getDiagonals(kingPos), getHorseSteps(kingPos)...), getLines(kingPos)...))
	for i := range moves {
		if state.isLegalMove(AbsoluteMove{moves[i], kingPos}, true) {
			return true
		}
	}
	return false
}

func (state *gameState) hasAllowedMoves(isWhite bool) bool {
	pieces := []Coordinate{}
	for i := range int8(8) {
		for j := range int8(8) {
			if state.gameboard[j][i].pieceType != None && (state.gameboard[j][i]).isWhite == isWhite {
				pieces = append(pieces, Coordinate{i, j})
			}
		}
	}
	for i := range pieces {
		pieceMoves := filterOnBoardMoves(getMoves(state.gameboard[pieces[i].y][pieces[i].x].pieceType, pieces[i], isWhite))
		for j := range pieceMoves {
			if state.isLegalMove(AbsoluteMove{pieces[i], pieceMoves[j]}, false) {
				return true
			}
		}
	}
	return false
}
