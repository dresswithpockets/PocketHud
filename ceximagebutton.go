package editor

type CExImageButton struct {
    BaseControl
    soundDepressed  string
    soundReleased   string
    borderDefault   *SchemeBorder
    borderArmed     *SchemeBorder
    paintBackground bool

    overrides *CExImageButton
}

func (c *CExImageButton) draw() {
    panic("Draw not implemented yet on CExImageButton")
}
