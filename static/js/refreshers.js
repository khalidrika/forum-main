// Utility function to refresh all reactions
function RefreshReactions() {
    // Get all elements that have data-id
    const allPostDivs = document.querySelectorAll(".post-card[data-id]");

    allPostDivs.forEach((postEl) => {
        // Grab the post ID from the data-attribute
        const postID = postEl.getAttribute("data-id");
        if (postID) {
            FetchReactions(postID, postEl, "post");
        }
    });
}

// Utility function to refresh comment counts for all posts
function RefreshCommentsCounts() {
    // Get all elements that have data-id
    const allPostDivs = document.querySelectorAll(".post-card[data-id]");

    allPostDivs.forEach((postEl) => {
        // Grab the post ID from the data-attribute
        const postID = postEl.getAttribute("data-id");
        if (postID) {
            FetchCommentsCount(postID, postEl);
        }
    });
}

// Refresh reactions, comments and posts on each page show (no-reload)
window.addEventListener("pageshow", async (e) => {
    preLoadTheme()
    // The page was restored from BFCache
    if (e.persisted) { 
        // User loggedout and he's trying to get back to profile.
        if  (window.location.href.includes("profile") && localStorage.getItem("loggedOut") === "true") {
            window.location.href = "/";
            localStorage.removeItem("loggedOut");
        }

        RefreshReactions()
        RefreshCommentsCounts()
        if (localStorage.getItem("postCreated") === "true") {
            if (window.location.href.includes("profile") && currentTab == "posts") {
                offset = 0;
                fetchUserPosts(offset)
                    .then(() => {
                        offset += ProfileLimit;
                    })
            } else if (window.location.pathname === "/") {
                offset = 0;
                LoadPosts(offset, selectedTags.join(",")).then(() => {
                    offset += HomeLimit;
                });
            }
            localStorage.removeItem("postCreated"); // clear the marker
        }
        if (localStorage.getItem("reactionUpdated") === "true" && window.location.href.includes("profile") && currentTab == "liked") {
            offset = 0;
            fetchLikedPosts(offset)
                .then(() => {
                    offset += ProfileLimit;
                })
            localStorage.removeItem("reactionUpdated"); // clear the marker
        }
    }
});
