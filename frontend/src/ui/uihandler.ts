import { GetImageData } from "../../wailsjs/go/app/App";

/**
 * Section of the UI below the search bar.
 *
 * A component looks like this in html:
 * <div id="component1" class="hide tooltip">
 *     <img id="component1-image" class="compImg">
 *     <div id="component1-text" class="compText">
 *         <div id="component1-name" class="compName"></div>
 *         <div id="component1-seperator" class="compSep"></div>
 *         <span id="component1-value" class="compValue"></span>
 *         <span id="component1-tooltip" class="tooltiptext"></span>
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
 * 
 * @param tooltip component to hold the tooltip text
 */
class Component {
	self: HTMLDivElement;

	image: HTMLImageElement;

	text: HTMLDivElement;

	name: HTMLDivElement;

	seperator: HTMLDivElement;

	value: HTMLSpanElement;

	tooltip: HTMLSpanElement;

	constructor(index: number) {
		const newWrapper = this.#makeElement("div", `component${index}`, ["hide", "tooltip"]) as HTMLDivElement;
		const newSubImage = this.#makeElement("img", `component${index}-image`, ["compImg"]) as HTMLImageElement;
		const newTextDiv = this.#makeElement("div", `component${index}-text`, ["compText"]) as HTMLDivElement;
		const newNameDiv = this.#makeElement("div", `component${index}-name`, ["compName"]) as HTMLDivElement;
		const newTextSeperator = this.#makeElement("div", `component${index}-seperator`, ["compSep"]) as HTMLDivElement;
		const newTextSpan = this.#makeElement("span", `component${index}-value`, ["compValue"]) as HTMLSpanElement;
		const newToolTipSpan = this.#makeElement("span", `component${index}-tooltip`, ["tooltiptext"]) as HTMLSpanElement;

		newTextDiv.appendChild(newNameDiv);
		newTextDiv.appendChild(newTextSeperator);
		newTextDiv.appendChild(newTextSpan);
		newTextDiv.appendChild(newToolTipSpan);
		newWrapper.appendChild(newSubImage);
		newWrapper.appendChild(newTextDiv);
		document.body.appendChild(newWrapper);

		this.self = newWrapper;
		this.image = newSubImage;
		this.text = newTextDiv;
		this.name = newNameDiv;
		this.seperator = newTextSeperator;
		this.value = newTextSpan;
		this.tooltip = newToolTipSpan;
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
	 * gets the index added to the id of the component
	 * 
	 * @returns index of the component
	 */
	getIndex(): number {
		return parseInt(this.self.id[9]);
	}
}

/**
 * Holds the basic properties and functions to manipulate the UI
 * 
 * @param topBarHeight height of the top bar in pixels
 * 
 * @param componentHeight standardised pixel height of a component
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
 * @param images map of the base64 image data needed to embed for the html
 */
class UIHandler {
	readonly topBarHeight = 45;

	readonly componentHeight = 40;

	readonly leftIcon = document.getElementById("left-icon") as HTMLImageElement;

	readonly searchBar = document.getElementById("search-bar") as HTMLInputElement;

	readonly rightSection = document.getElementById("right-section") as HTMLDivElement;

	readonly rightIcon = document.getElementById("right-icon") as HTMLImageElement;

	components: Array<Component> = [];

	images: Map<string, string> = new Map();

	constructor(max: number) {
		(async () => {
			const temp: Record<string, string> = await GetImageData();
			this.images = new Map(Object.entries(temp));
		})();

		if (max < 4) {
			max = 4;
		}

		for (let index = 0; index < max; index++) {
			this.components.push(new Component(index));
		}
	}

	/**
	 * gets an array of all currently displayed components
	 * 
	 * @returns an array of all indexes of all shown components
	 */
	getDisplayedComps(): Array<number> {
		const output: Array<number> = [];

		this.components.forEach(comp => {
			if (comp.self.classList.contains("show")) {
				output.push(comp.getIndex());
			}
		});

		return output;
	}

	/**
	 * Gets the index added to the id of the highlighted component. If there is none this sets it to be the first visible one.
	 * 
	 * @returns index of the highlighted component
	 */
	getHighlightedComp(): number {
		for (let index = 0; index < this.components.length; index++) {
			if (this.components[index].self.classList.contains("highlighted")) {
				return index;
			}

		}

		// if we didn't find a highlightedComp we reset it back to the first comp
		this.components[0].self.classList.add("highlighted");
		return 0;
	}

	/**
	 * Resets the UI of the application to the original starting point
	 */
	resetUI(): void {
		this.leftIcon.src = this.images.get("bolt") as string;
		this.searchBar.value = "";
		this.rightSection.classList.remove("loading-grid");
		this.rightIcon.src = "";
		this.rightIcon.classList.add("hide");

		this.displayComponents(undefined, Array.from({ length: 8 }, (_, i) => i));
	}

	/**
	 * Displays/hides the provided components, all other components stay unchanged.
	 * 
	 * @param showComps which components to change the state to be shown
	 * 
	 * @param hideComps which components to change the state to be hidden
	 */
	displayComponents(showComps: Array<number> = [], hideComps: Array<number> = []): void {
		showComps.forEach(index => {
			this.components[index].self.classList.remove("hide");
			this.components[index].self.classList.add("show");
		});

		hideComps.forEach(index => {
			this.components[index].self.classList.remove("show");
			this.components[index].self.classList.add("hide");
		});
	}

	/**
	 * Updates the higlighted componented by 1 in either direction (with overflow) or resets back to the first dispalyed comp.
	 *
	 * @param increase if the index of the highligthed comp is supposed to increase
	 * 
	 * @param reset resets back to the first displayed component
	 */
	updateHighlightedComp(increase?: boolean, reset: boolean = false): void {
		const displayedComps = this.getDisplayedComps();
		const highlightedComp = this.getHighlightedComp();

		if (displayedComps.length == 0) {
			return;
		}

		if (reset || displayedComps.length == 1 || displayedComps.indexOf(highlightedComp) < 0) {
			this.components[highlightedComp].self.classList.remove("highlighted");
			this.components[displayedComps[0]].self.classList.add("highlighted");
			return;
		}

		this.components[highlightedComp].self.classList.remove("highlighted");
		const current = displayedComps.indexOf(highlightedComp) as number;

		if (increase) {
			this.components[displayedComps[current + 1 <= displayedComps.length - 1 ? current + 1 : 0]].self.classList.add("highlighted");
		} else {
			this.components[displayedComps[current - 1 >= 0 ? current - 1 : displayedComps.length - 1]].self.classList.add("highlighted");
		}
	}
}

export { Component, UIHandler };