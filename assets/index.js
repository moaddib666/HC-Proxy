window.onload = function() {
    const element = document.getElementsByTagName("body");
    var a = document.createElement("a")
    a.href = "https://hellcorp.com.ua/";
    var tag = document.createElement("p")
    tag.innerText = "Powered by HellCorp Proxy";
    tag.style = "" +
        "border: solid 1px black;" +
        "-moz-border-radius: 6em;" +
        "-webkit-border-radius: 6em;" +
        "border-radius: 6em;" +
        "background: black;" +
        "color:white;" +
        "padding: 1em;" +
        "font-weight: bold;" +
        "width: fit-content;" +
        "-moz-box-shadow: rgb(0 0 0) 0.1em 0.1em 0.5em;" +
        "-webkit-box-shadow: rgb(0 0 0) 0.1em 0.1em 0.5em;" +
        "box-shadow: rgb(0 0 0) 0.1em 0.1em 0.5em;" +
        "cursor:pointer;" +
        "position: fixed;" +
        "left:1em;" +
        "bottom:1em;" +
        "font-size: 0.6em;";
    a.appendChild(tag);
    element.item(0).appendChild(a);
}
