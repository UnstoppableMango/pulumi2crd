package pulumi2crd

import (
	"fmt"
	"strings"
)

func DecomposeToken(tok string) (pack string, mod string, typ string, err error) {
	// https://github.com/pulumi/pulumi/blob/b0d15812cba1cd74b5441f20e3345ae63778a308/pkg/codegen/pcl/utilities.go#L51
	if components := strings.Split(tok, ":"); len(components) != 3 {
		return "", "", tok, fmt.Errorf("malformed token: %s", tok)
	} else {
		return components[0], components[1], components[2], nil
	}
}
