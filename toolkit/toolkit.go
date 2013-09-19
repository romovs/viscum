package toolkit

import (
	"image"
	"bufio"
	_ "image/png"
	"os"
	log "github.com/cihub/seelog"
)

type Toolkit struct {

	// checkbox. should be of same size.
	cbChecked			image.Image
	cbUnchecked			image.Image
	cbCheckedHover		image.Image
	cbUncheckedHover	image.Image
}

const(
	THEME_PATH	= "data/elements/"
)

var tk Toolkit

func init() {

	tk = Toolkit{}

	tk.cbChecked, _ = loadImage(THEME_PATH + "checkbox-checked.png")
	tk.cbUnchecked, _ = loadImage(THEME_PATH + "checkbox-unchecked.png")
	tk.cbCheckedHover, _ = loadImage(THEME_PATH + "checkbox-checked-hover.png")
	tk.cbUncheckedHover, _ = loadImage(THEME_PATH + "checkbox-unchecked-hover.png")
}

func loadImage(filename string) (image.Image, error) {

	f, err := os.Open(filename)
    if err != nil {
		log.Critical(err)
		return nil, err
    }
    defer f.Close()
    		
    img, _, err := image.Decode(bufio.NewReader(f))
    if err != nil  {
		log.Critical(err)
		return nil, err
	}

	return img, nil
}