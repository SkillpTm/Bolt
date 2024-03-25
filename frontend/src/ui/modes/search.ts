export { Search }

/* <----------------------------------------------------------------------------------------------------> */

import { GetImageData } from "../../../wailsjs/go/app/App";

import { Component, UIHandler } from "../uihandler";

/* <----------------------------------------------------------------------------------------------------> */

/**
 * Is in task of anything related to the search's UI.
 */
class Search {
    uiHandler!: UIHandler;

    #highlightedComp = 0;
    #results: Array<string> = [];
    #resultPage = 0;
    #searching = false;

    /**
     * Sets the uiHandler as a property. Also adds event listeners for the page buttons to the left and right icons.
     *
     * @param uiHandler the main uiHandler used throught the app
     */
    constructor(uiHandler: UIHandler) {
        this.uiHandler = uiHandler;

        this.uiHandler.leftIcon.addEventListener("mouseenter", async () => {await this.updateLeftRightIcons(true, true);});
        this.uiHandler.leftIcon.addEventListener("mouseleave", async () => {await this.updateLeftRightIcons(false, true);});
        this.uiHandler.leftIcon.addEventListener("click", async () => {
            await this.updatePage(-1);
            await this.updateLeftRightIcons(true, true);
        });

        this.uiHandler.rightIcon.addEventListener("mouseenter", async () => {await this.updateLeftRightIcons(true, false);});
        this.uiHandler.rightIcon.addEventListener("mouseleave", async () => {await this.updateLeftRightIcons(false, false);});
        this.uiHandler.rightIcon.addEventListener("click", async () => {
            await this.updatePage(1);
            await this.updateLeftRightIcons(true, false);
        });
    }

    /**
     * Removes any image from the right icon and starts the loading animation.
     */
    newSearch(): void {
        this.#searching = true;
        this.uiHandler.rightIcon.src = "";

        this.uiHandler.rightSection.classList.add("loading-grid");
    }

    /**
     * Stores the new results and displays the results.
     *
     * @param results the new results we received
     */
    async newResults(results: Array<string>): Promise<void> {
        this.#searching = false;
        this.#results = results;

        await this.updatePage(0)
    }

    /**
     * Updates the page number. If the change is in bounds of the results length. If the resultPage changes we also re-display the results.
     *
     * @param change increase/decrease to page count, 0 resets it back to 0
     */
    async updatePage(change: number): Promise<void> {
        if (change === 0) {
            this.#resultPage = 0;
            await this.displayResults();
            return;
        }

        // if the new page would be out of bounds we simply don't change the resultPage value
        if ((this.#resultPage + change) * this.uiHandler.components.length > this.#results.length - 1) {
            return;
        }

        if ((this.#resultPage + change) < 0) {
            return;
        }

        this.#resultPage += change;

        await this.displayResults();
    }

    /**
     * Adds the entry names, paths and icons to the components and displays them.
     */
    async displayResults(): Promise<void> {
        this.uiHandler.rightSection.classList.remove("loading-grid");

        if (this.uiHandler.searchBar.value.length === 0) {
            await this.uiHandler.resetUI();
            return;
        }

        if (this.#results.length > 0) {
            this.uiHandler.rightIcon.src = await GetImageData("tick");
        } else {
            this.uiHandler.rightIcon.src = await GetImageData("cross");
        }

        let displayAmount = 0;

        for (let index = 0; index < this.#results.length && index < this.uiHandler.components.length; index++) {
            const currentFile = index + (this.uiHandler.components.length * this.#resultPage);

            if (currentFile > this.#results.length-1) {
                break;
            }

            let filePath = this.#results[currentFile].split("/");

            // if the last element is empty, it means our string ended in a slash, indicating it was a folder.
            if (filePath[filePath.length-1] === "") {
                filePath.pop();

                this.uiHandler.components[index].image.src = await GetImageData("folder");
            } else {
                this.uiHandler.components[index].image.src = await GetImageData("file");
            }

            this.uiHandler.components[index].name.textContent = filePath.pop() as string;
            this.uiHandler.components[index].value.textContent = filePath.join("/") + "/";

            displayAmount++;
        }

        this.updateHighlightedComp(0);
        this.uiHandler.displayComponents(displayAmount);
    }

    /**
     * Updates the higlighted component
     *
     * @param change increase/decrease to component position, 0 resets it back to 0. Incase of an overflow to the max components it rolls back to 0 and the other way around.
     */
    updateHighlightedComp(change: number): void {
        this.uiHandler.components[this.#highlightedComp].self.classList.remove("highligthed");

        if (change === 0) {
            this.#highlightedComp = 0;
        } else if (change < 0) {
            this.#highlightedComp = (this.#highlightedComp + change + this.uiHandler.displayedComps) % this.uiHandler.displayedComps;
        } else {
            this.#highlightedComp = (this.#highlightedComp + change) % this.uiHandler.displayedComps;
        }

        this.uiHandler.components[this.#highlightedComp].self.classList.add("highligthed");
    }

    /**
     * Gets the path to the highlighted component's file.
     * 
     * @returns the full path of highlighted component.
     */
    getHighlightedFile(): string {
        return this.#results[this.#resultPage * this.uiHandler.components.length + this.#highlightedComp];
    }

    /**
     * Gets the path to the hovered over component's file.
     *
     * @param comp the component over which is being hovered
     * 
     * @returns the full path of hovered over component.
     */
    getHoveredFile(comp: Component): string {
        return this.#results[this.#resultPage * this.uiHandler.components.length + (parseInt(comp.self.id[9]) -1)];
    }

    /**
     * Changes the left or right icon to an arrow for pag changes when approriate.
     *
     * @param mouseEnter if we just entered or left the icon
     * 
     * @param left if we need to change the left or right icon
     */
    async updateLeftRightIcons(mouseEnter: boolean, left: boolean): Promise<void> {
        const resetIcons = async () => {
            if (left) {
                this.uiHandler.leftIcon.classList.remove("icon-clickable");

                this.uiHandler.leftIcon.src = await GetImageData("magnifying_glass");
            } else {
                this.uiHandler.rightIcon.classList.remove("icon-clickable");

                if (this.#results.length > 0) {
                    this.uiHandler.rightIcon.src = await GetImageData("tick");
                } else {
                    this.uiHandler.rightIcon.src = await GetImageData("cross");
                }
            }
        }

        if (this.uiHandler.searchBar.value.length === 0 || this.#results.length === 0 || this.#searching || !mouseEnter) {
            await resetIcons();
            return;
        }

        if (left) {
            this.uiHandler.leftIcon.classList.add("icon-clickable");

            if (this.#resultPage > 0) {
                this.uiHandler.leftIcon.src = await GetImageData("left");
            } else {
                this.uiHandler.leftIcon.src = await GetImageData("not-left");
            }
        } else {
            this.uiHandler.rightIcon.classList.add("icon-clickable");

            if (this.#resultPage * this.uiHandler.components.length < this.#results.length - 1) {
                this.uiHandler.rightIcon.src = await GetImageData("right");
            } else {
                this.uiHandler.rightIcon.src = await GetImageData("not-right");
            }
        }
    }
}