package main

func clamp(value float32, floor float32, ceil float32) (result float32) {
	if value < floor {
		return floor
	} else if value > ceil {
		return ceil
	}
	return value
}

func clampCoord(coord Coordinate) Coordinate {
	if coord.x > 7 {
		coord.x = 7
	} else if coord.x < 0 {
		coord.x = 0
	}
	if coord.y > 7 {
		coord.y = 7
	} else if coord.y < 0 {
		coord.y = 0
	}
	return coord
}
