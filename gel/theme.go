package gel

import (
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/monorepo/opts/binary"
	"github.com/p9c/monorepo/qu"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Theme struct {
	quit       qu.C
	shaper     text.Shaper
	collection Collection
	TextSize   unit.Value
	*Colors
	icons         map[string]*Icon
	scrollBarSize int
	Dark          *binary.Opt
	iconCache     IconCache
	WidgetPool    *Pool
}

// NewTheme creates a new theme to use for rendering a user interface
func NewTheme(dark *binary.Opt, fontCollection []text.FontFace, quit qu.C) (th *Theme) {
	textSize := unit.Sp(16)
	if runtime.GOOS == "linux" {
		var e error
		var b []byte
		runner := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "text-scaling-factor")
		if b, e = runner.CombinedOutput(); D.Chk(e) {
		}
		var factor float64
		numberString := strings.TrimSpace(string(b))
		if factor, e = strconv.ParseFloat(numberString, 10); D.Chk(e) {
		}
		I.Ln(factor)
		textSize = textSize.Scale(float32(factor))
	}
	th = &Theme{
		quit:          quit,
		shaper:        text.NewCache(fontCollection),
		collection:    fontCollection,
		TextSize:      textSize,
		Colors:        newColors(),
		scrollBarSize: 0,
		iconCache:     make(IconCache),
	}
	th.SetDarkTheme(dark.True())
	return
}
