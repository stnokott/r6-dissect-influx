//go:build pprof

package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CPU Profiling, currently only used for PGO (https://go.dev/doc/pgo)

const ENV_FILENAME string = "PPROF_FILENAME"

func init() {
	filename, ok := os.LookupEnv(ENV_FILENAME)
	if !ok {
		log.Fatalf("environment variable %s not found!", ENV_FILENAME)
	}

	onStartupFuncs = append(onStartupFuncs, func(_ context.Context) {
		cpuProfile, _ := os.Create("build/cpu_profiles/hosts/" + filename + ".pprof")
		pprof.StartCPUProfile(cpuProfile)
	})

	onDomReadyFuncs = append(onDomReadyFuncs, func(ctx context.Context) {
		runtime.WindowExecJS(ctx, `
		  const Sleep = (ms) => {
				return new Promise(r => setTimeout(r, ms));
			}

		  const f = async () => {
				// show window
				window.runtime.Show();
				await Sleep(5000);
				// scrolling
				let matchContainer = document.querySelector("#match-container");
				// scroll to last item
				matchContainer.children[matchContainer.children.length - 1].scrollIntoView({ behavior: "smooth" });
				await Sleep(10 * 1000);
				// scroll to first item
				matchContainer.children[0].scrollIntoView({ behavior: "smooth" });
				// wait 30s to gather some CPU profile data
				await Sleep(30 * 1000);
				// close application
				let btnCloseApp = document.querySelector('#root > #btn-close-application:last-child');
				btnCloseApp.click();
			}

			f();
		`)
	})

	onShutdownFuncs = append(onShutdownFuncs, func(_ context.Context) {
		pprof.StopCPUProfile()
	})
}
