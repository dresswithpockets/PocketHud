package main

// MainMenuOverride defined the area of the screen that the main menu occupies
type MainMenuOverride struct {
    BaseControl
    updateUrl     string `vgui:"update_url"`
    blogUrl       string `vgui:"blog_url"`
    buttonXOffset int    `vgui:"button_x_offset"`
    buttonY       int    `vgui:"button_y"`
    buttonYDelta  int    `vgui:"button_y_delta"`
}

// KvButton child control for MainMenuOverride
// defines the style and layout for the buttons in gamemenu.res that arent explicitly defined
type KvButton struct {
    BaseControl
}

type SaxxySettings struct {
    BaseControl
    autoResize              bool
    pinCorner               bool
    flashBounds             Bounds
    flashStartSizeMin       int
    flashStartSizeMax       int
    flashMaxScale           int
    flashLifeLengthMin      float32
    flashLifeLengthMax      float32
    curtainAnimDuration     float32
    curtainOpenTime         float32
    flashStartTime          float32
    initialFreakoutDuration float32
    clapSoundDuration       float32
    cameraFlashSettings     CameraFlashSettings
}

type CameraFlashSettings struct {
    BaseControl
    tileImage  bool
    scaleImage bool
}

type CPvPRankPanel struct {
    EditablePanel
    // TODO: functionality
}

type CSteamFriendsListPanel struct {
    BaseControl
    // TODO: what base type?
}

type ScrollBar struct {
    BaseControl
    // TODO: what base type?
}

type ScrollableEditablePanel struct {
    EditablePanel
    // TODO: functionality
}

type CMainMenuNotificationsControl struct {
    BaseControl
    // TODO: what base type?
}

type CTFStreamListPanel struct {
    BaseControl
    // TODO: what base type?
}

type CItemModelPanel struct {
    BaseControl
    // TODO: what base type?
}

// TODO: from vgui.Value

func (m *MainMenuOverride) draw() {
    panic("Draw is not implemented for MainMenuOverride yet.")
}

func (k *KvButton) draw() {
    panic("Draw is not implemented for KvButton yet.")
}

func (s *SaxxySettings) draw() {
    panic("Draw is not implemented for SaxxySettings yet.")
}

func (c *CPvPRankPanel) draw() {
    panic("Draw is not implemented for CPvPRankPanel yet.")
}

func (c *CSteamFriendsListPanel) draw() {
    panic("Draw is not implemented for CStreamFriendsListPanel yet.")
}

func (s *ScrollBar) draw() {
    panic("Draw is not implemented for ScrollBar yet.")
}

func (s *ScrollableEditablePanel) draw() {
    panic("Draw is not implemented for ScrollableEditablePanel yet.")
}

func (c *CMainMenuNotificationsControl) draw() {
    panic("Draw is not implemented for CMainMenuNotificationsControl yet.")
}

func (c *CTFStreamListPanel) draw() {
    panic("Draw is not implemented for CTFStreamListPanel yet.")
}

func (c *CItemModelPanel) draw() {
    panic("Draw is not implemented for CItemModelPanel yet.")
}
