export {UIHandler};

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
}