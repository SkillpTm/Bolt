export { LinkModule }

import { UIHandler } from "../uihandler";

/**
 * LinkModule is in charge of detection links and showing it's module
 * 
 * @param isWebsite regex to detect websites
 * 
 * @param linkComp index of the component the link is supposed to be shown for
 * 
 * @param uiHandler the main uiHandler used throught the app
 */
class LinkModule {
    readonly isWebsite = /^(https?:\/\/)?((?![-])[a-z0-9-]+\.)+[a-z]+(\/?)/i;

    #linkComp: number;
    uiHandler!: UIHandler;

    constructor(uiHandler: UIHandler, compIndex: number) {
        this.uiHandler = uiHandler;
        this.#linkComp = compIndex;
    }

    /**
     * Checks if the new input is a link/domain and shows the component, if it is
     */
    newInput(): void {
        if (this.isWebsite.test(this.uiHandler.searchBar.value.trim())) {
            this.showComp();
        }
    }

    /**
     * modifies the #linkComp to show the suggestion to open the link/domain in the browser
     */
    showComp(): void {
        this.uiHandler.components[this.#linkComp].image.src = this.uiHandler.images.get("file") as string;
        this.uiHandler.components[this.#linkComp].tooltip.textContent = this.uiHandler.searchBar.value;
        this.uiHandler.components[this.#linkComp].name.textContent = this.uiHandler.searchBar.value;
        this.uiHandler.components[this.#linkComp].value.textContent = "Open link in browser";
        this.uiHandler.components[this.#linkComp].value.classList.add("browserInteract")

        this.uiHandler.displayComponents([this.#linkComp]);
    }
}