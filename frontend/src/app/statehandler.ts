export {StateHandler}

/* <----------------------------------------------------------------------------------------------------> */

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
}

