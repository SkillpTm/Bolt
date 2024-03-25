export {StateHandler}

/* <----------------------------------------------------------------------------------------------------> */

import { OpenFileExplorer } from "../../wailsjs/go/app/App";
import { WindowHide } from "../../wailsjs/runtime/runtime";

import { UIHandler } from "../ui/uihandler";
import { Search } from "../ui/modes/search";

/* <----------------------------------------------------------------------------------------------------> */

/**
 * Holds the uiHandler and all the UI modes.
 */
class StateHandler {
    uiHandler!: UIHandler;
    search!: Search;

    /**
     * Sets the uiHandler as a property and adds the properties: Search.
     *
     * @param uiHandler the main uiHandler used throught the app
     */
    constructor(uiHandler: UIHandler) {
        this.uiHandler = uiHandler;
        this.search = new Search(uiHandler);

        this.uiHandler.components.forEach((comp) => {
            comp.self.addEventListener("click", async () => {
                await this.openFile(this.search.getHoveredFile(comp));
            });
        });
    }

    async reset(): Promise<void> {
        await this.uiHandler.resetUI();
        await this.search.newResults([] as Array<string>);
    }

    // openFile opens the given file and hides the app
    /**
     * Open the given file with a windows file explorer.
     *
     * @param filePath the path for the file to be opened
     */
    async openFile(filePath: string): Promise<void> {
        WindowHide();
        await this.reset();
        OpenFileExplorer(filePath);
    }
}