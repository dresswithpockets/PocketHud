package main

import (
    "golang.org/x/image/font"
    "image/color"
)

type Scheme struct {
    colors map[string]*SchemeColor
    fonts map[string]*SchemeColor
    borders map[string]*SchemeColor
}

type SchemeColor struct {
    name  string
    color color.Color
}

type SchemeFont struct {
    name string
    face font.Face
}

type SchemeBorder struct {
    name string
    // TODO: border definition
}