package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
)

const (
    VguiHeight = 480
    Vgui16x9Width = 852
    Vgui4x3Width = 640
    Vgui5x4Width = 600
)

// vguiToPixelCoords converts vgui coords to pixel coords for rendering
// vgui origin is top left, and pixel origin is bottom left
func vguiToPixelCoords(v pixel.Vec) pixel.Vec {
    return pixel.V(v.X, VguiHeight - v.Y)
}

// drawDashedLine helper function, draws a simple line with dashing parameters.
// Sourced from vgui2\vgui_controls\Label.cpp:629
func drawDashedLine(imd *imdraw.IMDraw, x0, y0, x1, y1, dashLen, gapLen int16) {
    if x1 - x0 > y1 - y0 {
        // horizontal direction line
        for {
            if x0 + dashLen > x1 {
                imd.Push(pixel.V(float64(x0), float64(y0)), pixel.V(float64(x1), float64(y1)))
                imd.Rectangle(0)
            } else {
                imd.Push(pixel.V(float64(x0), float64(y0)), pixel.V(float64(x0 + dashLen), float64(y1)))
                imd.Rectangle(0)
            }
            x0 += dashLen
            if x0 + gapLen > x1 {
                break
            }
            x0 += gapLen
        }
    } else {
        // vertical direction line
        for {
            if y0 + dashLen > y1 {
                imd.Push(pixel.V(float64(x0), float64(y0)), pixel.V(float64(x1), float64(y1)))
                imd.Rectangle(0)
            } else {
                imd.Push(pixel.V(float64(x0), float64(y0)), pixel.V(float64(x1), float64(y0 + dashLen)))
                imd.Rectangle(0)
            }
            y0 += dashLen
            if y0 + gapLen > y1 {
                break
            }
            y0 += gapLen
        }
    }
}