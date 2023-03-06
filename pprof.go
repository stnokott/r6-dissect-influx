//go:build pprof

package main

import (
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

func init() {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	beforeFuncs = append(beforeFuncs, func() {
		cpuProfile, _ := os.Create("build/cpu_profiles/" + host + ".pprof")
		pprof.StartCPUProfile(cpuProfile)
	})

	afterFuncs = append(afterFuncs, func() {
		pprof.StopCPUProfile()
	})
}
