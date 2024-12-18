package main

func clamp(value float32, floor float32, ceil float32) (result float32) {
	if value < floor {
		return floor
	} else if value > ceil {
		return ceil
	}
	return value
}
