export {StateHandler}

/* <----------------------------------------------------------------------------------------------------> */

import { OpenFileExplorer } from "../../wailsjs/go/app/App";
import { WindowHide } from "../../wailsjs/runtime/runtime";

import { UIHandler } from "./uihandler";

/* <----------------------------------------------------------------------------------------------------> */

// The StateHandler directs the UIHandler to display tur correct state of the application
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
    async handleResult(newResults: Array<String> | undefined, uiHandler: UIHandler): Promise<void> {
        if (newResults) {
            this.#results = newResults;
        }

        await uiHandler.displayResults(this.#query.length, this.#results);
    }

    /**
     * Updates the page number on the given UIHandler. If the change isn't 0 it also executes handleReslut.
     *
     * @param change increase/decrease to page count, 0 resets it back to 0
     * 
     * @param uiHandler the handler on which to update the page
     */
    async updatePage(change: number, uiHandler: UIHandler): Promise<void> {
        uiHandler.updatePage(change, this.#results.length);

        if (change !== 0) {
            await this.handleResult(undefined, uiHandler);
        }
    }

    // openFile opens the given file and hides the app
    openFile(uiHandler: UIHandler, fileIndex: number) {
        WindowHide();
        uiHandler.reset();
        OpenFileExplorer(this.#results[fileIndex].replaceAll("/", "\\"));
    }
}