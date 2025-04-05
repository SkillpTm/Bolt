export { StateHandler }

import { HideWindow, OpenFileExplorer } from "../../wailsjs/go/app/App";
import { BrowserOpenURL, WindowSetSize } from "../../wailsjs/runtime/runtime";

import { Component, UIHandler } from "../ui/uihandler";
import { SearchModule } from "../ui/modes/search";
import { LinkModule } from "../ui/modes/link";

/**
 * Holds the uiHandler and all the UI modes.
 * 
 * @param linkModule the module in charge of detection and displaying link opens
 * 
 * @param searchMode module used to change the UI depending on the search state
 * 
 * @param uiHandler the main uiHandler used to manipulate the UI
 */
class StateHandler {
    linkModule!: LinkModule;
    searchMode!: SearchModule;
    uiHandler!: UIHandler;

    constructor() {
        this.uiHandler = new UIHandler(8);
        this.searchMode = new SearchModule(this.uiHandler, 6);
        this.linkModule = new LinkModule(this.uiHandler, 0);

        this.uiHandler.components.forEach((comp) => {
            comp.self.addEventListener("click", async () => {
                await this.routeAction(comp);
            });
        });

        // send the current input to Go to search the file system
        this.uiHandler.searchBar.addEventListener("input", async () => {
            this.handleInput();
        });

        // if the input bar is not selected anymore the user selected another window, so we hide
        this.uiHandler.searchBar.addEventListener("blur", () => {
            setTimeout(() => {
                if (document.activeElement === this.uiHandler.searchBar) {
                    HideWindow();
                    this.reset();
                }
            }, 50);
        });

        this.reset();
    }

    /**
     * Essentially acts as an event to act upon a new input
     */
    async handleInput(): Promise<void> {
        this.uiHandler.displayComponents(undefined, Array.from({length: 8}, (_, i) => i));
        await this.searchMode.newInput();
        this.linkModule.newInput();
    }

    /**
     * Resets the ui and state of the frontend
     */
    reset(): void {
        this.uiHandler.resetUI();
        this.searchMode.newResults(new Array<string>);

        WindowSetSize(570, this.uiHandler.topBarHeight + this.uiHandler.getDisplayedComps().length * this.uiHandler.componentHeight);
    }

    /**
     * Handles enter/left click to open the file manager/browser
     * 
     * @param clickComp if this was started by a left click, this is the clicked component
     */
    async routeAction(clickComp?: Component): Promise<void> {
        HideWindow();

        let currentComp: Component;
        if (clickComp) {
            currentComp = clickComp;
        } else {
            currentComp = this.uiHandler.components[this.uiHandler.getHighlightedComp()];
        }

        if (this.linkModule.isWebsite.test(this.uiHandler.searchBar.value.trim())) {
            BrowserOpenURL(this.uiHandler.searchBar.value.trim());
        } else if (this.searchMode.results.length > 0) {
            await OpenFileExplorer(currentComp.tooltip.textContent as string);
        }

        this.reset();
    }
}