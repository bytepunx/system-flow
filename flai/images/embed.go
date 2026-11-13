// Package images carries the Dockerfiles of the harness images in the flai
// binary, so that an installed flai builds them on the host when it cannot
// pull them (ADR-0295 §1). go:embed reaches no parent folder, so the package
// sits beside the Dockerfiles it embeds.
package images

import (
	_ "embed"
	"fmt"
)

// Bedrock is the image every harness image builds over: flai, pilot, git, and
// tini.
const Bedrock = "bedrock"

//go:embed bedrock/Dockerfile
var bedrock []byte

//go:embed claude-code/Dockerfile
var claudeCode []byte

// embedded is each Dockerfile this binary carries, by image: bedrock, then
// each inner harness's.
var embedded = map[string][]byte{
	Bedrock:       bedrock,
	"claude-code": claudeCode,
}

// Dockerfile is the Dockerfile this binary carries for image, bedrock or an
// inner harness, such as claude-code, as flai/images/<image>/Dockerfile holds
// it; an image it carries none for is an error naming those it does.
func Dockerfile(image string) ([]byte, error) {
	data, ok := embedded[image]
	if !ok {
		return nil, fmt.Errorf("flai carries no Dockerfile for the agent image %q: it carries bedrock's and claude-code's; name one of them, or the operator's own image with flai serve agent harness container --image <inner>=<image>", image)
	}
	return data, nil
}

// Path is where image's Dockerfile is in a tree of this repository, from its
// root, in slash form: flai/images/<image>/Dockerfile.
func Path(image string) string {
	return "flai/images/" + image + "/Dockerfile"
}
