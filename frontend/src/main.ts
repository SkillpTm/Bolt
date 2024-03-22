import "./assets/styles/main.css";
import "./assets/styles/searchicons.css";
import "./assets/styles/component.css";

import { LaunchSearch } from "../wailsjs/go/app/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

import { StateHandler } from "./app/statehandler";
import { UIHandler } from "./app/uihandler";

/* <----------------------------------------------------------------------------------------------------> */

const stateHandler = new StateHandler;
const uiHandler = new UIHandler(7);

/* <----------------------------------------------------------------------------------------------------> */

// disable right click
document.oncontextmenu = () => {
    return false;
}

// focus the searchBar on load
window.onload = () => {
    uiHandler.searchBar.focus();
    uiHandler.reset();
}

// makes sure the searchBar is always clicked
document.addEventListener("click", () => {
    uiHandler.searchBar.focus();

});

// send the current input to Go to search the file system
uiHandler.searchBar.addEventListener("input", () => {
    stateHandler.handleSearch(uiHandler);
    LaunchSearch(uiHandler.searchBar.value);
});

// when Go found results receive, handle and display them
EventsOn("searchResult", (results: string[]) => {
    stateHandler.handleResult(results, uiHandler);
});

declare global {
    interface Window {

    }
}
