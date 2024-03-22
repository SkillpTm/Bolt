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