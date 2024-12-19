package main

import "fmt"

type Field struct {
	pieceType PieceType
	isWhite   bool
}

type PieceType byte

const (
	None PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type MoveType byte

const (
	InvalidMove MoveType = iota
	LinearMove
	DiagonalMove
	HorseMove
)

type AbsoluteMove struct {
	start Coordinate
	end   Coordinate
}

type Coordinate struct {
	y int8
	x int8
}

func (state *gameState) showInfo(coord Coordinate) string {
	return fmt.Sprintf("{\n\tX: %d, \n\tY: %d\n\tType: %d\n\tOwner: %t\n}\n", coord.x, coord.y, state.gameboard[coord.y][coord.x].pieceType, state.gameboard[coord.y][coord.x].isWhite)
}

func filterOnBoardMoves(moves []Coordinate) (ret []Coordinate) {
	for i := range moves {
		if isOnBoard(moves[i]) {
			ret = append(ret, moves[i])
		}
	}
	return ret
}

func getMoves(piece PieceType, start Coordinate, isWhite bool) []Coordinate {
	switch piece {
	case None:
		return []Coordinate{}
	case Knight:
		return getHorseSteps(start)
	case Pawn:
		if isWhite {
			if start.y == 1 {
				return []Coordinate{
					{start.y + 1, start.x},
					{start.y + 2, start.x},
					{start.y + 1, start.x - 1},
					{start.y + 1, start.x + 1},
				}
			} else {
				return []Coordinate{
					{start.y + 1, start.x},
					{start.y + 1, start.x - 1},
					{start.y + 1, start.x + 1},
				}
			}
		} else if !isWhite {
			if start.y == 6 {
				return []Coordinate{
					{start.y - 1, start.x},
					{start.y - 2, start.x},
					{start.y - 1, start.x - 1},
					{start.y - 1, start.x + 1},
				}
			} else {
				return []Coordinate{
					{start.y + 1, start.x},
					{start.y - 1, start.x - 1},
					{start.y - 1, start.x + 1},
				}
			}
		}
	case Bishop:
		return getDiagonals(start)
	case Rook:
		return getLines(start)
	case Queen:
		return append(getDiagonals(start), getLines(start)...)
	case King:
		return []Coordinate{{1, 1}, {1, 0}, {1, -1}, {0, 1}, {0, -1}, {-1, 1}, {-1, 0}, {-1, -1}}
	}
	panic(fmt.Sprintf("Piece was not accounted for in switch case: %d", piece))
}

func getHorseSteps(start Coordinate) []Coordinate {
	return []Coordinate{
		{start.y + 2, start.x + 1},
		{start.y + 2, start.x - 1},
		{start.y - 2, start.x + 1},
		{start.y - 2, start.x - 1},
		{start.y + 1, start.x + 2},
		{start.y - 1, start.x + 2},
		{start.y + 1, start.x - 2},
		{start.y - 1, start.x - 2},
	}
}

func getDiagonals(start Coordinate) []Coordinate {
	steps := []Coordinate{}
	// start at 1 to prevent standing still as option
	for i := int8(1); i < 8; i++ {
		steps = append(steps, Coordinate{start.y + i, start.x + i})
		steps = append(steps, Coordinate{(-1 * i) + start.y, i + start.x})
		steps = append(steps, Coordinate{i + start.y, (-1 * i) + start.x})
		steps = append(steps, Coordinate{(-1 * i) + start.y, (-1 * i) + start.x})
	}
	return steps
}

func getLines(start Coordinate) (steps []Coordinate) {
	// start at 1 to prevent standing still as option
	for i := int8(1); i < 8; i++ {
		steps = append(steps, Coordinate{start.y + i, 0 + start.x})
		steps = append(steps, Coordinate{start.y + 0, i + start.x})
		steps = append(steps, Coordinate{start.y + -1*i, 0 + start.x})
		steps = append(steps, Coordinate{start.y + 0, -1*i + start.x})
	}
	return steps
}

func getValue(piece PieceType) float32 {
	switch piece {
	case Pawn:
		return 1
	case Knight:
		return 3
	case Bishop:
		return 3
	case Rook:
		return 5
	case Queen:
		return 9
	case King:
		return 0 // king should not be counted
	}
	panic("unaccounted option")
}

func isOnBoard(coord Coordinate) bool {
	return (0 <= coord.y && coord.y < 8 && 0 <= coord.x && coord.x < 8)
}

func getMoveType(move AbsoluteMove) MoveType {
	xDiff := move.start.y - move.end.y
	yDiff := move.start.x - move.end.x

	if move.start == move.end {
		return InvalidMove
	} else if move.start.y == move.end.y && move.start.x != move.end.x || move.start.y != move.end.y && move.start.x == move.end.x {
		return LinearMove
	} else if xDiff == yDiff || move.end.y-move.start.y == yDiff {
		return DiagonalMove
	} else if ((xDiff == 2 || xDiff == -2) && (yDiff == 1 || yDiff == -1)) || ((xDiff == 1 || xDiff == -1) && (xDiff == 2 || xDiff == -1)) {
		return HorseMove
	} else {
		return InvalidMove
	}
}
