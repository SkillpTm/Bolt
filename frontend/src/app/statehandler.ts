export {StateHandler}

import { OpenFileExplorer } from "../../wailsjs/go/app/App";
import { BrowserOpenURL, WindowHide } from "../../wailsjs/runtime/runtime";

import { UIHandler } from "../ui/uihandler";
import { SearchMode } from "../ui/modes/search";

/**
 * Holds the uiHandler and all the UI modes.
 * 
 * @param uiHandler the main uiHandler used to manipulate the UI
 * 
 * @param searchMode mode used to change the UI depending on the search state
 */
class StateHandler {
    uiHandler!: UIHandler;
    searchMode!: SearchMode;

    /**
     * Sets the uiHandler as a property and adds the property: Search.
     */
    constructor() {
        this.uiHandler = new UIHandler(8);
        this.searchMode = new SearchMode(this.uiHandler);

        this.uiHandler.components.forEach((comp) => {
            comp.self.addEventListener("click", async () => {
                await this.openFile(this.searchMode.getHoveredFile(comp));
            });
        });

        this.reset()
    }

    /**
     * resets the ui and state of the frontend
     */
    reset(): void {
        this.uiHandler.resetUI();
        this.searchMode.newResults([] as Array<string>);
    }

    /**
     * Open the given file with a file manager/the search result in the browser and hide the search window.
     *
     * @param filePath the path for the file to be opened. If "<web-search>" opens the search result in the browser
     */
    async openFile(filePath: string): Promise<void> {
        WindowHide();

        if (filePath === "<web-search>") {
            BrowserOpenURL(`https://www.google.com/search?q=${this.uiHandler.searchBar.value}`);
        } else {
            await OpenFileExplorer(filePath);
        }

        this.reset();
    }
}