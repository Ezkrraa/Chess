package main

import "fmt"

type color string

const (
	colorRed  = "\033[0;31m"
	colorNone = "\033[0m"
)

type gameState struct {
	pieces   [8][8]Field
	moves    []AbsoluteMove
	selected *Coordinate
}

func getDefaultBoard() *gameState {
	return &gameState{
		[8][8]Field{
			{{Rook, true}, {Bishop, true}, {Knight, true}, {Queen, true}, {King, true}, {Knight, true}, {Bishop, true}, {Rook, true}},
			{{Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}, {Pawn, true}},
			{},
			{},
			{},
			{},
			{{Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}, {Pawn, false}},
			{{Rook, false}, {Bishop, false}, {Knight, false}, {Queen, false}, {King, false}, {Knight, false}, {Bishop, false}, {Rook, false}},
		},
		[]AbsoluteMove{},
		nil,
	}
}

func printBoard(state *gameState) {
	for i := int8(0); i < 8; i++ {
		for j := int8(0); j < 8; j++ {
			if state.selected != nil && state.selected.y == i && state.selected.x == j {
				fmt.Print(colorRed)
			}
			if state.pieces[i][j].isWhite {
				switch state.pieces[i][j].pieceType {
				case None:
					if (i+j)%2 == 0 {
						fmt.Print("■")
					} else {
						fmt.Print("□")
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
				switch state.pieces[i][j].pieceType {
				case None:
					if (i+j)%2 == 0 {
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
		}
		fmt.Print("\n")
	}
}

// if not mate, relative difference of white material and black material (between -1-1)
func evaluateState(state *gameState) float32 {
	whitePoints := float32(0)
	blackPoints := float32(0)
	for y := range state.pieces {
		for x := range state.pieces[y] {
			switch state.pieces[y][x].pieceType {
			case None:
				continue
			case King:
				if isInCheck(state, Coordinate{int8(y), int8(x)}) {
					return -1
				}
			default:
				if state.pieces[y][x].isWhite {
					whitePoints += getValue(state.pieces[y][x].pieceType)
				} else {
					blackPoints += getValue(state.pieces[y][x].pieceType)
				}
			}
		}
	}
	return clamp((whitePoints-blackPoints)*0.5, -.99, .99)
}

func move(state *gameState, move AbsoluteMove) bool {
	if !isAllowedMove(state, move, false) {
		// this shouldn't happen so I throw before the problem gets worse
		panic("Move was not allowed")
	}
	state.pieces[move.end.y][move.end.x] = state.pieces[move.start.y][move.start.x]
	state.pieces[move.start.y][move.start.x] = Field{None, true}
	return true
}

// checks if the direct line between two points is clear
func isBlockedMove(state *gameState, move AbsoluteMove) bool {
	pointsBetween := getPointsAlongMove(move)
	for i := range *pointsBetween {
		point := (*pointsBetween)[i]
		if state.pieces[point.y][point.x].pieceType != None {
			return true
		}
	}
	return false
}

// checks:
// if there is a piece to move
// whether the piece can move like that (ie. rooks can't go diagonally)
// whether the destination isn't filled by a friendly piece or king
// whether a move is within bounds (coordinates within the 8x8 board)
// whether this move results in a check
func isAllowedMove(state *gameState, move AbsoluteMove, allowKingHits bool) bool {
	start := state.pieces[move.start.y][move.start.x]
	end := state.pieces[move.end.y][move.end.x]
	if end.pieceType == King && !allowKingHits {
		return false // cannot kill kings
	} else if end.pieceType != None && start.isWhite == end.isWhite {
		return false // cannot kill piece of same color
	} else if start.pieceType == None {
		return false // nonexistant piece cannot move
	}
	if start.pieceType != Knight && isBlockedMove(state, move) {
		return false
	}
	allowedEndCoords := filterOnBoardMoves(getMoves((*state).pieces[move.start.y][move.start.x].pieceType, move.start, (*state).pieces[move.start.y][move.start.x].isWhite))
	isAllowedStep := false
	for i := range allowedEndCoords {
		fmt.Println("Coord: ", allowedEndCoords[i])
		if allowedEndCoords[i] == move.end {
			isAllowedStep = true
		}
	}
	if !isAllowedStep {
		return false
	}
	// if this move puts you in check, NOT allowed
	for i := range int8(8) {
		for j := range int8(8) {
			if state.pieces[j][i].pieceType == King && state.pieces[j][i].isWhite == start.isWhite {
				return !isInCheck(state, Coordinate{y: j, x: i})
			}
		}
	}
	panic("No king found on the board")
}

// gets all the points along the line between two points
// this only works for straight and diagonal lines, which is good enough
func getPointsAlongMove(move AbsoluteMove) (coords *[]Coordinate) {
	coords = &[]Coordinate{}
	sX, sY := int8(1), int8(1)
	if move.start.y > move.end.y {
		sX = int8(-1)
	} else if move.start.y == move.end.y {
		sX = int8(0)
	}
	if move.start.x > move.end.x {
		sY = int8(-1)
	} else if move.start.x == move.end.x {
		sY = int8(0)
	}

	current := move.start
	for {
		if current == move.end {
			return coords
		}
		current.y += sX
		current.x += sY
		*coords = append(*coords, current)
	}
}

func isInCheck(state *gameState, kingPos Coordinate) bool {
	king := state.pieces[kingPos.y][kingPos.x]
	if king.pieceType != King {
		panic("Tried to check for mate with invalid piece")
	}
	moves := filterOnBoardMoves(append(append(getDiagonals(kingPos), getHorseSteps(kingPos)...), getLines(kingPos)...))
	for i := range moves {
		if isAllowedMove(state, AbsoluteMove{moves[i], kingPos}, true) {
			return true
		}
	}
	return false
}

func hasAllowedMoves(state *gameState, isWhite bool) bool {
	pieces := []Coordinate{}
	for i := range int8(8) {
		for j := range int8(8) {
			if state.pieces[j][i].pieceType != None && (state.pieces[j][i]).isWhite == isWhite {
				pieces = append(pieces, Coordinate{i, j})
			}
		}
	}
	for i := range pieces {
		pieceMoves := getMoves(state.pieces[pieces[i].y][pieces[i].x].pieceType, pieces[i], isWhite)
		for j := range pieceMoves {
			if isAllowedMove(state, AbsoluteMove{pieces[i], pieceMoves[j]}, false) {
				return true
			}
		}
	}
	return false
}
