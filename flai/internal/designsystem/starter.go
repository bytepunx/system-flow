package designsystem

import _ "embed"

// StarterTokens is the token file flai design init writes: a resolver
// document with every required token in light and dark, whose role pairs
// pass their contrast thresholds in both themes.
//
//go:embed starter.tokens.json
var StarterTokens []byte

// StarterGuidance is the guidance flai design init writes, before it fills
// the front matter's updated and topics and the token file's path, which the
// body gives as {{tokens}}.
//
//go:embed starter.md
var StarterGuidance string

// StarterTokensMark is the mark in StarterGuidance that stands for the token
// file's path.
const StarterTokensMark = "{{tokens}}"
