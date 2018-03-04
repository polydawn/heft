package interpret

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sk "github.com/google/skylark"

	"go.polydawn.net/heft/skyform"
)

// Makes a new globals dict with our favorite custom bits in it.
func newGlobals() sk.StringDict {
	return sk.StringDict{
		"iamheft": sk.NewBuiltin("iamheft", func(thread *sk.Thread, fn *sk.Builtin, args sk.Tuple, kwargs []sk.Tuple) (sk.Value, error) {
			if thread.Print != nil {
				thread.Print(thread, "yes")
			} else {
				fmt.Fprintln(os.Stderr, "yes")
			}
			return sk.None, nil
		}),
		"formula":       sk.NewBuiltin("formula", skyform.MakeFormulaUnion),
		"basting":       sk.NewBuiltin("basting", skyform.NewBasting),
		"releaseItemID": sk.NewBuiltin("releaseItemID", skyform.NewReleaseItemID),
	}
}

func ExecFile(filename string) {
	thread := &sk.Thread{}
	globals := newGlobals()
	if err := sk.ExecFile(thread, filename, nil, globals); err != nil {
		fmt.Fprintf(os.Stderr, "larking: %s\n", err)
		os.Exit(4)
	}
}

type Loader struct {
	Filesystem string            // set "." for no offset; "" is file load disabled.
	Psuedofs   map[string]string // map module name to script.  remember keys should end in ".sk" extension.

	// map of module name to results, built as we evaluate.
	// nil entries are a sentinel value meaning "wip" (and if found, mean "cycle").
	evaluations map[string]*evaluation
}

type evaluation struct {
	globals sk.StringDict
	err     error
}

func (l *Loader) EvalScript(src string) (sk.StringDict, error) {
	thread := &sk.Thread{Load: l.load}
	globals := newGlobals()
	l.evaluations = make(map[string]*evaluation)
	err := sk.Exec(sk.ExecOptions{
		Thread:   thread,
		Filename: "__main__", Source: src,
		Globals: globals,
	})
	for name := range newGlobals() {
		delete(globals, name)
	}
	return globals, err
}

// per sk.Thread#Load signature.
func (l *Loader) load(parentThread *sk.Thread, module string) (sk.StringDict, error) {
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
	globals := make(sk.StringDict)
	err = sk.Exec(sk.ExecOptions{
		Thread:   thread,
		Filename: module, Source: src,
		Globals: globals,
	})

	// Censor our builtin funcs back out of the results.
	for name := range newGlobals() {
		delete(globals, name)
	}

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
func (l *Loader) getSource(path string) (interface{}, error) {
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
