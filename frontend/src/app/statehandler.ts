export {StateHandler}

/* <----------------------------------------------------------------------------------------------------> */

import { OpenFileExplorer } from "../../wailsjs/go/app/App";
import { WindowHide } from "../../wailsjs/runtime/runtime";

import { UIHandler } from "./uihandler";

/* <----------------------------------------------------------------------------------------------------> */

// The StateHandler directs the UIHandler to display tur correct state of teh application
class StateHandler {

    #query: String;
    #results: Array<String>;

    // constructor sets zero values for the state parameters
    constructor() {
        this.#query = "";
        this.#results = [];
    }

    // handleSearch updates the nav-bar for the newly started search
    handleSearch(uiHandler: UIHandler): void {
        this.#query = uiHandler.searchBar.value;

        uiHandler.startedSearch();
    }

    // handleResult updates the nav-bar and components for the finished search
    handleResult(newResults: Array<String>, uiHandler: UIHandler): void {
        this.#results = newResults;

        uiHandler.displayResults(this.#query.length, this.#results);
    }

    // openFile opens the given file and hides the app
    openFile(uiHandler: UIHandler, fileIndex: number) {
        WindowHide();
        uiHandler.reset();
        OpenFileExplorer(this.#results[fileIndex].replaceAll("/", "\\"));
    }
}