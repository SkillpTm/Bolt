export {UIHandler};

/* <----------------------------------------------------------------------------------------------------> */

import { WindowSetSize } from "../../wailsjs/runtime/runtime";
import { GetImageData } from "../../wailsjs/go/app/App";

/* <----------------------------------------------------------------------------------------------------> */

/*
A component looks like this in html:

<div id="component1" class="hideComp">
    <img id="component1-image" class="compImg">
    <div id="component1-text" class="compText">
        <div id="component1-name" class="compName"></div>
        <div id="component1-seperator" class="compSep"></div>
        <span id="component1-value" class="compValue"></span>
    </div>
</div>
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

// UIHandler updates the UI of the application to the current state
class UIHandler {

    static TOP_BAR_SIZE = 45;
    static COMPONENT_SIZE = 40;

    searchBar = document.getElementById("search-bar") as HTMLInputElement;
    #loadingIcon = document.getElementById("loading-icon") as HTMLDivElement;
    #resultStatus = document.getElementById("result-status") as HTMLImageElement;

    components = [] as Array<Component>;
    #highlightedComp = 0;
    #maxComponents = 0;
    #page = 0;

    // constructor adds as many components to the application as specififed. This will be the max amount for dispalying anything.
    constructor(max: number) {
        this.#maxComponents = max;
        this.#generateMaxComponents(max);
        this.components[0].self.classList.add("highligthed");
    }

    // #generateMaxComponents generates interface Components appended to the application
    #generateMaxComponents(max: number): void {
        const body = document.body;

        for (let index = 0; index < max; index++) {
            const newBodyElement = this.#makeElement("div", `component${index+1}`, ["hideComp"]) as HTMLDivElement;
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

    // #makeElement makes an HTML element with an id and classes
    #makeElement(tagName: string, id: string, classes: Array<string>): HTMLElement {
        const newElement = document.createElement(tagName);
        newElement.id = id;
        classes.forEach(newClass => {
            newElement.classList.add(newClass);
        });
        
        return newElement;
    }

    // startedSearch updates the loading icon
    startedSearch(): void {
        this.#resultStatus.src = "";

        this.#loadingIcon.classList.add("loading-grid");
    }

    // getCurrentComp gets you the index of the currently highligthed component relative to the results on the stateHandler
    getCurrentComp(): number {
        return this.#page * this.#maxComponents + this.#highlightedComp;
    }

    // getCurrentComp gets you the index of the currently hovered over component relative to the results on the stateHandler
    getHoverComp(comp: Component): number {
        return this.#page * this.#maxComponents + (parseInt(comp.self.id[9]) -1);
    }

    updateHighlightedComp(change: number): void {
        this.components[this.#highlightedComp].self.classList.remove("highligthed");

        if (change === 0) {
            this.#highlightedComp = 0;
        } else if (change < 0) {
            this.#highlightedComp = (this.#highlightedComp + change + this.components.length) % this.components.length;
        } else {
            this.#highlightedComp = (this.#highlightedComp + change) % this.components.length;
        }

        this.components[this.#highlightedComp].self.classList.add("highligthed");
    }

    // #displayComponents unhides the specified amount of components from top to bottom, if any remain they get hidden. It also resizes the window accordingly
    #displayComponents(amount: number): void {
        if (amount > this.#maxComponents) {
            amount = this.#maxComponents;
        }

        WindowSetSize(570, UIHandler.TOP_BAR_SIZE + (amount * UIHandler.COMPONENT_SIZE));

        for (let index = 0; index < this.components.length; index++) {
            if (index + 1 <= amount) {
                this.components[index].self.classList.remove("hideComp");
                this.components[index].self.classList.add("showComp");
            } else {
                this.components[index].self.classList.remove("showComp");
                this.components[index].self.classList.add("hideComp");
            }  
        }

        this.updateHighlightedComp(0);
    }

    // reset resets the UI of the application to the original starting point
    reset(): void {
        this.#loadingIcon.classList.remove("loading-grid");
        this.#resultStatus.src = "";
        this.searchBar.value = "";
        this.#page = 0;
        this.#displayComponents(0);
    }

    /**
     * Updates the page number. If the change would be invalid (i.e. below 0) nothing changes.
     *
     * @param change increase/decrease to page count, 0 resets it back to 0
     * 
     * @param resultsLength the length of the results from the StateHandler, used to determine if the change is valid
     */
    updatePage(change: number, resultsLength: number): void {
        if (change === 0) {
            this.#page = 0;
            return;
        }

        if ((this.#page + change) * this.#maxComponents > resultsLength-1) {
            return;
        }

        if ((this.#page + change) === 0) {
            return;
        }

        this.#page += change;
    }

    // displayResults adds the entry names, paths and icons to the components and displays them
    async displayResults(queryLength: number, results: Array<String>): Promise<void> {
        this.#loadingIcon.classList.remove("loading-grid");

        if (queryLength === 0) {
            this.reset();
        } else {
            if (results.length > 0) {
                this.#resultStatus.src = await GetImageData("tick");
            } else {
                this.#resultStatus.src = await GetImageData("cross");
            }

            let displayAmount = 0;

            for (let index = 0; index < results.length && index < this.#maxComponents; index++) {
                const currentFile = index + (this.#maxComponents * this.#page);

                if (currentFile > results.length-1) {
                    break;
                }

                let filePath = results[currentFile].split("/");
    
                // if the last element is empty, it means our string ended in a slash, indicating it was a folder
                if (filePath[filePath.length-1] === "") {
                    filePath.pop();
                    this.components[index].image.src = await GetImageData("folder");
                } else {
                    this.components[index].image.src = await GetImageData("file");
                }
    
                this.components[index].name.textContent = filePath.pop() as string;
                this.components[index].value.textContent = filePath.join("/") + "/";

                displayAmount++;
            }

            this.#displayComponents(displayAmount);
        }
    }
}