package interpret

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sk "github.com/google/skylark"
)

type Interpreter struct {
	Filesystem string            // set "." for no offset; "" is file load disabled.
	Psuedofs   map[string]string // map module name to script.  remember keys should end in ".sk" extension.

	// ModulePredeclared is a set of builtins, functions, and values that will be
	// pre-declared universally: in the main script and any modules loaded.
	//
	// You can also add additional pre-declared entities for the main script
	// when you evaluate it; these will be merged into this set, overriding it,
	// and will not be provided to any further modules loaded.
	// (Therefore: if you're going to provide side-effecting or
	// referentially-opaque functions, consider doing it as a parameter of
	// your main script eval call; it'll keep them from spreading.)
	ModulePredeclared sk.StringDict

	// map of module name to results, built as we evaluate.
	// nil entries are a sentinel value meaning "wip" (and if found, mean "cycle").
	evaluations map[string]*evaluation
}

type evaluation struct {
	globals sk.StringDict
	err     error
}

// Eval evaluates a "main" script.
//
// The filename parameter is optional and used only for logs and error messages.
//
// Additional predeclared values can be provided; these will be merged with the
// set of ModulePredeclared configured in the interpreter as a whole and made
// available to this "main" script (but not available to any subsequently
// loaded modules).
//
// Any modules loaded by `load` in the "main" script (or subsequently loaded
// recursively by other modules) will be memoized by the interpreter, but
// this "main" script itself is never memoized.
func (l *Interpreter) Eval(src string, filename string, additionalPredeclared sk.StringDict) (sk.StringDict, error) {
	if filename == "" {
		filename = "__main__"
	}
	if l.evaluations == nil {
		l.evaluations = make(map[string]*evaluation)
	}
	thread := &sk.Thread{Load: l.load}
	predeclared := mergeStringDict(l.ModulePredeclared, additionalPredeclared)
	globals, err := sk.Exec(sk.ExecOptions{
		Thread:   thread,
		Filename: filename, Source: src,
		Predeclared: predeclared,
	})
	return globals, err
}

// 'b' dominates
func mergeStringDict(a, b sk.StringDict) (c sk.StringDict) {
	c = make(sk.StringDict, len(a)+len(b))
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return
}

// per sk.Thread#Load signature.
func (l *Interpreter) load(parentThread *sk.Thread, module string) (sk.StringDict, error) {
	// Normalize the module path/name.
	//  Module names are roughly paths, but we don't allow them to up-dir,
	//  and they always end in ".sk" extension for consistency.
	module = filepath.Clean(module)
	if module[0] == '.' {
		return sk.StringDict{}, fmt.Errorf("module names may not start with '.' or '..'")
	}
	if !strings.HasSuffix(module, ".sk") {
		module += ".sk"
	}

	// Return memoized result if possible (or exit on cycle detect!).
	memo, ok := l.evaluations[module]
	if memo != nil {
		return memo.globals, memo.err
	}
	if ok && memo == nil {
		return nil, fmt.Errorf("cycle in load graph")
	}

	// Load the source (or a handle to it), or error.
	src, err := l.getSource(module)
	if err != nil {
		return sk.StringDict{}, err
	}
	if src == nil {
		return sk.StringDict{}, fmt.Errorf("no module %q found", module)
	}

	// Let's run!
	//  The `globals` var will be mutated by the exec.
	thread := &sk.Thread{Load: l.load}
	globals, err := sk.Exec(sk.ExecOptions{
		Thread:   thread,
		Filename: module, Source: src,
		Predeclared: l.ModulePredeclared,
	})

	// Remember remember the exec of module...
	l.evaluations[module] = &evaluation{globals, err}

	return globals, err
}

// Return something that skylark will be pleased to consider as a source.
// It's either a string or a reader.
// The return may also be (nil,nil), which means no source found.
//
// If a file reader is returned, it's not particularly easy to close it;
// we consider it fine to disregard this leakage, as no interpreter should be
// living long enough individually for this to become problematic.
func (l *Interpreter) getSource(path string) (interface{}, error) {
	// Check in-memory psuedo-fs first.
	if src, ok := l.Psuedofs[path]; ok {
		return src, nil
	}

	// If filesystem not enabled, ded.
	if l.Filesystem == "" {
		return nil, nil
	}

	// If enabled, try the filesystem.
	filename := filepath.Join(l.Filesystem, path)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err == nil {
		return f, nil
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}
