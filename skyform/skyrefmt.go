package skyform

import (
	"fmt"

	sk "github.com/google/skylark"
	. "github.com/polydawn/refmt/tok"
)

type valueTokenizer struct {
	stack   []decoderStep // When empty, and step returns done, all done.
	current decoderStep   // Shortcut to end of stack.
}

func NewValueTokenizer() (d *valueTokenizer) {
	d = &valueTokenizer{
		stack: make([]decoderStep, 0, 10),
	}
	return
}

func (d *valueTokenizer) Bind(val sk.Value) {
	d.stack = d.stack[0:0]
	d.current = decoderStep{val, -1, false}
}

type decoderStep struct {
	val   sk.Value // in practice either Dict or List or Tuple, since leafs finish instantly
	index int      // For lists, duh; for dict, it's in order of AttrNames.
	value bool     // for dicts: true if it's time to emit the value
}

func (d *valueTokenizer) Step(tokenSlot *Token) (done bool, err error) {
	//fmt.Printf(":depths=%d ", len(d.stack))
	done, err = d.step(tokenSlot)
	// If the step errored: out, entirely.
	if err != nil {
		return true, err
	}
	// If the step wasn't done, return same status.
	if !done {
		return false, nil
	}
	// If it WAS done, and stack empty, we're entirely done.
	nSteps := len(d.stack) - 1
	if nSteps < 0 {
		return true, nil // that's all folks
	}
	// Pop the stack.
	d.current = d.stack[nSteps]
	d.stack = d.stack[0:nSteps]
	return false, nil
}

func (d *valueTokenizer) step(tokenSlot *Token) (done bool, err error) {
	switch v := d.current.val.(type) {
	case sk.NoneType:
		tokenSlot.Type = TNull
		return true, nil
	case sk.Bool:
		tokenSlot.Type = TBool
		tokenSlot.Bool = v == sk.True
		return true, nil
	case sk.Int:
		tokenSlot.Type = TInt
		tokenSlot.Int, _ = v.Int64()
		return true, nil
	case sk.Float:
		tokenSlot.Type = TFloat64
		tokenSlot.Float64 = float64(v)
		return true, nil
	case sk.String:
		tokenSlot.Type = TString
		tokenSlot.Str = string(v)
		return true, nil
	case sk.Indexable: // Tuple, List
		switch {
		case d.current.index < 0:
			//fmt.Printf(":: open arr (len=%d)\n", v.Len())
			tokenSlot.Type = TArrOpen
			tokenSlot.Length = sk.Len(v)
			d.current.index++
			return false, nil
		case d.current.index == sk.Len(v):
			//fmt.Printf(":: close arr\n")
			tokenSlot.Type = TArrClose
			return true, nil
		default:
			//fmt.Printf(":: arr step %d\n", d.current.index)
			d.current.index++
			d.pushPhase(decoderStep{v.Index(d.current.index - 1), -1, false})
			return d.step(tokenSlot)
		}
	case *sk.Dict:
		switch {
		case d.current.index < 0:
			//fmt.Printf(":: open map (len=%d)\n", v.Len())
			tokenSlot.Type = TMapOpen
			tokenSlot.Length = v.Len()
			d.current.index++
			return false, nil
		case d.current.index == v.Len():
			//fmt.Printf(":: close map\n")
			tokenSlot.Type = TMapClose
			return true, nil
		default:
			//fmt.Printf(":: map step %d+%v\n", d.current.index, d.current.value)
			if !d.current.value {
				tokenSlot.Type = TString
				tokenSlot.Str = string(v.Items()[d.current.index][0].(sk.String))
				d.current.value = true
				return false, nil
			} else {
				d.current.value = false
				d.current.index++
				d.pushPhase(decoderStep{v.Items()[d.current.index-1][1], -1, false})
				return d.step(tokenSlot)
			}
		}
	default:
		// function, builtin_function_or_method, set, and all user-defined types.
		return true, fmt.Errorf("cannot convert %s to JSON", v.Type())
	}
}

func (d *valueTokenizer) pushPhase(newPhase decoderStep) {
	d.stack = append(d.stack, d.current)
	d.current = newPhase
}
