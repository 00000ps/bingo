package cv

import "bingo/pkg/utils"

type (
	// Rect means a rectangle
	Rect struct {
		X, Y, Width, Height int
	}
	// Location means x & y int value
	Location struct {
		X, Y int
	}
)

// Area returns the rect area
func (r *Rect) Area() int { return r.Width * r.Height }

// LB returns the rect left-bottom location
func (r *Rect) LB() Location {
	return Location{
		X: r.X,
		Y: r.Y + r.Height,
	}
}

// RT returns the rect left-bottom location
func (r *Rect) RT() Location {
	return Location{
		X: r.X + r.Width,
		Y: r.Y,
	}
}

// RB returns the rect left-bottom location
func (r *Rect) RB() Location {
	return Location{
		X: r.X + r.Width,
		Y: r.Y + r.Height,
	}
}

// IoU means Intersection over Union of rectangle
func (r *Rect) IoU(a *Rect) float32 { return IoU(r, a) }

// IoU means Intersection over Union of rectangle
func IoU(a, b *Rect) float32 {
	W := utils.Min(a.RT().X, b.RT().X) - utils.Max(a.LB().X, b.LB().X)
	H := utils.Min(a.RT().Y, b.RT().Y) - utils.Max(a.LB().Y, b.LB().Y)
	if W <= 0 || H <= 0 {
		return 0
	}
	SA := (a.RT().X - a.LB().X) * (a.RT().Y - a.LB().Y)
	SB := (b.RT().X - b.LB().X) * (b.RT().Y - b.LB().Y)
	cross := W * H
	return float32(cross) / float32(SA+SB-cross)
	// return 0
}
