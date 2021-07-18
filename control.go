package main

import (
    . "github.com/ahmetb/go-linq"
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "math"
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

const FullSize = math.MaxInt16

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
    x int16
    y int16
}

type Size struct {
    width  int16
    height int16
}

type Bounds struct {
    Position
    Size
}

type VguiImage struct {
    name    string
    picture pixel.Picture
}

type ControlBuilder func(object *vgui.Object) Control

type ControlProvider struct {
    builders map[string]ControlBuilder
}

type Control interface {
    applySettings(object *vgui.Object)
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
    tabPosition     int
    dullText        bool
    brightText      bool
    image           *VguiImage
    paintBackground bool
    paintBorder     bool
    _default         int

    parent   Control
    children []Control

    // indicates that we need to recalculate the absolute bounds of this control
    dirty          bool
    absoluteBounds pixel.Rect

    baseOverride *BaseControl
}

func (p Position) Vec() pixel.Vec {
    return pixel.V(float64(p.x), float64(p.y))
}

func (s Size) Vec() pixel.Vec {
    return pixel.V(float64(s.width), float64(s.height))
}

func (c *ControlProvider) setBuilder(controlName string, builder ControlBuilder) {
    c.builders[controlName] = builder
}

func (c *ControlProvider) resolveNewControlFromObject(object *vgui.Object) (Control, error) {
    controlName, ok := object.Get("ControlName")
    if !ok {
        panic("ControlName not found on object when resolving new control from object. Default behaviour not well defined.")
    }

    if !controlName.IsValue {
        panic("ControlName must always be a single value, not an object with properties.")
    }

    if builder, ok := c.builders[controlName.Value]; ok {
        return builder(object), nil
    }

    return nil, &ErrUnknownControlName{controlName.Value}
}

func (c *BaseControl) applySettings(object *vgui.Object) {
    panic("applySettings not implemented on abstract BaseControl")
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
    left := parentCenter.X - parentSize.X/2
    right := parentCenter.X + parentSize.Y/2
    top := parentCenter.Y - parentSize.Y/2
    bottom := parentCenter.Y + parentSize.Y/2

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
    c.absoluteBounds = pixel.R(xpos, ypos, xpos+float64(c.size.width), ypos+float64(c.size.height))
}

/*func defaultControlBuilder(t reflect.Type, object *vgui.Object) Control {
    var control Control
    t.AssignableTo(reflect.TypeOf(control))
    value := reflect.New(t)
    for i := 0; i < value.NumField(); i++ {
        field := value.Field(i)
    }
    // TODO: parse object values into fields
}

func parseIntValue(value string) (reflect.Value, error) {
    i, err := strconv.ParseInt(value, 10, 16)
    return reflect.ValueOf(i), err
}*/

// TODO parse bools, floats, special string types, etc