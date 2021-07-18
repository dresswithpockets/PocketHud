package main

import (
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "image/color"
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

type RoundCorner int

const (
    RoundCornerTopLeft RoundCorner = 1 << iota
    RoundCornerTopRight
    RoundCornerBottomLeft
    RoundCornerBottomRight
    RoundCornerAll = RoundCornerTopLeft | RoundCornerTopRight | RoundCornerBottomLeft | RoundCornerBottomRight
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

type Inset struct {
    left, top, right, bottom int
}

type Panel struct {
    surface  *Surface
    parent   *Panel
    children []*Panel
    target   pixel.Target
    scheme   *Scheme

    panelName      string
    panelFlags     Flag
    buildModeFlags Flag
    enabled        bool

    border              *SchemeBorder
    paintBackgroundType BackgroundType

    size           Size
    pos            Position
    zpos           int16
    visible        bool
    inset          Inset
    roundedCorners RoundCorner
    fgColor        color.Color
    bgColor        color.Color

    mouseInput  bool
    kbInput     bool
    tabPosition int
}

func (p *Panel) ApplySettings(object *vgui.Object) {
    if p.panelFlags.Has(PanelNeedsDefaultSettingsApplied) {
        // TODO InternalInitDefaultValues. We don't know what the defaults are from GetAnimMap()
    }

    // TODO InternalApplySettings seems to ultimately set hud textures, not sure

    p.buildModeFlags.Clear(BuildModeSaveXposRightAligned |
        BuildModeSaveXposCenterAligned |
        BuildModeSaveYposBottomAligned |
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
    // TODO screenSize := alignScreenSize for proportional/test title safe area
    // TODO fullscreen dimensions by removing override?

    // TODO parentPos := Position{0, 0} for proportional/test title safe area

    // flag to cause windows to get screenSize from their parents,
    // this allows children windows to use fill and right/bottom alignment even
    // if their parent does not use the full screen.
    if object.GetBoolD("proportionalToParent", false) {
        p.buildModeFlags.Set(BuildModeSaveProportionalToParent)
        if p.parent != nil {
            bounds := p.parent.GetBounds()
            // TODO parentPos = bounds.Position for proportional/test title safe area
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

    // TODO usedTitleSafeArea, panel title safe area, x360 mode

    // TODO navigation simulation (SetNavX where X is Up/Down/Left/Right/etc) for `navX` resource properties

    p.SetPos(pos)

    if zpos, ok := object.GetInt("zpos"); ok {
        p.SetZPos(int16(zpos))
    }

    // TODO: if UsesTitleSafeArea, handle WIDE_FULL, TALL_FULL build flags & mutate size based on those flags

    p.SetSize(size)

    // NOTE this has to happen after pos + size is set
    // TODO ApplyAutoResizeSettings(object)

    if object.GetBoolD("IgnoreScheme", false) {
        // TODO PerformApplySchemeSettings
    }

    // panel state
    p.SetVisible(object.GetBoolD("visible", true))
    p.SetEnabled(object.GetBoolD("enabled", true))

    p.SetMouseInputEnabled(object.GetBoolD("mouseinputenabled", true))

    p.SetTabPosition(object.GetIntD("tabPosition", 0))

    // TODO tooltiptext

    paintBackground := object.GetIntD("paintbackground", -1)
    if paintBackground >= 0 {
        p.SetPaintBackgroundEnabled(paintBackground != 0)
    }

    paintBorder := object.GetIntD("paintborder", -1)
    if paintBorder >= 0 {
        p.SetPaintBorderEnabled(paintBorder != 0)
    }

    if border, ok := object.GetString("border"); ok {
        p.SetBorder(p.GetScheme().GetBorder(border))
    }

    if newName, ok := object.GetString("fieldName"); ok {
        p.SetName(newName)
    }

    // TODO actionsignallevel (telemetry)

    // TODO forceStereoRenderToFrameBuffer

    // this is a flag int, can be anything between 0b0 and 0b1000. See: type RoundCorner
    if roundedCorners, ok := object.GetInt("RoundedCorners"); ok {
        p.roundedCorners = RoundCorner(roundedCorners)
    }

    // TODO pin corners to siblings/siblings to corners
    // TODO overridableColorEntries

    p.SetKeyboardInputEnabled(object.GetBoolD("keyboardinputenabled", true))

    // TODO OnChildSettingsApplied event
}

func (p *Panel) PaintTraverse() {
    // TODO
}

func (p *Panel) PaintBorder() {
    if p.border == nil {
        return
    }
    p.border.PaintFromPanel(p)
}

func (p *Panel) PaintBackground() {
    size := p.GetSize()
    // TODO SkipChild?
    // if ( m_SkipChild.Get() && m_SkipChild->IsVisible() ) {}
    // else
    {
        vTopLeft := pixel.ZV
        vBottomRight := vTopLeft.Add(size.Vec())
        topLeft := vguiToPixelCoords(vTopLeft)
        bottomRight := vguiToPixelCoords(vBottomRight)
        topRight := pixel.V(bottomRight.X, topLeft.Y)
        bottomLeft := pixel.V(topLeft.X, bottomRight.Y)
        switch p.paintBackgroundType {
        case BgFilled:
            p.DrawFilledBox(p.bgColor, topLeft, topRight, bottomRight, bottomLeft)
        case BgTextured:
            p.DrawTexturedBox(p.bgColor, 1, topLeft, topRight, bottomRight, bottomLeft)
        case BgRoundedCorners:
            p.DrawBox(p.bgColor, 1, topLeft, topRight, bottomRight, bottomLeft)
            // TODO case 3 DrawBoxFade? this seems to be explicitly unsupported by the vgui Panel type, but maybe its allowed in child types
        }
    }
}

func (p *Panel) DrawFilledBox(col color.Color, corners ...pixel.Vec) {
    d := imdraw.New(nil)
    d.SetColorMask(col)
    // push the corners of the rectangle to be filled
    d.Push(corners...)
    // draw filled rectangle
    d.Rectangle(0)
}

func (p *Panel) DrawTexturedBox(col color.Color, normalizedAlpha float64, corners ...pixel.Vec) {

    r, g, b, a := col.RGBA()
    a = uint32(float64(a) / 255.0 * normalizedAlpha)
    col = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

    // TODO background texture/image
    d := imdraw.New(nil)
    d.SetColorMask(col)
    d.Push(corners...)
    d.Draw(p.target)
}

func (p *Panel) DrawBox(col color.Color, normalizedAlpha float64, corners ...pixel.Vec) {
    // TODO: Draw Textured Rounded Box
}

func (p *Panel) Paint() {}

func (p *Panel) SetName(panelName string) {
    p.panelName = panelName
}

func (p *Panel) GetScheme() *Scheme {
    if p.scheme != nil {
        return p.scheme
    }

    if p.parent != nil {
        return p.parent.GetScheme()
    }

    return GetDefaultScheme()
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

        wideInt, err := strconv.ParseInt(wide, 10, 16)
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
            warningLogger.Println("Cannot handle Proportional width when computing panel size: GetProportionalScaledValueEx not implemented.")
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportional) {
            // TODO GetProportionalScaledValueEx
            warningLogger.Println("Cannot handle Proportional width when computing panel size: GetProportionalScaledValueEx not implemented.")
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportionalSelf) {
            width = p.GetSize().width * int16(wideInt)
        } else {
            if p.IsProportional() {
                // scale the width up to our screen coords
                // TODO GetProportionalScaledValueEx
                warningLogger.Println("Cannot handle Proportional width when computing panel size: GetProportionalScaledValueEx not implemented.")
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
                    // TODO panel name
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

        tallInt, err := strconv.ParseInt(tall, 10, 16)
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
            warningLogger.Println("Cannot handle Proportional height when computing panel size: GetProportionalScaledValueEx not implemented.")
        } else if p.buildModeFlags.Has(BuildModeSaveTallProportional) {
            // TODO GetProportionalScaledValueEx
            warningLogger.Println("Cannot handle Proportional height when computing panel size: GetProportionalScaledValueEx not implemented.")
        } else if p.buildModeFlags.Has(BuildModeSaveWideProportionalSelf) {
            height = p.GetSize().height * int16(tallInt)
        } else {
            if p.IsProportional() {
                // scale the height up to our screen coords
                // TODO GetProportionalScaledValueEx
                warningLogger.Println("Cannot handle Proportional height when computing panel size: GetProportionalScaledValueEx not implemented.")
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
            warningLogger.Println("Cannot handle Proportional layout when computing the panels position: GetProportionalScaledValueEx not implemented.")
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
            newPos = parentSize/2 + posDelta
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

func (p *Panel) SetSize(size Size) {
    // TODO vpanel minimum size
    if p.size == size {
        return
    }
    p.size = size
    // TODO OnSizeChanged event
}

func (p *Panel) GetPos() Position {
    return p.pos
}

func (p *Panel) SetPos(pos Position) {
    p.pos = pos
}

func (p *Panel) GetZPos() int16 {
    return p.zpos
}

func (p *Panel) SetZPos(z int16) {
    p.zpos = z
}

func (p *Panel) SetVisible(visible bool) {
    if p.visible == visible {
        return
    }

    // TODO surface->SetPanelVisible (in case special window processing needs to occur)?

    p.visible = visible

    // TODO if IsPopup, CalculateMouseVisible() ?
}

func (p *Panel) SetEnabled(enabled bool) {
    if p.enabled == enabled {
        return
    }

    p.enabled = enabled
    // TODO InvalidateLayout(false)
    // TODO Repaint()
}

func (p *Panel) SetMouseInputEnabled(enabled bool) {
    p.mouseInput = enabled
    // TODO surface()->CalculateMouseVisible()
}

func (p *Panel) SetKeyboardInputEnabled(enabled bool) {
    p.kbInput = enabled
    for _, child := range p.children {
        child.SetKeyboardInputEnabled(enabled)
    }
    // TODO if turning kb input off, make sure this panel is not the current focus of a parent panel
}

func (p *Panel) GetTabPosition() int {
    return p.tabPosition
}

func (p *Panel) SetTabPosition(position int) {
    p.tabPosition = position
}

func (p *Panel) SetPaintBackgroundEnabled(enabled bool) {
    if enabled {
        p.panelFlags.Set(PanelPaintBackgroundEnabled)
    } else {
        p.panelFlags.Clear(PanelPaintBackgroundEnabled)
    }
}

func (p *Panel) SetPaintBorderEnabled(enabled bool) {
    if enabled {
        p.panelFlags.Set(PanelPaintBorderEnabled)
    } else {
        p.panelFlags.Clear(PanelPaintBorderEnabled)
    }
}

func (p *Panel) GetInset() Inset {
    return p.inset
}

func (p *Panel) SetInset(inset Inset) {
    p.inset = inset
}

func (p *Panel) SetBorder(border *SchemeBorder) {
    p.border = border

    if border != nil {
        p.SetInset(border.inset)

        // update background type based on the border
        p.SetPaintBackgroundType(border.backgroundType)
    } else {
        p.SetInset(Inset{0, 0, 0, 0})
    }
}

func (p *Panel) SetPaintBackgroundType(backgroundType BackgroundType) {
    p.paintBackgroundType = backgroundType
}
