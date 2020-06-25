/*
Package proftimer provides a quick and dirty, process-global timers for ad-hoc
profiling of sections of a codebase.

Using

	proftimer.Resume("mytimer")
	// Code that takes some time...
	proftimer.Accum("mytimer")
	proftimer.Report(os.Stdout, "mytimer")

You can resume and accumulate multiple timers at the same time.

*/
package proftimer
