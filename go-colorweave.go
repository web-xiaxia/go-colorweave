package main

import (
	gwc "github.com/jyotiska/go-webcolors"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"sort"
)

const (
	ColorModelCss3  = "css3"
	ColorModelCss21 = "css21"
	ColorThemeLight = "light"
	ColorThemeDark  = "dark"
)

type ColorInfo struct {
	Name       string
	Hex        string
	Counter    int
	Proportion float64
	Theme      string
}
type ColorInfoList []*ColorInfo

func (c ColorInfoList) Theme() string {
	themeLightProportion := float64(0)
	themeDarkProportion := float64(0)
	for _, info := range c {
		if info.Theme == ColorThemeLight {
			themeLightProportion += info.Proportion
		} else {
			themeDarkProportion += info.Proportion
		}
	}
	if themeLightProportion > themeDarkProportion {
		return ColorThemeLight
	}
	return ColorThemeDark
}

func HexToHSL(hex string) (float64, float64, float64) {
	rgb := gwc.HexToRGB(hex)
	rInt, gInt, bInt := rgb[0], rgb[1], rgb[2]
	r := float64(rInt) / 255
	g := float64(gInt) / 255
	b := float64(bInt) / 255
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	var h, s, l float64
	l = (max + min) / 2
	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		s = d / (2 - max - min)
		if l <= 0.5 {
			s = d / (max + min)
		}
		switch max {
		case r:
			e := float64(0)
			if g < b {
				e = 6
			}
			h = (g-b)/d + e
			break
		case g:
			h = (b-r)/d + 2
			break
		case b:
			h = (r-g)/d + 4
			break
		}
		h /= 6
	}
	return h, s, l
}

func IsColorDarkOrLight(hex string) string {
	_, _, l := HexToHSL(hex)
	if l > 0.5 {
		return ColorThemeLight
	}
	return ColorThemeDark
}

// FindClosestColor This method finds the closest color for a given RGB tuple and returns the name of the color in given mode
func FindClosestColor(requestedColor []int, mode string) (string, string) {
	minColors := make(map[int]string)
	var colorMap map[string]string

	// css3 gives the shades while css21 gives the primary or base colors
	if mode == ColorModelCss3 {
		colorMap = gwc.CSS3NamesToHex
	} else {
		colorMap = gwc.HTML4NamesToHex
	}

	for name, hexCode := range colorMap {
		rgbTriplet := gwc.HexToRGB(hexCode)
		rd := math.Pow(float64(rgbTriplet[0]-requestedColor[0]), float64(2))
		gd := math.Pow(float64(rgbTriplet[1]-requestedColor[1]), float64(2))
		bd := math.Pow(float64(rgbTriplet[2]-requestedColor[2]), float64(2))
		minColors[int(rd+gd+bd)] = name
	}

	keys := make([]int, 0, len(minColors))
	for key := range minColors {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return minColors[keys[0]], colorMap[minColors[keys[0]]]
}

func ListDominantColors(image image.Image, limit int) ColorInfoList {
	return ListDominantColorsWithMode(image, limit, ColorModelCss21)
}
func ListDominantColorsWithMode(image image.Image, limit int, colorModel string) ColorInfoList {
	// Resize the image to smaller scale for faster computation
	image = resize.Resize(100, 0, image, resize.Bilinear)
	bounds := image.Bounds()
	totalPixels := bounds.Max.X * bounds.Max.Y

	colorInfoList := make([]*ColorInfo, 0, 6)
	colorInfoMap := make(map[string]*ColorInfo, 6)
	for i := 0; i <= bounds.Max.X; i++ {
		for j := 0; j <= bounds.Max.Y; j++ {
			pixel := image.At(i, j)
			red, green, blue, _ := pixel.RGBA()
			rgbTuple := []int{int(red / 255), int(green / 255), int(blue / 255)}
			colorName, colorHex := FindClosestColor(rgbTuple, colorModel)

			colorInfo, present := colorInfoMap[colorName]
			if !present {
				colorInfo = &ColorInfo{
					Name:       colorName,
					Hex:        colorHex,
					Counter:    0,
					Proportion: 0,
					Theme:      IsColorDarkOrLight(colorHex),
				}
				colorInfoList = append(colorInfoList, colorInfo)
				colorInfoMap[colorName] = colorInfo
			}
			colorInfo.Counter += 1
		}
	}

	sort.Slice(colorInfoList, func(i, j int) bool {
		return colorInfoList[i].Counter > colorInfoList[j].Counter
	})

	if len(colorInfoList) > limit {
		colorInfoList = colorInfoList[:limit]
	}

	// Display the top N dominant colors from the image
	for _, val := range colorInfoList {
		val.Proportion = (float64(val.Counter) / float64(totalPixels)) * 100
		// fmt.Printf("%s %s %.2f%%\n", val.Name, val.Hex, val.Proportion)
	}
	return colorInfoList
}
