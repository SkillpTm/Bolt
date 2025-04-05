export { SearchModule }

import { LaunchSearch } from "../../../wailsjs/go/app/App";

import { Component, UIHandler } from "../uihandler";

/**
 * Is in task of anything related to the search's UI.
 * 
 * @param maxResultsDispalyed private property, how many results can be displayed at once
 * 
 * @param resultPage private property, current page of the results we're on
 * 
 * @param results private property, holds all current results
 * 
 * @param searching private property, state of if we're searching right now
 * 
 * @param uiHandler main uiHandler user to access its base functions and properties
 */
class SearchModule {
    #maxDisplayedResults = 0;
    #resultPage = 0;
    results: Array<string> = [];
    #searching = false;
    uiHandler!: UIHandler;

    constructor(uiHandler: UIHandler, maxDisplayedResults: number) {
        this.uiHandler = uiHandler;
        this.#maxDisplayedResults = maxDisplayedResults;

        this.uiHandler.leftIcon.addEventListener("mouseenter", () => {this.updateLeftRightIcons(true, true)});
        this.uiHandler.leftIcon.addEventListener("mouseleave", () => {this.updateLeftRightIcons(false, true)});
        this.uiHandler.leftIcon.addEventListener("click", () => {
            this.updatePage(-1);
            this.updateLeftRightIcons(true, true);
        });

        this.uiHandler.rightIcon.addEventListener("mouseenter", () => {this.updateLeftRightIcons(true, false)});
        this.uiHandler.rightIcon.addEventListener("mouseleave", () => {this.updateLeftRightIcons(false, false)});
        this.uiHandler.rightIcon.addEventListener("click", () => {
            this.updatePage(1);
            this.updateLeftRightIcons(true, false);
        });
    }

    /**
     * Removes any image from the right icon and starts the loading animation.
     */
    async newInput(): Promise<void> {
        this.#searching = true;
        this.uiHandler.rightIcon.src = "";

        this.uiHandler.rightSection.classList.add("loading-grid");
        this.uiHandler.rightSection.classList.remove("hide");
        this.uiHandler.rightIcon.classList.add("hide");

        await LaunchSearch(this.uiHandler.searchBar.value);
    }

    /**
     * Stores the new results and displays them.
     *
     * @param results the new results we received
     */
    newResults(results: Array<string>): void {
        this.#searching = false;
        this.results = results;

        this.updatePage(0);
    }

    /**
     * Updates the page number. If the change is in bounds of the results length. If the resultPage changes we also re-display the results.
     *
     * @param change increase/decrease to page count, 0 resets it back to 0
     */
    updatePage(change: number): void {
        if (change === 0) {
            this.#resultPage = 0;
            this.displayResults();
            return;
        }

        // if the new page would be out of bounds we simply don't change the resultPage value
        if ((this.#resultPage + change) * (this.#maxDisplayedResults) > (this.results.length - 1)) {
            return;
        }

        if ((this.#resultPage + change) < 0) {
            return;
        }

        this.#resultPage += change;

        this.displayResults();
    }

    /**
     * Adds the entry names, paths and icons to the components and displays them.
     */
    displayResults(): void {
        this.uiHandler.rightSection.classList.remove("loading-grid");
        this.uiHandler.rightSection.classList.add("hide");
        this.uiHandler.rightIcon.classList.remove("hide");

        if (this.uiHandler.searchBar.value.length === 0) {
            this.uiHandler.resetUI();
            return;
        }

        if (this.results.length > 0) {
            this.uiHandler.rightIcon.src = this.uiHandler.images.get("tick") as string;
        } else {
            this.uiHandler.rightIcon.src = this.uiHandler.images.get("cross") as string;
        }

        const displayComps: Array<number> = [];

        for (let index = 0; index < this.results.length - (this.#resultPage * this.#maxDisplayedResults) && index < (this.#maxDisplayedResults); index++) {
            const currentFile = index + (this.#resultPage * (this.#maxDisplayedResults));

            const filePath = this.results[currentFile].split("/");

            if (this.results[currentFile].endsWith("/")) {
                // pop empty element
                filePath.pop();

                this.uiHandler.components[index+1].image.src = this.uiHandler.images.get("folder") as string;
            } else {
                this.uiHandler.components[index+1].image.src = this.uiHandler.images.get("file") as string;
            }

            this.uiHandler.components[index+1].tooltip.textContent = filePath.join("/") as string;
            this.uiHandler.components[index+1].name.textContent = filePath.pop() as string;
            this.uiHandler.components[index+1].value.textContent = filePath.join("/") + "/";

            displayComps.push(index+1);
        }

        // the 2nd input produces an array with all values between 1-7 that aren't in displayComps
        this.uiHandler.displayComponents(displayComps, Array.from({length: 6}, (_, i) => i+1).filter(item => !displayComps.includes(item)));
        this.uiHandler.updateHighlightedComp(undefined, true);
    }

    /**
     * Gets the path to the highlighted component's file.
     * 
     * @returns the full path of highlighted component.
     */
    getHighlightedFile(): string {
        return this.uiHandler.components[this.uiHandler.getHighlightedComp()].tooltip.textContent as string;
    }

    /**
     * Gets the path to the hovered over component's file.
     *
     * @param comp the component over which is being hovered
     * 
     * @returns the full path of hovered over component.
     */
    getHoveredFile(comp: Component): string {
        return comp.tooltip.textContent as string;
    }

    /**
     * Changes the left or right icon to an arrow for page changes when approriate.
     *
     * @param mouseEnter if we just entered or left the icon
     * 
     * @param left if we need to change the left or right icon
     */
    updateLeftRightIcons(mouseEnter: boolean, left: boolean): void {
        // check, if we're in a state, in which the arrows shouldn't appear
        if (this.uiHandler.searchBar.value.length === 0 || this.results.length === 0 || this.#searching || !mouseEnter) {
            if (left) {
                this.uiHandler.leftIcon.classList.remove("icon-clickable");

                this.uiHandler.leftIcon.src = this.uiHandler.images.get("magnifying_glass") as string;
            } else {
                this.uiHandler.rightIcon.classList.remove("icon-clickable");

                if (this.results.length > 0) {
                    this.uiHandler.rightIcon.src = this.uiHandler.images.get("tick") as string;
                } else {
                    this.uiHandler.rightIcon.src = this.uiHandler.images.get("cross") as string;
                }
            }
        } else {
            if (left) {
                this.uiHandler.leftIcon.classList.add("icon-clickable");
    
                if (this.#resultPage > 0) {
                    this.uiHandler.leftIcon.src = this.uiHandler.images.get("left") as string;
                } else {
                    this.uiHandler.leftIcon.src = this.uiHandler.images.get("not-left") as string;
                }
            } else {
                this.uiHandler.rightIcon.classList.add("icon-clickable");
    
                if (this.#resultPage * this.#maxDisplayedResults + this.uiHandler.getDisplayedComps().length < (this.results.length-1)) {
                    this.uiHandler.rightIcon.src = this.uiHandler.images.get("right") as string;
                } else {
                    this.uiHandler.rightIcon.src = this.uiHandler.images.get("not-right") as string;
                }
            }
        }
    }
}