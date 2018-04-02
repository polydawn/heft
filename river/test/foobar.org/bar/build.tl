#!/usr/bin/env heft

load("std", "step", "reference", "action", "concatBash", "output")

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
