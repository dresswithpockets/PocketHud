package main

import (
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "image"
    "os"
)

type PictureProvider interface {
    getPicture(path string) (pixel.Picture, error)
}

type VguiProvider interface {
    getObject(path string) (*vgui.Object, error)
}

type SourcePictureProvider struct {
    pictures           map[string]pixel.Picture
    fileSourceProvider vgui.FileSourceProvider
}

type SourceVguiProvider struct {
    roots map[string]*vgui.Object
    fileSourceProvider vgui.FileSourceProvider
}

// getPicture returns a pixel.Picture from the path, and caches it in pictures.
// If the pixel.Picture is not already loaded & mapped to the path, it will load it according to fileSourceProvider
func (p *SourcePictureProvider) getPicture(path string) (pixel.Picture, error) {
    abs, err := p.fileSourceProvider.GetAbsolute(path)
    if err != nil {
        return nil, err
    }
    if pic, ok := p.pictures[abs]; ok {
        return pic, nil
    }
    pic, err := loadPicture(abs)
    if err != nil {
        return nil, err
    }
    p.pictures[abs] = pic
    return pic, nil
}

func (p * SourceVguiProvider) getObject(path string) (*vgui.Object, error) {
    return vgui.FromFileSourceProvider(path, p.fileSourceProvider)
}

func loadPicture(path string) (pixel.Picture, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }
    return pixel.PictureDataFromImage(img), nil
}
