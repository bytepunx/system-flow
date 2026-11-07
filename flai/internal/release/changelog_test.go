package release

import (
	"strings"
	"testing"
)

func TestSettleChangelog(t *testing.T) {
	const base = "# Changelog\n\n## Unreleased\n\n- a\n\n## 1.1.14 - d\n\n- x\n"
	tests := []struct {
		name               string
		main, base, branch string
		want               string
	}{{
		// I-0251's first instance: release 1.1.17 named the section the
		// story's line was under; union would bring the old heading back.
		name:   "a release named the section the line was under",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.14 - d\n\n- x\n",
		base:   base,
		branch: "# Changelog\n\n## Unreleased\n\n- a\n- mine\n\n## 1.1.14 - d\n\n- x\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.14 - d\n\n- x\n",
	}, {
		// I-0251's second instance: the branch's line already sat under the
		// top Unreleased, which release 1.1.18 named with S-0422's lines.
		name:   "another release over a settled branch",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.18 - f\n\n- other\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.3 - c\n\n- y\n\n## 1.1.3 - c\n\n- z\n",
		base:   "# Changelog\n\n## Unreleased\n\n- other\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.3 - c\n\n- y\n\n## 1.1.3 - c\n\n- z\n",
		branch: "# Changelog\n\n## Unreleased\n\n- other\n- mine\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.3 - c\n\n- y\n\n## 1.1.3 - c\n\n- z\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.18 - f\n\n- other\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.3 - c\n\n- y\n\n## 1.1.3 - c\n\n- z\n",
	}, {
		name:   "after another story's line under the top Unreleased",
		main:   "# Changelog\n\n## Unreleased\n\n- a\n- theirs\n\n## 1.1.14 - d\n\n- x\n",
		base:   base,
		branch: "# Changelog\n\n## Unreleased\n\n- a\n- mine\n- mine too\n\n## 1.1.14 - d\n\n- x\n",
		want:   "# Changelog\n\n## Unreleased\n\n- a\n- theirs\n- mine\n- mine too\n\n## 1.1.14 - d\n\n- x\n",
	}, {
		name:   "main has no Unreleased",
		main:   "# Changelog\n\n## 1.1.17 - e\n\n- a\n",
		base:   "# Changelog\n\n## 1.1.14 - d\n\n- x\n",
		branch: "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.14 - d\n\n- x\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n",
	}, {
		name:   "main has no section at all",
		main:   "# Changelog\n",
		base:   "# Changelog\n",
		branch: "# Changelog\n\n## Unreleased\n\n- mine\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n",
	}, {
		name:   "the branch added no line",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n\n- a\n",
		base:   base,
		branch: base,
		want:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n\n- a\n",
	}, {
		name:   "the line is already under main's Unreleased",
		main:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n",
		base:   base,
		branch: "# Changelog\n\n## Unreleased\n\n- a\n- mine\n\n## 1.1.14 - d\n\n- x\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n",
	}, {
		name:   "a line the branch removed is not carried",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.14 - d\n\n- x\n",
		base:   base,
		branch: "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.14 - d\n\n- x\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n\n## 1.1.14 - d\n\n- x\n",
	}, {
		name:   "a line the branch wrote under a version heading",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n- a\n",
		base:   "# Changelog\n\n## 1.1.17 - e\n\n- a\n",
		branch: "# Changelog\n\n## 1.1.17 - e\n\n- a\n- mine\n",
		want:   "# Changelog\n\n## Unreleased\n\n- mine\n\n## 1.1.17 - e\n\n- a\n",
	}, {
		name:   "a line added twice is carried twice",
		main:   "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n- a\n",
		base:   "# Changelog\n\n## Unreleased\n",
		branch: "# Changelog\n\n## Unreleased\n\n- b\n- b\n",
		want:   "# Changelog\n\n## Unreleased\n\n- b\n- b\n\n## 1.1.17 - e\n\n- a\n",
	}, {
		name:   "main ends without a newline",
		main:   "# Changelog\n\n## Unreleased",
		base:   "# Changelog\n\n## Unreleased",
		branch: "# Changelog\n\n## Unreleased\n\n- mine",
		want:   "# Changelog\n\n## Unreleased\n\n- mine",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettleChangelog(tt.main, tt.base, tt.branch)
			if got != tt.want {
				t.Fatalf("SettleChangelog:\n got %q\nwant %q", got, tt.want)
			}
			if strings.Count(got, "\n"+unreleasedHeading+"\n") > 1 {
				t.Errorf("more than one %q in %q", unreleasedHeading, got)
			}
			if again := SettleChangelog(got, tt.base, tt.branch); again != got {
				t.Errorf("settling again changed it:\n got %q\nwant %q", again, got)
			}
		})
	}
}

// TestSettleChangelogIsWhatAReleaseNamesNext checks that the lines settled
// under the top Unreleased are what the next release names with its version.
func TestSettleChangelogIsWhatAReleaseNamesNext(t *testing.T) {
	main := "# Changelog\n\n## Unreleased\n\n## 1.1.17 - e\n\n- a\n"
	got := SettleChangelog(main, "# Changelog\n\n## Unreleased\n\n- a\n", "# Changelog\n\n## Unreleased\n\n- a\n- mine\n")
	rest, body, ok := cutUnreleased(got)
	if !ok || body != "- mine\n" {
		t.Fatalf("cutUnreleased(%q) = %q, %v; want the body - mine", got, body, ok)
	}
	if strings.Contains(rest, "mine") {
		t.Errorf("the line is outside the Unreleased section too: %q", rest)
	}
}
