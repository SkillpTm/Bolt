import "../styles/main.css";
import "../styles/searchicons.css";
import "../styles/component.css";

import { LogErrorTS } from "../wailsjs/go/app/App";
import { EventsOn, WindowSetSize } from "../wailsjs/runtime/runtime";

import { StateHandler } from "./app/statehandler";

const stateHandler = new StateHandler();

// disable right click
document.oncontextmenu = () => {
	return false;
};

// focus the searchBar on load
window.onload = () => {
	stateHandler.uiHandler.searchBar.focus();
	stateHandler.reset();
};

// makes sure the searchBar is always focused
document.addEventListener("click", () => {
	stateHandler.uiHandler.searchBar.focus();
});

document.addEventListener("keydown", async (event) => {
	if (event.ctrlKey && event.key === "ArrowUp") {
		event.preventDefault();
		stateHandler.searchMode.updatePage(-1);
	} else if (event.ctrlKey && event.key === "ArrowDown") {
		event.preventDefault();
		stateHandler.searchMode.updatePage(1);
	} else if (event.key === "ArrowUp" && stateHandler.uiHandler.getDisplayedComps().length > 0) {
		event.preventDefault();
		stateHandler.uiHandler.updateHighlightedComp(false);
	} else if (event.key === "ArrowDown" && stateHandler.uiHandler.getDisplayedComps().length > 0) {
		event.preventDefault();
		stateHandler.uiHandler.updateHighlightedComp(true);
	} else if (event.key === "Enter" && stateHandler.uiHandler.getDisplayedComps().length > 0) {
		await stateHandler.routeAction();
	}
});

// when Go found results receive, handle and display them
EventsOn("searchResult", (results: string[]) => {
	stateHandler.searchMode.newResults(results);
	WindowSetSize(570, stateHandler.uiHandler.topBarHeight + stateHandler.uiHandler.getDisplayedComps().length * stateHandler.uiHandler.componentHeight);
});

// catches all synchronous errors and passes them for error logging to Go
window.onerror = function (_message, source, lineno, colno, error) {
	LogErrorTS(
		`${source ?? "Unknown Source"}: ${lineno ?? "?"}, ${colno ?? "?"}:`.replaceAll("\n", " ") + "\n" +
		`--> ${error?.message ?? "unknown error"}:`.replaceAll("\n", " ") + "\n" +
		`[TS] ${error?.name ?? "Error"}`.replaceAll("\n", " "),
	);
	return true;
};

// catches all Promise errors and passes them for error logging to Go
window.addEventListener("unhandledrejection", (event) => {
	LogErrorTS(
		`${event.reason?.stack ?? "Unknown Source"}`.replaceAll("\n", " ") + "\n" +
		`--> ${event.reason?.message ?? event.reason ?? "unknown error"}`.replaceAll("\n", " ") + "\n" +
		`[TS] Promise Rejection: ${event.reason?.name ?? "Error"}`.replaceAll("\n", " "),
	);
});