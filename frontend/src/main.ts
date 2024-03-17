import "./assets/styles/main.css";

import {LaunchSearch} from "../wailsjs/go/app/App";
import {EventsOn} from "../wailsjs/runtime/runtime";

// disable right click
document.oncontextmenu = function() {return false;}

const loadingIcon = document.getElementById("loading-icon") as HTMLDivElement;
const resultStatus = document.getElementById("result-status") as HTMLImageElement;
const searchBar = document.getElementById("search-bar") as HTMLInputElement;

// focus the searchBar on load
window.onload = () => {
    searchBar.focus();
};

// makes sure the searchBar is always clicked
document.addEventListener("click", () => {
    searchBar.focus();
});

// send the current input to Go to search the file system
searchBar.addEventListener("input", () => {
    resultStatus.src = "";
    loadingIcon.classList.add("loading-grid");
    LaunchSearch(searchBar.value);
});

// when go found results receive them
EventsOn("searchResult", (results: string[]) => {
    loadingIcon.classList.remove("loading-grid");

    if (results.length > 0) {
        resultStatus.src = "src/assets/images/tick.png";
    } else {
        if (searchBar.value.length > 0) {
            resultStatus.src = "src/assets/images/cross.png";
        } else {
            resultStatus.src = "";
        }
    }

    console.log(results);
});

declare global {
    interface Window {

    }
}
