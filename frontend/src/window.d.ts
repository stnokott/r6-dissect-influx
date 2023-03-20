import * as runtime from "../wailsjs/runtime/runtime";

export { };

declare global {
	interface Window {
		runtime: runtime;
	}
}

window.runtime = window.runtime || {};
