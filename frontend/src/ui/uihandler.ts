export { UIHandler, type Component };

/* <----------------------------------------------------------------------------------------------------> */

import { WindowSetSize } from "../../wailsjs/runtime/runtime";

import { GetImageData } from "../../wailsjs/go/app/App";

/* <----------------------------------------------------------------------------------------------------> */

/**
 * Is in task of doing general work around the components.
 *
 * A component looks like this in html:
 * <div id="component1" class="hide">
 *     <img id="component1-image" class="compImg">
 *     <div id="component1-text" class="compText">
 *         <div id="component1-name" class="compName"></div>
 *         <div id="component1-seperator" class="compSep"></div>
 *         <span id="component1-value" class="compValue"></span>
 *     </div>
 * </div>
 */
interface Component {
    self: HTMLDivElement;
    image: HTMLImageElement;
    text: HTMLDivElement;
    name: HTMLDivElement;
    seperator: HTMLDivElement;
    value: HTMLSpanElement;
}

/* <----------------------------------------------------------------------------------------------------> */

/**
 * Is in task of doing general work around the components.
 */
class UIHandler {

    readonly TOP_BAR_SIZE = 45;
    readonly COMPONENT_SIZE = 40;

    readonly leftIcon = document.getElementById("left-icon") as HTMLImageElement;
    readonly searchBar = document.getElementById("search-bar") as HTMLInputElement;
    readonly rightSection = document.getElementById("right-section") as HTMLDivElement;
    readonly rightIcon = document.getElementById("right-icon") as HTMLImageElement;

    components = [] as Array<Component>;
    displayedComps = 0;

    /**
     * Creates the components, adds them to the DOM and stores them on the property components.
     *
     * @param max the maximum components the app should have
     */
    constructor(max: number) {
        const body = document.body;

        for (let index = 0; index < max; index++) {
            const newBodyElement = this.#makeElement("div", `component${index+1}`, ["hide"]) as HTMLDivElement;
            const newSubImage = this.#makeElement("img", `component${index+1}-image`, ["compImg"]) as HTMLImageElement;
            const newTextDiv = this.#makeElement("div", `component${index+1}-text`, ["compText"]) as HTMLDivElement;
            const newNameDiv = this.#makeElement("div", `component${index+1}-name`, ["compName"]) as HTMLDivElement;
            const newTextSeperator = this.#makeElement("div", `component${index+1}-seperator`, ["compSep"]) as HTMLDivElement;
            const newTextSpan = this.#makeElement("span", `component${index+1}-value`, ["compValue"]) as HTMLSpanElement;

            newTextDiv.appendChild(newNameDiv);
            newTextDiv.appendChild(newTextSeperator);
            newTextDiv.appendChild(newTextSpan);
            newBodyElement.appendChild(newSubImage);
            newBodyElement.appendChild(newTextDiv);
            body.appendChild(newBodyElement);

            const newComponent: Component = {
                self: newBodyElement,
                image: newSubImage,
                text: newTextDiv,
                name: newNameDiv,
                seperator: newTextSeperator,
                value: newTextSpan,
            }

            this.components.push(newComponent);
        }
    }

    /**
     * Makes an HTML element with an id and classes.
     *
     * @param tagName the html tag you want to create
     * 
     * @param id the id this element should have
     * 
     * @param classes the classes the element should have
     * 
     * @returns an HTMLElement with the id and classes attached.
     */
    #makeElement(tagName: string, id: string, classes: Array<string>): HTMLElement {
        const newElement = document.createElement(tagName);
        newElement.id = id;

        classes.forEach(newClass => {
            newElement.classList.add(newClass);
        });
        
        return newElement;
    }

    /**
     * Resets the UI of the application to the original starting point
     */
    async resetUI(): Promise<void> {
        this.leftIcon.src = await GetImageData("magnifying_glass");
        this.searchBar.value = "";
        this.rightSection.classList.remove("loading-grid");
        this.rightIcon.src = "";
        this.rightIcon.classList.add("hide");

        this.displayComponents(0);
    }

    // #displayComponents unhides the specified amount of components from top to bottom, if any remain they get hidden. It also resizes the window accordingly
    /**
     * Displays the provided amount of components from top to bottom.
     *
     * @param amount how many compnents should be displayed, if the number provided is larger than the max components, it gets reset to that
     */
    displayComponents(amount: number): void {
        if (amount > this.components.length) {
            amount = this.components.length;
        }

        this.displayedComps = 0;

        WindowSetSize(570, this.TOP_BAR_SIZE + (amount * this.COMPONENT_SIZE));

        for (let index = 0; index < this.components.length; index++) {
            if (index + 1 <= amount) {
                this.displayedComps++;

                this.components[index].self.classList.remove("hide");
                this.components[index].self.classList.add("showComp");
            } else {
                this.components[index].self.classList.remove("showComp");
                this.components[index].self.classList.add("hide");
            }  
        }
    }
}