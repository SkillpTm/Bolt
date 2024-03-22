export {UIHandler};

/* <----------------------------------------------------------------------------------------------------> */

import { WindowSetSize } from "../../wailsjs/runtime/runtime";

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

    // constructor adds as many components to the application as specififed. This will be the max amount for dispalying anything.
    constructor(max: number) {
        this.#generateMaxComponents(max);
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

    // #displayComponents unhides the specified amount of components from top to bottom, if any remain they get hidden. It also resizes the window accordingly
    #displayComponents(amount: number): void {
        if (amount > this.components.length) {
            amount = this.components.length
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
    }

    // reset resets the UI of the application to the original starting point
    reset(): void {
        this.#loadingIcon.classList.remove("loading-grid");
        this.#resultStatus.src = "";
        this.#displayComponents(0);
    }

    // displayResults addes the entry names, paths and icons to the components and displays them
    displayResults(queryLength: number, results: Array<String>): void {
        this.#loadingIcon.classList.remove("loading-grid");

        if (queryLength === 0) {
            this.reset();
        } else {
            if (results.length > 0) {
                this.#resultStatus.src = "src/assets/images/tick.png";
            } else {
                this.#resultStatus.src = "src/assets/images/cross.png";
            }

            for (let index = 0; index < results.length && index < 7; index++) {
                let fs = results[index].split("/");
    
                // if the last element is empty, it means our string ended in a slash, indicating it was a folder
                if (fs[fs.length-1] === "") {
                    fs.pop();
                    this.components[index].image.src = "src/assets/images/folder.png";
                } else {
                    this.components[index].image.src = "src/assets/images/file.png";
                }
    
                this.components[index].name.textContent = fs.pop() as string;
                this.components[index].value.textContent = fs.join("/") + "/";
            }

            this.#displayComponents(results.length);
        }
    }
}