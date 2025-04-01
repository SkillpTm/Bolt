export { UIHandler, type Component };

import { WindowSetSize } from "../../wailsjs/runtime/runtime";

/**
 * Section of the UI below the search bar.
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
 * 
 * @param self component wrapper
 * 
 * @param image component image on the right side
 * 
 * @param text wrapper for name, sperator and value
 * 
 * @param name component name, left of seperator
 * 
 * @param seperator small line between text and name
 * 
 * @param value component value, right of seperator
 */
interface Component {
    self: HTMLDivElement;
    image: HTMLImageElement;
    text: HTMLDivElement;
    name: HTMLDivElement;
    seperator: HTMLDivElement;
    value: HTMLSpanElement;
}

/**
 * Holds the basic properties and functions to manipulate the UI
 * 
 * @param TOP_BAR_SIZE pixel size of the top bar
 * 
 * @param COMPONENT_SIZE standardised pixel size of a component
 * 
 * @param leftIcon html element for the left icon of the top bar
 * 
 * @param searchBar html element for the search bar of the top bar
 * 
 * @param rightSection html element for the right section of the top bar
 * 
 * @param rightIcon html element for the right icon of the top bar
 * 
 * @param components array of all components, even the hidden ones
 * 
 * @param displayedComps how many components are currently visible
 * 
 * @param images map of the base64 image data needed to embed for the html
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
    images = new Map<string, string>();

    /**
     * Creates the components, adds them to the DOM and stores them on the property components.
     *
     * @param max the maximum components the app should be able to display, minimum 3
     */
    constructor(max: number) {
        if (max < 3) {
            max = 3;
        }

        this.#regenerateComponents(max);
    }

    /**
     * Makes an HTML element with an id and classes.
     *
     * @param max the amount components the app should be able to display
     */
    #regenerateComponents(max: number): void {
        this.components = [] as Array<Component>;

        for (let index = 0; index < max; index++) {
            const newWrapper = this.#makeElement("div", `component${index+1}`, ["hide"]) as HTMLDivElement;
            const newSubImage = this.#makeElement("img", `component${index+1}-image`, ["compImg"]) as HTMLImageElement;
            const newTextDiv = this.#makeElement("div", `component${index+1}-text`, ["compText"]) as HTMLDivElement;
            const newNameDiv = this.#makeElement("div", `component${index+1}-name`, ["compName"]) as HTMLDivElement;
            const newTextSeperator = this.#makeElement("div", `component${index+1}-seperator`, ["compSep"]) as HTMLDivElement;
            const newTextSpan = this.#makeElement("span", `component${index+1}-value`, ["compValue"]) as HTMLSpanElement;

            newTextDiv.appendChild(newNameDiv);
            newTextDiv.appendChild(newTextSeperator);
            newTextDiv.appendChild(newTextSpan);
            newWrapper.appendChild(newSubImage);
            newWrapper.appendChild(newTextDiv);
            document.body.appendChild(newWrapper);

            this.components.push(
                {
                    self: newWrapper,
                    image: newSubImage,
                    text: newTextDiv,
                    name: newNameDiv,
                    seperator: newTextSeperator,
                    value: newTextSpan,
                } as Component
            );
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
     * @returns HTMLElement with the id and classes attached.
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
    resetUI(): void {
        this.leftIcon.src = this.images.get("magnifying_glass") as string;
        this.searchBar.value = "";
        this.rightSection.classList.remove("loading-grid");
        this.rightIcon.src = "";
        this.rightIcon.classList.add("hide");

        this.displayComponents(0);
    }

    /**
     * Displays the provided amount of components from top to bottom.
     *
     * @param amount how many compnents should be displayed, if the number provided is larger than the property components length, it gets set to that
     */
    displayComponents(amount: number): void {
        if (amount > this.components.length) {
            amount = this.components.length;
        }

        this.displayedComps = amount;

        WindowSetSize(570, this.TOP_BAR_SIZE + (amount * this.COMPONENT_SIZE));

        for (let index = 0; index < this.components.length; index++) {
            if (index + 1 <= amount) {
                this.components[index].self.classList.remove("hide");
                this.components[index].self.classList.add("showComp");
            } else {
                this.components[index].self.classList.remove("showComp");
                this.components[index].self.classList.add("hide");
            }  
        }
    }
}