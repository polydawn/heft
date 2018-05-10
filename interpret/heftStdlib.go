package interpret

var heftLibsPsuedofs = map[string]string{
	"std.sk": `
def step(*components):
	result = components[0]
	for comp in components[1:]:
		result += comp
	return result

# Convert an array of strings into a "bash -c ${cmds}"-style array
# suitable for handing to "action(...)".
def concatBash(*cmds):
	return ["/bin/bash", "-c", "\n".join(cmds)]

# Yield a FormulaUnion fragment containing imports.
def reference(path, importableID):
	if type(importableID) != "ReleaseItemID":
		importableID = releaseItemID(importableID)
	if importableID.version == "":
		pass # help how do i error in skylark
	return formula({
		"imports":{
			path: importableID,
		},
	})

# Yield a FormulaUnion fragment containing an exec action.
def action(fragment):
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

# Yield a FormulaUnion fragment containing outputs.
def output(path, **kwargs):
	return formula({
		"formula":{"outputs":{
			path: {"packtype": kwargs['packtype']}
		}},
	})
`,
}
