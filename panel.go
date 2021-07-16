package main

import (
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "strconv"
    "strings"
    "unicode"
)

const (
    PanelMarkedForDeletion Flag = 1 << iota
    PanelNeedsRepaint
    PanelPaintBorderEnabled
    PanelPaintBackgroundEnabled
    PanelPaintEnabled
    PanelPostChildPaintEnabled
    PanelAutoDeleteEnabled
    PanelNeedsLayout
    PanelNeedsSchemeUpdate
    PanelNeedsDefaultSettingsApplied
    PanelAllowChainKeybindingToParent
    PanelInPerformLayout
    PanelIsProportional
    PanelTriplePressAllowed
    PanelDragRequiresPanelExit
    PanelIsMouseDisabledForThisPanelOnly
)

const (
    BuildModeEditable Flag = iota
    BuildModeDeletable
    BuildModeSaveXposRightAligned
    BuildModeSaveXposCenterAligned
    BuildModeSaveYposBottomAligned
    BuildModeSaveYposCenterAligned
    BuildModeSaveWideFull
    BuildModeSaveTallFull
    BuildModeSaveProportionalToParent
    BuildModeSaveWideProportional
    BuildModeSaveTallProportional
    BuildModeSaveXposProportionalSelf
    BuildModeSaveYposProportionalSelf
    BuildModeSaveWideProportionalTall
    BuildModeSaveTallProportionalWide
    BuildModeSaveXposProportionalParent
    BuildModeSaveYposProportionalParent
    BuildModeSaveWideProportionalSelf
    BuildModeSaveTallProportionalSelf
)

type AspectRatio int

const (
    Aspect16x9 AspectRatio = iota
    Aspect4x3
    Aspect5x4
)

type ComputeOperator int

const (
    OpAdd ComputeOperator = iota
    OpSub
    OpSet
)

type Surface struct {
    target pixel.Picture
    aspect AspectRatio
}

// GetSize returns the VGUI size of the surface.
// The height will always be VguiHeight.
// Depending on the aspect ratio, the width will be one of Vgui16x9Width, Vgui4x3Width, or Vgui5x4Width
func (s *Surface) GetSize() Size {
    h := int16(VguiHeight)
    w := map[AspectRatio]int16{Aspect16x9: Vgui16x9Width, Aspect4x3: Vgui4x3Width, Aspect5x4: Vgui5x4Width}[s.aspect]
    return Size{w, h}
}

type Panel struct {
    surface *Surface
    parent  *Panel
    target  pixel.Picture

    panelFlags     Flag
    buildModeFlags Flag

    size Size
    pos Position
}

func (p *Panel) ApplySettings(object *vgui.Object) {
    if p.panelFlags.Has(PanelNeedsDefaultSettingsApplied) {
        // TODO InternalInitDefaultValues. We don't know what the defaults are from GetAnimMap()
    }

    // TODO InternalApplySettings seems to ultimately set hud textures, not sure

    p.buildModeFlags.Clear(BuildModeSaveXposRightAligned |
        BuildModeSaveXposCenterAligned |
        BuildModeSaveYposRightAligned |
        BuildModeSaveYposCenterAligned |
        BuildModeSaveWideFull |
        BuildModeSaveTallFull |
        BuildModeSaveProportionalToParent |
        BuildModeSaveWideProportional |
        BuildModeSaveTallProportional |
        BuildModeSaveXposProportionalSelf |
        BuildModeSaveYposProportionalSelf |
        BuildModeSaveWideProportionalTall |
        BuildModeSaveTallProportionalWide |
        BuildModeSaveXposProportionalParent |
        BuildModeSaveYposProportionalParent |
        BuildModeSaveWideProportionalSelf |
        BuildModeSaveTallProportionalSelf)

    // get the position
    alignScreenSize := p.surface.GetSize()
    screenSize := alignScreenSize
    // TODO fullscreen dimensions by removing override?

    parentPos := Position{0,0}

    // flag to cause windows to get screenSize from their parents,
    // this allows children windows to use fill and right/bottom alignment even
    // if their parent does not use the full screen.
    if object.GetBoolD("proportionalToParent", false) {
        p.buildModeFlags.Set(BuildModeSaveProportionalToParent)
        if p.parent != nil {
            bounds := p.parent.GetBounds()
            parentPos = bounds.Position
            alignScreenSize = bounds.Size
        }
    }

    width := p.ComputeWidth(object, alignScreenSize, false)
    height := p.ComputeHeight(object, alignScreenSize, false)
    size := Size{width, height}

    pos := p.GetPos()
    {
        x, buildFlagsX := p.ComputePos(object, pos.x, size.width, alignScreenSize.width, true, OpSet)
        y, buildFlagsY := p.ComputePos(object, pos.y, size.height, alignScreenSize.height, false, OpSet)
        p.buildModeFlags.Set(buildFlagsX | buildFlagsY)
        pos.x = x
        pos.y = y
    }

    usesTitleSafeArea := false

    // TODO panel_test_title_safe

    // TODO panel ApplySettings
}

func (p *Panel) GetBounds() Bounds {
    panic("GetBounds not implemented yet")
}

func (p *Panel) getInset() (left, top, right, bottom int16) {
    // TODO panel inset is based on set border, otherwise its (0,0,0,0). calculate this when border is set, not here.
    return 0, 0, 0, 0
}

func (p *Panel) ComputeWidth(object *vgui.Object, parentSize Size, computingOther bool) int16 {
    width := p.GetSize().width

    if wide, ok := object.GetString("wide"); ok {
        lower := strings.ToLower(wide)
        if lower[0] == 'f' {
            p.buildModeFlags.Set(BuildModeSaveWideFull)
            wide = wide[1:]
        } else {
            if lower[0] == 'o' {
                wide = wide[1:]
                if computingOther {
                    // TODO Wide and Tall of panel %s are set to be each other! (see: vgui2\vgui_controls\Panel.cpp:8656
                    return 0
                }

                p.buildModeFlags.Set(BuildModeSaveWideProportionalTall)
                width = p.ComputeHeight(object, parentSize, true)

                if p.IsProportional() {
                    // TODO GetProportionalNormalizedValue(width)
                }
            } else if lower[0] == 'p' {
                p.buildModeFlags.Set(BuildModeSaveWideProportional)
                wide = wide[1:]
            } else if lower[0] == 's' {
                p.buildModeFlags.Set(BuildModeSaveWideProportionalSelf)
                wide = wide[1:]
            }
        }

        wideInt, err := strconv.ParseInt(wide, 10,16)
        if err != nil {
            // just do what atof does - vgui source uses atof and atoi, but the input is always an integer so just mix
            // the two here basically
            wideInt = 0
        }

        if !p.buildModeFlags.Has(BuildModeSaveWideProportionalTall) {
            width = int16(wideInt)
        }

        if p.buildModeFlags.Has(BuildModeSaveWideProportionalTall) {
            // TODO: GetProportionalScaledValueEx
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportional) {
            // TODO GetProportionalScaledValueEx
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportionalSelf) {
            width = p.GetSize().width * int16(wideInt)
        } else {
            if p.IsProportional() {
                // scale the width up to our screen coords
                // TODO GetProportionalScaledValueEx
            }

            // correct the alignment
            if p.buildModeFlags.Has(BuildModeSaveWideFull) {
                width = parentSize.width - width
            }
        }
    }

    return width
}

func (p *Panel) ComputeHeight(object *vgui.Object, parentSize Size, computingOther bool) int16 {
    height := p.GetSize().height

    if tall, ok := object.GetString("tall"); ok {
        lower := strings.ToLower(tall)
        if lower[0] == 'f' {
            p.buildModeFlags.Set(BuildModeSaveTallFull)
            tall = tall[1:]
        } else {
            if lower[0] == 'o' {
                tall = tall[1:]
                if computingOther {
                    // TODO Wide and Tall of panel %s are set to be each other! (see: vgui2\vgui_controls\Panel.cpp:8656
                    return 0
                }

                p.buildModeFlags.Set(BuildModeSaveTallProportionalWide)
                height = p.ComputeWidth(object, parentSize, true)

                if p.IsProportional() {
                    // TODO GetProportionalNormalizedValue(width)
                }
            } else if lower[0] == 'p' {
                p.buildModeFlags.Set(BuildModeSaveTallProportional)
                tall = tall[1:]
            } else if lower[0] == 's' {
                p.buildModeFlags.Set(BuildModeSaveTallProportionalSelf)
                tall = tall[1:]
            }
        }

        tallInt, err := strconv.ParseInt(tall, 10,16)
        if err != nil {
            // just do what atof does - vgui source uses atof and atoi, but the input is always an integer so just mix
            // the two here basically
            tallInt = 0
        }

        if !p.buildModeFlags.Has(BuildModeSaveTallProportionalWide) {
            height = int16(tallInt)
        }

        if p.buildModeFlags.Has(BuildModeSaveTallProportionalWide) {
            // TODO: GetProportionalScaledValueEx
        } else if p.buildModeFlags.Has(BuildModeSaveTallProportional) {
            // TODO GetProportionalScaledValueEx
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportionalSelf) {
            height = p.GetSize().height * int16(tallInt)
        } else {
            if p.IsProportional() {
                // scale the width up to our screen coords
                // TODO GetProportionalScaledValueEx
            }

            // correct the alignment
            if p.buildModeFlags.Has(BuildModeSaveTallFull) {
                height = parentSize.height - height
            }
        }
    }

    return height
}

func (p *Panel) ComputePos(object *vgui.Object, pos int16, size int16, parentSize int16, isX bool, op ComputeOperator) (outPos int16, flags Flag) {
    outPos = pos
    flags = 0

    flagRightAlign := BuildModeSaveYposBottomAligned
    flagCenterAlign := BuildModeSaveYposCenterAligned
    flagProportionalSelf := BuildModeSaveYposProportionalSelf
    flagProportionalParent := BuildModeSaveYposProportionalParent
    propertyName := "ypos"
    if isX {
        flagRightAlign = BuildModeSaveXposRightAligned
        flagCenterAlign = BuildModeSaveXposCenterAligned
        flagProportionalSelf = BuildModeSaveXposProportionalSelf
        flagProportionalParent = BuildModeSaveXposProportionalParent
        propertyName = "xpos"
    }

    var posDelta int16 = 0
    inputStr, ok := object.GetString(propertyName)
    if ok {
        lower := strings.ToLower(inputStr)
        // look for alignment flags
        if lower[0] == 'r' {
            flags.Set(flagRightAlign)
            inputStr = inputStr[1:]
        } else if lower[0] == 'c' {
            flags.Set(flagCenterAlign)
            inputStr = inputStr[1:]
        }

        if inputStr[0] == 's' {
            flags.Set(flagProportionalSelf)
            inputStr = inputStr[1:]
        } else {
            flags.Set(flagProportionalParent)
            inputStr = inputStr[1:]
        }

        parsedPos, err := strconv.ParseInt(inputStr, 10, 16)
        if err != nil {
            parsedPos = 0
        }
        newPos := int16(parsedPos)

        if p.IsProportional() {
            // TODO: GetProportionalScaledValueEx
        }

        if flags.Has(flagProportionalSelf) {
            posDelta = size * newPos
        } else if flags.Has(flagProportionalParent) {
            posDelta = parentSize * newPos
        } else {
            posDelta = newPos
        }

        // correct alignment
        if flags.Has(flagRightAlign) {
            newPos = parentSize - posDelta
        } else if flags.Has(flagCenterAlign) {
            newPos = parentSize / 2 + posDelta
        } else {
            newPos = posDelta
        }

        switch op {
        case OpAdd:
            outPos += newPos
        case OpSub:
            outPos -= newPos
        case OpSet:
            outPos = newPos
        }

        if inputStr[0] == '-' || inputStr[0] == '+' {
            inputStr = inputStr[1:]
        }

        // TODO handle floating point x/ypos
        for len(inputStr) > 0 && unicode.IsDigit(rune(inputStr[0])) {
            inputStr = inputStr[1:]
        }

        if len(inputStr) > 0 {
            switch inputStr[0] {
            case '+':
                outPos, flags = p.ComputePos(object, outPos, size, parentSize, isX, OpAdd)
            case '-':
                outPos, flags = p.ComputePos(object, outPos, size, parentSize, isX, OpSub)
            }
        }
    }
    return
}

func (p *Panel) IsProportional() bool {
    return p.panelFlags.Has(PanelIsProportional)
}

func (p *Panel) GetSize() Size {
    return p.size
}

func (p *Panel) GetPos() Position {
    return p.pos
}