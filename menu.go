package main

import "github.com/dresswithpockets/go-vgui"

type Element struct {
    name     string
    controls map[string]Control
}

type Menu struct {
    file     string
    elements map[string]*Element
}

func (c *ControlProvider) newElementFromObject(object *vgui.Object) *Element {
    controls := map[string]Control{}
    for k, v := range object.Properties {
        control, err := c.resolveNewControlFromObject(v)
        if err != nil {
            // TODO: handle erroneous control from object
        }
        controls[k] = control
    }
    return &Element{object.Name, controls}
}

func (c *ControlProvider) newMenuFromObject(object *vgui.Object) *Menu {
    elements := map[string]*Element{}
    for k, v := range object.Properties {
        elements[k] = c.newElementFromObject(v)
    }
    // TODO: get file from object
    return &Menu{"", elements}
}

func (b *Menu) drawMenu(app *App) {
    panic("Draw not implemented on Menu.")
}
