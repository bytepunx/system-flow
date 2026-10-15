package designsystem

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

func oklch(l, c, h float64) Color {
	return Color{Space: ColorSpaceOKLCH, Components: [3]float64{l, c, h}, Alpha: 1}
}

// rgbOf is the colour a #rrggbb gives.
func rgbOf(t *testing.T, hex string) RGB {
	t.Helper()
	var ch [3]float64
	for i := range ch {
		b, err := strconv.ParseUint(hex[1+2*i:3+2*i], 16, 8)
		if err != nil {
			t.Fatalf("rgbOf(%q): %v", hex, err)
		}
		ch[i] = float64(b) / 255
	}
	return RGB{ch[0], ch[1], ch[2]}
}

func near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestSRGBConversion(t *testing.T) {
	cases := []struct {
		name    string
		c       Color
		hex     string
		clipped bool
	}{
		{"srgb is its components", Color{Space: ColorSpaceSRGB, Components: [3]float64{0.96, 0.96, 0.97}, Alpha: 1}, "#f5f5f7", false},
		{"CSS Color 4's oklch for red", oklch(0.628, 0.2577, 29.23), "#ff0000", false},
		{"oklch white", oklch(1, 0, 0), "#ffffff", false},
		{"oklch black", oklch(0, 0, 0), "#000000", false},
		{"an achromatic oklch with no hue", Color{Space: ColorSpaceOKLCH, Components: [3]float64{0.5, 0, 0}, NoHue: true, Alpha: 1}, "#636363", false},
		{"an oklch green inside sRGB", oklch(0.5, 0.13, 150), "#137738", false},
		{"an oklch blue just outside sRGB", oklch(0.52, 0.18, 255), "#0066cd", true},
		{"an oklch green far outside sRGB", oklch(0.9, 0.4, 145), "#00ff00", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rgb, clipped := c.c.SRGB()
			if got := rgb.Hex(); got != c.hex || clipped != c.clipped {
				t.Errorf("SRGB() = %s, clipped %v; want %s, clipped %v", got, clipped, c.hex, c.clipped)
			}
			for _, v := range []float64{rgb.R, rgb.G, rgb.B} {
				if v < 0 || v > 1 {
					t.Errorf("SRGB() = %+v, a channel outside 0 to 1", rgb)
				}
			}
		})
	}
}

func TestOklchToSRGBChannels(t *testing.T) {
	// oklch(0.5 0.13 150), in gamut: its unclamped channels by Ottosson's
	// matrices, computed once with them in JavaScript.
	rgb, _ := oklch(0.5, 0.13, 150).SRGB()
	want := RGB{0.075080, 0.466092, 0.219843}
	if !near(rgb.R, want.R, 1e-6) || !near(rgb.G, want.G, 1e-6) || !near(rgb.B, want.B, 1e-6) {
		t.Errorf("SRGB() = %+v, want %+v", rgb, want)
	}
}

func TestLuminanceAndRatio(t *testing.T) {
	white, black := RGB{1, 1, 1}, RGB{0, 0, 0}
	if l := white.Luminance(); !near(l, 1, 1e-12) {
		t.Errorf("white's luminance = %v, want 1", l)
	}
	if l := black.Luminance(); l != 0 {
		t.Errorf("black's luminance = %v, want 0", l)
	}
	if r := ContrastRatio(white, black); !near(r, 21, 1e-9) {
		t.Errorf("white on black = %v, want 21", r)
	}
	if a, b := ContrastRatio(black, white), ContrastRatio(white, black); a != b {
		t.Errorf("ContrastRatio(black, white) = %v, ContrastRatio(white, black) = %v; want the same", a, b)
	}
}

// TestRatiosAgreeWithFlaiover checks the ratio of pairs of the dashboard's
// own colours against what flaiover's contrast() in src/lib/theme.test.ts
// gives: the expected ratios were computed by running that function's
// arithmetic, unchanged, in node on these hex values, which are copied from
// flaiover/src/routes/layout.css (--t-<name>, :root for light and
// [data-theme='dark'] for dark) and from theme.test.ts's own brand checks.
func TestRatiosAgreeWithFlaiover(t *testing.T) {
	cases := []struct {
		name   string
		fg, bg string
		want   float64
	}{
		{"light ink on ground", "#1c1917", "#f7f3e3", 15.7264153227},
		{"light muted on ground", "#645853", "#f7f3e3", 6.1659706554},
		{"light muted on surface", "#645853", "#fdfbf3", 6.6183014471},
		{"light on-primary on primary", "#ffffff", "#054a91", 8.7688841343},
		{"light good on good-soft", "#0c6e19", "#d1fad6", 5.6306980707},
		{"light danger on danger-soft", "#9b1c1c", "#f9e1dc", 6.5324889947},
		{"light line-strong on ground", "#8b7b74", "#f7f3e3", 3.6433854475},
		{"dark ink on ground", "#f7f3e3", "#191615", 16.1829136955},
		{"dark muted on surface", "#a79a95", "#272220", 5.7632478065},
		{"dark good on good-soft", "#62c96f", "#084910", 5.1265509268},
		{"dark danger on danger-soft", "#f28b82", "#4a1a17", 6.0572189020},
		{"dark line-strong on ground", "#7b6c66", "#191615", 3.5812673297},
		{"the brand green on the light ground", "#119822", "#f7f3e3", 3.4098055814},
		{"the brand cyan on the light ground", "#23b5d3", "#f7f3e3", 2.1903337956},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ContrastRatio(rgbOf(t, c.fg), rgbOf(t, c.bg)); !near(got, c.want, 1e-9) {
				t.Errorf("ContrastRatio(%s, %s) = %.10f, want flaiover's %.10f", c.fg, c.bg, got, c.want)
			}
		})
	}
}

func TestContrastOfAFullTheme(t *testing.T) {
	sys := load(t, fullFile(t))
	got := Contrast(ThemeLight, tokens(t, sys, ThemeLight))
	// The colours by Ottosson's matrices and the ratios by flaiover's
	// arithmetic, computed once in node from the fixture's values.
	want := []struct {
		fg, bg          string
		fgHex, bgHex    string
		ratio, thresh   float64
		passes, clipped bool
	}{
		{"text", "background", "#141822", "#ffffff", 17.74, 4.5, true, false},
		{"text", "surface", "#141822", "#f5f5f7", 16.29, 4.5, true, false},
		{"text-muted", "background", "#505561", "#ffffff", 7.46, 4.5, true, false},
		{"text-muted", "surface", "#505561", "#f5f5f7", 6.85, 4.5, true, false},
		{"success", "background", "#137738", "#ffffff", 5.63, 4.5, true, false},
		{"success", "surface", "#137738", "#f5f5f7", 5.17, 4.5, true, false},
		{"success", "success-soft", "#137738", "#dcf7e1", 4.95, 4.5, true, false},
		{"warning", "background", "#9d6300", "#ffffff", 4.97, 4.5, true, true},
		{"warning", "surface", "#9d6300", "#f5f5f7", 4.57, 4.5, true, true},
		{"warning", "warning-soft", "#9d6300", "#fff1cc", 4.43, 4.5, false, true},
		{"danger", "background", "#be222a", "#ffffff", 6.08, 4.5, true, false},
		{"danger", "surface", "#be222a", "#f5f5f7", 5.58, 4.5, true, false},
		{"danger", "danger-soft", "#be222a", "#ffe7e4", 5.15, 4.5, true, true},
		{"info", "background", "#006aa0", "#ffffff", 5.88, 4.5, true, true},
		{"info", "surface", "#006aa0", "#f5f5f7", 5.4, 4.5, true, true},
		{"info", "info-soft", "#006aa0", "#ddf2ff", 5.1, 4.5, true, true},
		{"on-primary", "primary", "#ffffff", "#0066cd", 5.55, 4.5, true, true},
		{"on-secondary", "secondary", "#ffffff", "#315e60", 7.23, 4.5, true, false},
		{"border", "background", "#ccccd1", "#ffffff", 1.59, 3, false, false},
		{"focus", "background", "#007df3", "#ffffff", 4.01, 3, true, true},
	}
	if len(got) != len(want) {
		t.Fatalf("Contrast gave %d pairs, want %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		wp := Pair{
			Theme: ThemeLight, Foreground: "color." + w.fg, Background: "color." + w.bg,
			ForegroundHex: w.fgHex, BackgroundHex: w.bgHex,
			Ratio: w.ratio, Threshold: w.thresh, Passes: w.passes, Clipped: w.clipped,
		}
		if got[i] != wp {
			t.Errorf("pair %d = %+v\nwant     %+v", i, got[i], wp)
		}
	}
}

func TestContrastOfEachTheme(t *testing.T) {
	dark := `{ "color": { "background": { "$value": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000" } } } }`
	sys := load(t, resolverDoc(t, dark))
	if len(sys.Problems) != 0 {
		t.Fatalf("problems: %+v", sys.Problems)
	}
	p := Contrast(ThemeDark, tokens(t, sys, ThemeDark))[0]
	if p.Theme != ThemeDark || p.Foreground != "color.text" || p.BackgroundHex != "#000000" || p.Ratio != 1.18 || p.Passes {
		t.Errorf("dark's first pair = %+v, want color.text on the dark background #000000 at 1.18, failing", p)
	}
}

// colorTokens is a token set of srgb colour roles, each named with its hex.
func colorTokens(t *testing.T, roles ...string) *TokenSet {
	t.Helper()
	var members []string
	for i := 0; i < len(roles); i += 2 {
		c := rgbOf(t, roles[i+1])
		members = append(members, fmt.Sprintf(`%q: { "$value": { "colorSpace": "srgb", "components": [%v, %v, %v], "hex": %q } }`,
			roles[i], c.R, c.G, c.B, roles[i+1]))
	}
	root, err := ParseJSON([]byte(`{ "color": { "$type": "color", ` + strings.Join(members, ", ") + ` } }`))
	if err != nil {
		t.Fatal(err)
	}
	return ReadTokens(root)
}

func TestContrastAtTheThresholds(t *testing.T) {
	// On white, #767676 is 4.542:1, #777777 4.478:1, #949494 3.033:1, and
	// #959595 2.995:1, either side of 4.5:1 and of 3:1.
	s := colorTokens(t,
		"background", "#ffffff",
		"text", "#767676",
		"text-muted", "#777777",
		"border", "#949494",
		"focus", "#959595",
	)
	got := Contrast(ThemeLight, s)
	want := []struct {
		fg     string
		ratio  float64
		passes bool
	}{
		{"color.text", 4.54, true},
		{"color.text-muted", 4.47, false},
		{"color.border", 3.03, true},
		{"color.focus", 2.99, false},
	}
	if len(got) != len(want) {
		t.Fatalf("Contrast = %+v, want the %d pairs whose roles are given", got, len(want))
	}
	for i, w := range want {
		if got[i].Foreground != w.fg || got[i].Background != "color.background" || got[i].Ratio != w.ratio || got[i].Passes != w.passes {
			t.Errorf("pair %d = %+v, want %s on color.background at %v, passes %v", i, got[i], w.fg, w.ratio, w.passes)
		}
	}
}

func TestContrastLeavesOutAPairItCannotRead(t *testing.T) {
	cases := []struct {
		name string
		edit func(f *File)
		gone string // the role whose pairs are left out
		left int
	}{
		{"a missing role", func(f *File) { without(f.Root, "color", "surface") }, "color.surface", 14},
		{"a colour flai does not read", func(f *File) {
			f.Root.Get("color").Get("info").Get("$value").Get("colorSpace").Text = "hsl"
		}, "color.info", 17},
		{"an alias that does not resolve", func(f *File) {
			f.Root.Get("color").Get("border").Members[0].Value = &Node{Kind: KindString, Text: "{color.nope}"}
		}, "color.border", 19},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := parsePlain(t, fullFile(t))
			c.edit(f)
			got := Contrast(ThemeLight, Resolve(&File{Root: f.Root, Tokens: ReadTokens(f.Root)}).Themes[0].Tokens)
			if len(got) != c.left {
				t.Errorf("Contrast gave %d pairs, want %d", len(got), c.left)
			}
			for _, p := range got {
				if p.Foreground == c.gone || p.Background == c.gone {
					t.Errorf("Contrast kept %+v, which needs %s", p, c.gone)
				}
			}
		})
	}
}
