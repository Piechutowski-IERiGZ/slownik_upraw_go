"use strict";
function toggleCollapse(event) {
    const container = event.parentElement.children[1];
    if (container.classList.contains("hidden")) {
        // Open
        container.classList.remove("hidden", "max-h-0");
        container.classList.add("flex", "max-h-fit");
        event.textContent = "📂" + event.textContent.slice(2);
    }
    else {
        // Close
        container.classList.remove("flex", "max-h-fit");
        container.classList.add("hidden", "max-h-0");
        event.textContent = "📁" + event.textContent.slice(2);
    }
}
//# sourceMappingURL=collapsible.js.map