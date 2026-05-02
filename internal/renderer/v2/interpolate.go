package v2

import (
	"github.com/cnkei/gospline"
)

// interpolate resamples data to targetWidth using cubic spline interpolation.
// Falls back to linear (or constant) for very short inputs.
func interpolate(data []float64, targetWidth int) []float64 {
	if len(data) == 0 {
		return make([]float64, targetWidth)
	}
	if targetWidth <= 0 {
		return []float64{}
	}
	if len(data) == 1 {
		result := make([]float64, targetWidth)
		for i := range result {
			result[i] = data[0]
		}
		return result
	}
	// For 2 points, linear is sufficient and avoids degenerate spline
	if len(data) == 2 {
		return interpolate(data, targetWidth) // reuse original linear if you keep it
	}

	// Original x-coordinates: equally spaced 0, 1, ..., n-1
	n := len(data) - 1
	x := make([]float64, len(data))
	for i := range x {
		x[i] = float64(i)
	}
	y := data // copy not needed if we don't mutate

	// Create natural cubic spline (2nd derivative = 0 at endpoints)
	s := gospline.NewCubicSpline(x, y) // or NewNaturalCubicSpline, NewMonotoneSpline, etc.

	// Evaluate at new equally spaced points
	result := make([]float64, targetWidth)
	for i := range result {
		t := float64(i) / float64(targetWidth-1) * float64(n) // map [0..targetWidth-1] → [0..n]
		result[i] = s.At(t)
	}

	return result
}
