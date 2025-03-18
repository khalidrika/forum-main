// Utility function to calculate relative time
function timeAgo(date) {
    const seconds = Math.floor((new Date() - new Date(date)) / 1000);
    let interval = Math.floor(seconds / 31536000);
    if (interval >= 1) return `${interval}y ago`;

    interval = Math.floor(seconds / 2592000);
    if (interval >= 1) return `${interval}mo ago`;

    interval = Math.floor(seconds / 86400);
    if (interval >= 1) return `${interval}d ago`;

    interval = Math.floor(seconds / 3600);
    if (interval >= 1) return `${interval}h ago`;

    interval = Math.floor(seconds / 60);
    if (interval >= 1) return `${interval}m ago`;

    return `seconds ago`;
}

// Updates all time-ago spans periodically
function updateTimeAgo() {
    const timeAgoSpans = document.querySelectorAll('.time-ago');
    timeAgoSpans.forEach(span => {
        // Get the timestamp from data attribute
        const timestamp = span.getAttribute('data-timestamp');
        if (timestamp) {
            span.textContent = "\u00A0• " + timeAgo(timestamp);
        }
    });
}

// Sorts an array of strings based on whether they start with a query string, 
// prioritizing matches, and then sorts alphabetically.
function sortByQuery(array, query) {
    return array.sort((a, b) => {
        const aLower = a.toLowerCase();
        const bLower = b.toLowerCase();
        const aStarts = aLower.startsWith(query.toLowerCase());
        const bStarts = bLower.startsWith(query.toLowerCase());

        // If a starts with query and b doesn't, a should come first
        if (aStarts && !bStarts) return -1;
        // If b starts with query and a doesn't, b should come first
        if (!aStarts && bStarts) return 1;
        // Otherwise, sort alphabetically
        return aLower.localeCompare(bLower);
    });
}

// Utility function to truncate long posts.
function truncateContent(text, maxLength) {
    if (!text || text.length <= maxLength) return text;
    return text.slice(0, maxLength) + "...";
}
