import { UIHandler } from "../uihandler";
import { bangs } from "../../../res/bang";
import { tlds } from "../../../res/tlds";

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
	readonly websiteRegex = /^(https?:\/\/)?((?![-])[a-z0-9-]+\.)+[a-z]+(\/?)/i;

	bangsMap: Map<string, Array<string>> = new Map;

	#linkComp: number;

	#tlds: Map<string, boolean> = new Map;

	uiHandler!: UIHandler;

	constructor(uiHandler: UIHandler, compIndex: number) {
		this.uiHandler = uiHandler;
		this.#linkComp = compIndex;
		tlds.forEach(tld => {
			this.#tlds.set(tld.toLowerCase(), true);
		});

		bangs.forEach(bang => {
			this.bangsMap.set(bang.t.toLowerCase(), [bang.s, bang.u]);
		});
	}

	/**
	 * Checks if the new input is a link/domain/bang and shows the component, if it is
	 */
	newInput(): void {
		if (this.isBang()) {
			let input = this.uiHandler.searchBar.value.trim().toLowerCase();
			const bangParts = input.split("!");
			const startBang = bangParts[1].split(" ")[0].toLowerCase(); // for start bangs we assume the input starts with an !, so we need to cut the first empty element
			const endBang = bangParts[bangParts.length - 1].toLowerCase();
			let bangInput: string;

			if (this.bangsMap.has(endBang)) {
				bangInput = endBang;
			} else {
				bangInput = startBang;
			}

			this.showComp(this.bangsMap.get(bangInput)![1].replace("{{{s}}}", encodeURIComponent(input.replace(`!${bangInput}`, ""))), `Search with ${this.bangsMap.get(bangInput)![0]}`);
			return;
		}

		if (this.isWebiste()) {
			this.showComp(this.uiHandler.searchBar.value, "Open link in browser");
		}
	}

	/**
	 * Checks if the input starts/ends with a valid duckduckgo bang
	 */
	isBang(): boolean {
		let input = this.uiHandler.searchBar.value.trim().toLowerCase();
		const bangParts = input.split("!");
		const startBang = bangParts[1].split(" ")[0].toLowerCase();  // for start bangs we assume the input starts with an !, so we need to cut the first empty element
		const endBang = bangParts[bangParts.length - 1].toLowerCase();

		if (this.bangsMap.has(startBang) || this.bangsMap.has(endBang)) {
			return true;
		}

		return false;
	}

	/**
	 * Checks if the input is a valid domain/website
	 */
	isWebiste(): boolean {
		let input = this.uiHandler.searchBar.value.trim().toLowerCase();

		if (this.websiteRegex.test(input)) {
			const urlString = input.startsWith("http") ? input : `https://${input}`;
			const parts = new URL(urlString).hostname.split(".");
			const tld = parts[parts.length - 1];
			if (this.#tlds.has(tld)) {
				return true;
			}
		}

		return false;
	}

	/**
	 * modifies the #linkComp to show the suggestion to open the link/domain/bang in the browser
	 */
	showComp(tooltip: string, value: string): void {
		this.uiHandler.components[this.#linkComp].image.src = this.uiHandler.images.get("file") as string;
		this.uiHandler.components[this.#linkComp].tooltip.textContent = tooltip;
		this.uiHandler.components[this.#linkComp].name.textContent = this.uiHandler.searchBar.value;
		this.uiHandler.components[this.#linkComp].value.textContent = value;
		this.uiHandler.components[this.#linkComp].value.classList.add("browserInteract");

		this.uiHandler.displayComponents([this.#linkComp]);
	}
}

export { LinkModule };