import "./assets/styles/main.css";
import "./assets/styles/searchicons.css";
import "./assets/styles/component.css";

import { LaunchSearch } from "../wailsjs/go/app/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

import { StateHandler } from "./app/statehandler";
import { UIHandler } from "./app/uihandler";

/* <----------------------------------------------------------------------------------------------------> */

const maxComponents = 7;

const stateHandler = new StateHandler;
const uiHandler = new UIHandler(maxComponents);

uiHandler.components.forEach((comp) => {
    comp.self.addEventListener("click", () => {
        stateHandler.openFile(uiHandler, uiHandler.getHoverComp(comp));
    });
});

/* <----------------------------------------------------------------------------------------------------> */

// disable right click
document.oncontextmenu = () => {
    return false;
}

// focus the searchBar on load
window.onload = () => {
    uiHandler.searchBar.focus();
    uiHandler.searchBar.select();
    uiHandler.reset();
}

// makes sure the searchBar is always clicked
document.addEventListener("click", () => {
    uiHandler.searchBar.focus();
});

EventsOn("hidApp", () => {
    uiHandler.reset();
});

document.addEventListener('keydown', (event) => {
    if (event.key === 'ArrowDown') {
        event.preventDefault()
        uiHandler.updateHighlightedComp(1);
    } else if (event.key === 'ArrowUp') {
        event.preventDefault()
        uiHandler.updateHighlightedComp(-1);
    } else if (event.key === 'Enter') {
        stateHandler.openFile(uiHandler, uiHandler.getCurrentComp());
    }
});

// send the current input to Go to search the file system
uiHandler.searchBar.addEventListener("input", async () => {
    await stateHandler.updatePage(0, uiHandler);
    stateHandler.handleSearch(uiHandler);
    LaunchSearch(uiHandler.searchBar.value);
});

// when Go found results receive, handle and display them
EventsOn("searchResult", async (results: string[]) => {
    await stateHandler.handleResult(results, uiHandler);
});

EventsOn("pageForward", async () => {
    await stateHandler.updatePage(1, uiHandler);
});

EventsOn("pageBackward", async () => {
    await stateHandler.updatePage(-1, uiHandler);
});

declare global {
    interface Window {

    }
}
