import "../styles/main.css";
import "../styles/searchicons.css";
import "../styles/component.css";

import { LaunchSearch } from "../wailsjs/go/app/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

import { StateHandler } from "./app/statehandler";

const stateHandler = new StateHandler();

// disable right click
document.oncontextmenu = () => {
    return false;
}

// focus the searchBar on load
window.onload = () => {
    stateHandler.uiHandler.searchBar.focus();
    stateHandler.reset();
}

// makes sure the searchBar is always focused
document.addEventListener("click", () => {
    stateHandler.uiHandler.searchBar.focus();
});

// move the highlighted section with arrow keys and open a file with enter
document.addEventListener("keydown", async (event) => {
    if (event.key === "ArrowDown") {
        event.preventDefault()
        stateHandler.searchMode.updateHighlightedComp(1);
    } else if (event.key === "ArrowUp") {
        event.preventDefault()
        stateHandler.searchMode.updateHighlightedComp(-1);
    } else if (event.key === "Enter") {
        await stateHandler.openFile(stateHandler.searchMode.getHighlightedFile());
    }
});

// send the current input to Go to search the file system
stateHandler.uiHandler.searchBar.addEventListener("input", async () => {
    stateHandler.searchMode.newSearch();
    await LaunchSearch(stateHandler.uiHandler.searchBar.value);
});

// store the base64 imageData on the uiHandler
EventsOn("imageData", (imageData: Map<string, string>) => {
    stateHandler.uiHandler.images = imageData;
});

// when Go found results receive, handle and display them
EventsOn("searchResult", (results: string[]) => {
    stateHandler.searchMode.newResults(results);
});