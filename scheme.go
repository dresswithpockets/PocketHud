package main

import (
    "golang.org/x/image/font"
    "image/color"
)

type Scheme struct {
    colors  map[string]*SchemeColor
    fonts   map[string]*SchemeFont
    borders map[string]*SchemeBorder

    baseBorder *SchemeBorder
}

type SchemeColor struct {
    name  string
    color color.Color
}

type SchemeFont struct {
    name string
    face font.Face
}

type BackgroundType int

const (
    BgFilled BackgroundType = iota
    BgTextured
    BgRoundedCorners
)

type SchemeBorder struct {
    name string

    inset          Inset
    backgroundType BackgroundType
    // TODO: border definition
}

var defaultScheme *Scheme

// TODO initializing default schemes and such, see LoadSchemeFromFileEx

func GetDefaultScheme() *Scheme {
    return defaultScheme
}

// TODO LoadBorders, baseBorder

func (s *Scheme) GetBorder(border string) *SchemeBorder {
    if b, ok := s.borders[border]; ok {
        return b
    }
    return s.baseBorder
}

func (s *SchemeBorder) PaintFromPanel(panel *Panel) {
    panic("border paint functionality has not been implemented")
}
