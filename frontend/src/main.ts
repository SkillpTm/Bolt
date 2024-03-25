import "./assets/styles/main.css";
import "./assets/styles/searchicons.css";
import "./assets/styles/component.css";

import { LaunchSearch } from "../wailsjs/go/app/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

import { StateHandler } from "./app/statehandler";
import { UIHandler } from "./ui/uihandler";

/* <----------------------------------------------------------------------------------------------------> */

const stateHandler = new StateHandler(new UIHandler(8));

/* <----------------------------------------------------------------------------------------------------> */

// disable right click
document.oncontextmenu = () => {
    return false;
}

// focus the searchBar on load
window.onload = async () => {
    stateHandler.uiHandler.searchBar.focus();
    stateHandler.uiHandler.searchBar.select();
    await stateHandler.reset();
}

// makes sure the searchBar is always clicked
document.addEventListener("click", () => {
    stateHandler.uiHandler.searchBar.focus();
});

// reset the app once it has been hidden
EventsOn("hidApp", async () => {
    await stateHandler.reset();
});

// move the highlighted section with arrow keys and open afile with enter
document.addEventListener("keydown", async (event) => {
    if (event.key === "ArrowDown") {
        event.preventDefault()
        stateHandler.search.updateHighlightedComp(1);
    } else if (event.key === "ArrowUp") {
        event.preventDefault()
        stateHandler.search.updateHighlightedComp(-1);
    } else if (event.key === "Enter") {
        await stateHandler.openFile(stateHandler.search.getHighlightedFile());
    }
});

// send the current input to Go to search the file system
stateHandler.uiHandler.searchBar.addEventListener("input", async () => {
    stateHandler.search.newSearch();
    await LaunchSearch(stateHandler.uiHandler.searchBar.value);
});

// when Go found results receive, handle and display them
EventsOn("searchResult", async (results: string[]) => {
    await stateHandler.search.newResults(results);
});

// change page forwards with shortcut
EventsOn("pageForward", async () => {
    await stateHandler.search.updatePage(1);
});

// change page backwards with shortcut
EventsOn("pageBackward", async () => {
    await stateHandler.search.updatePage(-1);
});