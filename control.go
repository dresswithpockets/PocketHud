package main

import (
    . "github.com/ahmetb/go-linq"
    "github.com/faiface/pixel"
)

type Alignment int

const (
    AlignCenter Alignment = iota
    AlignLeft
    AlignRight
    AlignNorth
    AlignSouth
    AlignWest
    AlignEast
    AlignNorthWest
    AlignNorthEast
    AlignSouthWest
    AlignSouthEast
)

type RelativeTo int

const (
    RelativeLeftOrTop RelativeTo = iota
    RelativeRightOrBottom
    RelativeCenter
)

type RelativeInt struct {
    value      int
    relativeTo RelativeTo
}

type ControlPosition struct {
    x RelativeInt
    y RelativeInt
}

type Position struct {
    x int
    y int
}

type Size struct {
    width  int
    height int
}

type VguiImage struct {
    name    string
    picture pixel.Picture
}

type Control interface {
    zOrder() int
    draw()
    drawChildren()
    setParent(other Control)
    getBounds() pixel.Rect
    recalculateBounds()
}

type BaseControl struct {
    app             *App
    name            string
    fieldName       string
    pos             ControlPosition
    z               int
    size            Size
    visible         bool
    enabled         bool
    fgColor         *SchemeColor
    bgColor         *SchemeColor
    border          *SchemeBorder
    font            *SchemeFont
    labelText       string
    textAlignment   Alignment
    textInset       Position
    image           *VguiImage
    paintBackground bool
    paintBorder     bool

    parent   Control
    children []Control

    // indicates that we need to recalculate the absolute bounds of this control
    dirty          bool
    absoluteBounds pixel.Rect

    // TODO: overrides (props with _override suffix)
}

func (c *BaseControl) zOrder() int {
    return c.z
}

// TODO: from vgui.Value
func (c *BaseControl) draw() {
    panic("Draw not implemented on the abstract BaseControl.")
}

// drawChildren draws all sub-controls of this control according to their zOrder
//goland:noinspection SpellCheckingInspection
func (c *BaseControl) drawChildren() {
    var sortedChildren []Control
    From(c.children).Sort(func(a, b interface{}) bool {
        actl := a.(Control)
        bctl := b.(Control)
        return actl.zOrder() < bctl.zOrder()
    }).ToSlice(&sortedChildren)
    for _, child := range sortedChildren {
        child.draw()
        child.drawChildren()
    }
}

func (c *BaseControl) setParent(other Control) {
    if c.parent == other {
        return
    }

    if c.parent != nil {
        parentBase := c.parent.(*BaseControl)

        // get the index of our control in the parent's children slice and the remove it from the slice
        // we have to do it this way in order to maintain order which is at least *somewhat* important for
        // controls on the same z-plane
        for i, v := range parentBase.children {
            if v == c {
                parentBase.children = append(parentBase.children[:i], parentBase.children[i+1:]...)
                break
            }
        }
    }

    // set the parent on our control & append our control to the parent's children list
    c.parent = other
    base := c.parent.(*BaseControl)
    base.children = append(base.children, c)

    c.dirty = true
}

func (c *BaseControl) getBounds() pixel.Rect {
    if c.dirty {
        c.recalculateBounds()
        c.dirty = false
    }
    return c.absoluteBounds
}

func (c *BaseControl) recalculateBounds() {
    viewport := c.app.window.Bounds()
    if c.parent != nil {
        viewport = c.parent.getBounds()
    }

    parentCenter := viewport.Center()
    parentSize := viewport.Size()
    left := parentCenter.X - parentSize.X / 2
    right := parentCenter.X + parentSize.Y / 2
    top := parentCenter.Y - parentSize.Y / 2
    bottom := parentCenter.Y + parentSize.Y / 2

    var xpos, ypos float64

    switch c.pos.x.relativeTo {
    case RelativeLeftOrTop:
        xpos = left
        break
    case RelativeCenter:
        xpos = parentCenter.X
        break
    case RelativeRightOrBottom:
        xpos = right
        break
    }

    switch c.pos.y.relativeTo {
    case RelativeLeftOrTop:
        ypos = top
        break
    case RelativeCenter:
        ypos = parentCenter.Y
        break
    case RelativeRightOrBottom:
        ypos = bottom
        break
    }

    xpos += float64(c.pos.x.value)
    ypos += float64(c.pos.y.value)
    c.absoluteBounds = pixel.R(xpos, ypos, xpos + float64(c.size.width), ypos + float64(c.size.height))
}