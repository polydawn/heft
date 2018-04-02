#!/usr/bin/env heft

load("std", "step", "reference", "action", "concatBash", "output")

## TODO start using hitch view funcs to pick versions to make this interesting!!

pipeline = basting(
	step1=step(
		reference("/",   "ports.timeless.io/busybash:201801:linux-amd64"),
		reference("/bar","foobar.org/bar:xyz:???"), ## FIXME exports in that module...
		reference("/baz","foobar.org/baz:xyz:???"), ## FIXME exports in that module...
		action(concatBash(
			"echo synthesizing something would be fun",
		)),
		output("/task/out", packtype='tar'),
	),
)
