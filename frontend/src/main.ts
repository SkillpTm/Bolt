import './style.css';

import {LaunchSearch} from '../wailsjs/go/app/App';
import {EventsOn} from '../wailsjs/runtime/runtime';

// disable right click
document.oncontextmenu = function() {return false;}

const searchBar = document.getElementById("search-bar") as HTMLInputElement;

// focus the searchBar on load
window.onload = function() {
    searchBar.focus();
};

// makes sure the searchBar is always clicked
document.addEventListener("click", function() {
    searchBar.focus();
});

// send the current input to Go to search the file system
searchBar.addEventListener("input", () => {
    LaunchSearch(searchBar.value);
});

// when go found results receive them
EventsOn("searchResult", function(results: string[]) {
    console.log(results)
});

declare global {
    interface Window {

    }
}
