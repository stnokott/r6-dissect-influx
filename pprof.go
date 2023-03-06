//go:build pprof

package main

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime/pprof"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CPU Profiling, currently only used for PGO (https://go.dev/doc/pgo)

const ENV_FILENAME string = "PPROF_FILENAME"
const ROOT_PATH string = "build/cpu_profiles/hosts/"

func init() {
	filename, ok := os.LookupEnv(ENV_FILENAME)
	if !ok {
		log.Fatalf("environment variable %s not found!", ENV_FILENAME)
	}

	onStartupFuncs = append(onStartupFuncs, func(_ context.Context) {
		if err := os.MkdirAll(ROOT_PATH, 0644); err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatalf("could not create directory structure %s: %v", ROOT_PATH, err)
		}
		cpuProfile, _ := os.Create(ROOT_PATH + filename + ".pprof")
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
		log.Println()
		log.Println("############################# IMPORTANT ###########################")
		log.Println("# Remember to merge profiles:                                     #")
		log.Println("# > go install github.com/google/pprof@latest                     #")
		log.Println("# > cd build/cpu_profiles                                         #")
		log.Println("# > pprof -proto -output merged.pprof hosts/a.pprof hosts/b.pprof #")
		log.Println("###################################################################")
	})
}
