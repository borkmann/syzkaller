// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package report

import (
	"bytes"
	"regexp"

	"github.com/google/syzkaller/pkg/symbolizer"
)

type fuchsia struct {
	kernelSrc string
	kernelObj string
	symbols   map[string][]symbolizer.Symbol
	ignores   []*regexp.Regexp
}

func ctorFuchsia(kernelSrc, kernelObj string, symbols map[string][]symbolizer.Symbol,
	ignores []*regexp.Regexp) (Reporter, error) {
	ctx := &fuchsia{
		kernelSrc: kernelSrc,
		kernelObj: kernelObj,
		symbols:   symbols,
		ignores:   ignores,
	}
	return ctx, nil
}

func (ctx *fuchsia) ContainsCrash(output []byte) bool {
	return bytes.Contains(output, []byte("ZIRCON KERNEL PANIC")) ||
		bytes.Contains(output, []byte("Supervisor Page Fault"))
}

func (ctx *fuchsia) Parse(output []byte) *Report {
	title, pos := "", 0
	if pos = bytes.Index(output, []byte("ZIRCON KERNEL PANIC")); pos != -1 {
		title = "ZIRCON KERNEL PANIC"
	} else if pos = bytes.Index(output, []byte("Supervisor Page Fault")); pos != -1 {
		title = "Supervisor Page Fault"
	} else {
		return nil
	}
	return &Report{
		Title:    title,
		Report:   output,
		Output:   output,
		StartPos: pos,
		EndPos:   pos + len(title),
	}
}

func (ctx *fuchsia) Symbolize(rep *Report) error {
	return nil
}
