#!/usr/bin/env heft

#
# look at all this preamble
#

def step(*components):
	result = components[0]
	for comp in components[1:]:
		result += comp
	return result

def concatBash(*cmds): # returns FormulaUnion fragment
	return ["bash", "-c", "\n".join(cmds)]

def reference(path, importableID): # returns FormulaUnion fragment
	if type(importableID) != "ReleaseItemID":
		importableID = releaseItemID(importableID)
	if importableID.version == "":
		pass # help how do i error in skylark
	return formula({
		"imports":{
			path: importableID,
		},
	})

def action(fragment): # returns FormulaUnion fragment
	t = type(fragment)
	if t == "string":
		return formula({
			"formula":{"action":{
				"exec": [fragment],
			}},
		})
	if t == "list":
		return formula({
			"formula":{"action":{
				"exec": fragment,
			}},
		})
	if t == "obj":
		return formula({
			"formula":{"action":fragment},
		})

def output(path, **kwargs):
	return formula({
		"formula":{"outputs":{
			path: {"packtype": kwargs['packtype']}
		}},
	})

#
# assign a pipeline!  this is what we "return".
#

pipeline = basting(
	step1=step(
		reference("/","ports.timeless.io/busybash:201801:linux-amd64"),
		action(concatBash(
			"mkdir out",
			"ls -la /app/go | tee > out/records",
		)),
		output("/task/out", packtype='tar'),
	),
)
