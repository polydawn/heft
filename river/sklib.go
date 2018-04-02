package river

var LibPsuedofs = map[string]string{
	"std.sk": `
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
`,
}
