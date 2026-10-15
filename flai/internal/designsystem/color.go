package designsystem

import (
	"fmt"
	"math"
)

// RGB is a gamma-encoded sRGB colour, each channel 0 to 1.
type RGB struct {
	R, G, B float64
}

// clipTolerance is how far outside 0 to 1 an encoded channel may fall and
// still count as in gamut: half of one 8-bit step, less than its hex can
// show. oklch values written to a few decimals, such as CSS Color 4's
// oklch(0.628 0.2577 29.23) for red, land a hair outside sRGB, and calling
// them clipped would be noise.
const clipTolerance = 0.5 / 255

// SRGB converts the colour to sRGB and reports whether it had to be clipped
// to the sRGB gamut. Alpha is left out.
//
// An srgb colour is its components. An oklch colour goes oklch → OKLab →
// linear sRGB by Björn Ottosson's matrices, the ones CSS Color 4 uses, then
// through the sRGB transfer function. A channel outside 0 to 1 is clamped to
// it. CSS Color 4's gamut mapping reduces chroma at constant lightness and
// hue instead, which keeps the colour's look closer; clamping is enough here,
// where the converted colour feeds only the contrast report, which says the
// colour was clipped, and the file keeps its oklch.
func (c Color) SRGB() (RGB, bool) {
	if c.Space == ColorSpaceSRGB {
		return RGB{c.Components[0], c.Components[1], c.Components[2]}, false
	}
	l, a, b := oklab(c)
	lr, lg, lb := oklabToLinearSRGB(l, a, b)
	rgb := RGB{encode(lr), encode(lg), encode(lb)}
	clipped := false
	for _, v := range []*float64{&rgb.R, &rgb.G, &rgb.B} {
		if *v < -clipTolerance || *v > 1+clipTolerance {
			clipped = true
		}
		*v = min(max(*v, 0), 1)
	}
	return rgb, clipped
}

// oklab is an oklch colour in OKLab: lightness and its a and b axes.
func oklab(c Color) (l, a, b float64) {
	l, ch, h := c.Components[0], c.Components[1], c.Components[2]
	if c.NoHue {
		h = 0
	}
	rad := h * math.Pi / 180
	return l, ch * math.Cos(rad), ch * math.Sin(rad)
}

// oklabToLinearSRGB is Ottosson's OKLab to linear sRGB, through LMS.
func oklabToLinearSRGB(l, a, b float64) (r, g, bl float64) {
	l1 := l + 0.3963377774*a + 0.2158037573*b
	m1 := l - 0.1055613458*a - 0.0638541728*b
	s1 := l - 0.0894841775*a - 1.2914855480*b
	lc, mc, sc := l1*l1*l1, m1*m1*m1, s1*s1*s1
	return +4.0767416621*lc - 3.3077115913*mc + 0.2309699292*sc,
		-1.2684380046*lc + 2.6097574011*mc - 0.3413193965*sc,
		-0.0041960863*lc - 0.7034186147*mc + 1.7076147010*sc
}

// encode is the sRGB transfer function, extended to negative values by
// symmetry as CSS Color 4 extends it, so that clipping sees how far out a
// channel falls.
func encode(v float64) float64 {
	x := math.Abs(v)
	e := 12.92 * x
	if x > 0.0031308 {
		e = 1.055*math.Pow(x, 1/2.4) - 0.055
	}
	return math.Copysign(e, v)
}

// Hex is the colour as a six-digit #rrggbb, lower case, each channel rounded
// to 8 bits.
func (r RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", byte8(r.R), byte8(r.G), byte8(r.B))
}

// byte8 is a 0 to 1 channel as the nearest 8-bit value.
func byte8(v float64) int { return int(math.Round(min(max(v, 0), 1) * 255)) }

// rounded is the colour its Hex gives.
func (r RGB) rounded() RGB {
	return RGB{float64(byte8(r.R)) / 255, float64(byte8(r.G)) / 255, float64(byte8(r.B)) / 255}
}

// Luminance is the colour's WCAG 2 relative luminance, 0 for black to 1 for
// white.
//
// Each channel is linearised with the threshold 0.03928 that WCAG 2's
// definition and flaiover's own contrast test (src/lib/theme.test.ts) use;
// the sRGB standard's 0.04045 gives the same value for every 8-bit channel,
// since none falls between the two.
func (r RGB) Luminance() float64 {
	lin := func(v float64) float64 {
		if v <= 0.03928 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(r.R) + 0.7152*lin(r.G) + 0.0722*lin(r.B)
}

// ContrastRatio is the WCAG 2 contrast ratio of two colours, 1 to 21, the
// same whichever is given first.
func ContrastRatio(a, b RGB) float64 {
	hi, lo := a.Luminance(), b.Luminance()
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

// The contrast thresholds of ADR-0148, the ones flaiover's own test uses.
const (
	ThresholdText    = 4.5 // text, and a state's colour read as text
	ThresholdControl = 3.0 // a border or a focus ring against what it sits on
)

// Pair is one role pair's contrast in one theme (ADR-0148 § Decision 4).
type Pair struct {
	Theme         string  `json:"theme"`          // ThemeLight or ThemeDark
	Foreground    string  `json:"foreground"`     // token path, such as color.text
	Background    string  `json:"background"`     // token path, such as color.background
	ForegroundHex string  `json:"foreground_hex"` // the foreground in sRGB, #rrggbb
	BackgroundHex string  `json:"background_hex"` // the background in sRGB, #rrggbb
	Ratio         float64 `json:"ratio"`          // the contrast ratio, truncated to two decimals so that it never reads as passing when it fails
	Threshold     float64 `json:"threshold"`      // the least ratio that passes: ThresholdText or ThresholdControl
	Passes        bool    `json:"passes"`         // the ratio, before truncation, is at least the threshold
	Clipped       bool    `json:"clipped"`        // either colour is an oklch colour outside sRGB, clipped to it
}

// rolePair is a foreground role, the role it sits on, and the threshold.
type rolePair struct {
	fg, bg    string
	threshold float64
}

// rolePairs are the pairs ADR-0148 § Decision 4 names, in its order: "text,
// muted text, and each state on `background`, `surface`, or its tint at
// 4.5:1; `on-primary` and `on-secondary` on their fill at 4.5:1; `border` and
// `focus` on `background` at 3:1." A tint is a state's own -soft role, so
// text and muted text are paired with background and surface only.
var rolePairs = func() []rolePair {
	var out []rolePair
	for _, fg := range []string{"text", "text-muted"} {
		out = append(out, rolePair{fg, "background", ThresholdText}, rolePair{fg, "surface", ThresholdText})
	}
	for _, state := range []string{"success", "warning", "danger", "info"} {
		out = append(out,
			rolePair{state, "background", ThresholdText},
			rolePair{state, "surface", ThresholdText},
			rolePair{state, state + "-soft", ThresholdText})
	}
	return append(out,
		rolePair{"on-primary", "primary", ThresholdText},
		rolePair{"on-secondary", "secondary", ThresholdText},
		rolePair{"border", "background", ThresholdControl},
		rolePair{"focus", "background", ThresholdControl})
}()

// Contrast reports the contrast of each role pair ADR-0148 names, in its
// order, for one theme's tokens. Each ratio is of the colours as their hex
// gives them, as flaiover's test computes it, so that it can be checked from
// the hex the pair reports. A pair whose role is missing, or not a colour
// flai reads, is left out; Validate reports those. Contrast is a report,
// never a problem.
func Contrast(theme string, tokens *TokenSet) []Pair {
	var out []Pair
	for _, rp := range rolePairs {
		fg, fgClipped, ok := roleColor(tokens, rp.fg)
		if !ok {
			continue
		}
		bg, bgClipped, ok := roleColor(tokens, rp.bg)
		if !ok {
			continue
		}
		ratio := ContrastRatio(fg.rounded(), bg.rounded())
		out = append(out, Pair{
			Theme:         theme,
			Foreground:    GroupColor + "." + rp.fg,
			Background:    GroupColor + "." + rp.bg,
			ForegroundHex: fg.Hex(),
			BackgroundHex: bg.Hex(),
			Ratio:         math.Floor(ratio*100) / 100,
			Threshold:     rp.threshold,
			Passes:        ratio >= rp.threshold,
			Clipped:       fgClipped || bgClipped,
		})
	}
	return out
}

// roleColor is the colour role name in tokens, in sRGB, and whether it was
// clipped; false when the role is missing or not a colour flai reads.
func roleColor(tokens *TokenSet, name string) (RGB, bool, bool) {
	t, ok := tokens.Get(GroupColor + "." + name)
	if !ok || t.Value == nil {
		return RGB{}, false, false
	}
	c, err := DecodeColor(t.Value)
	if err != nil {
		return RGB{}, false, false
	}
	rgb, clipped := c.SRGB()
	return rgb, clipped, true
}
