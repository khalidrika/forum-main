let selectedTags = []; // Global array of current tags
let offset = 0;
let isLoading = false;
let lastScrollY = window.scrollY;

// Initial load of unfiltered posts
LoadPosts(offset).then(() => {
    offset += HomeLimit;
});

window.addEventListener('scroll', handleScroll, { passive: true });

async function LoadPosts(offset = 0, tagsStr = "") {
    try {
        let url = `/api/get-posts?offset=${offset}`;
        if (tagsStr) {
            url += `&tags=${encodeURIComponent(tagsStr)}`;
        }

        const res = await fetch(url);
        if (!res.ok) {
            const errData = await res.json()
            console.log(errData.msg);
            PopError("Something went wrong. Failed to get posts.");
            return;
        }
        const posts = await res.json();

        // If offset == 0, we assume we're refreshing from start
        RenderPosts(posts, offset);

    } catch (err) {
        console.error(err);
        PopError("Something went wrong. Network error.");
    }
}

// Check scroll position on each "scroll" event.
// If close to bottom, load more posts (lazy loading).
function handleScroll() {
    // If we're already loading, skip
    if (isLoading) return;

    const windowHeight = window.innerHeight; // Related to device
    const scrollY = window.scrollY; // How much scrolled down through document
    const docHeight = document.documentElement.scrollHeight || document.body.scrollHeight;

    // Detect strict scroll direction
    let scrollingDown = scrollY > lastScrollY;
    lastScrollY = scrollY; // Update position

    // Only trigger loading when strictly scrolling down & near the bottom
    if (scrollingDown && scrollY + windowHeight >= docHeight - 400) { // If remaining bellow <= 400px
        isLoading = true;

        LoadPosts(offset, selectedTags.join(","))
            .then(() => {
                offset += HomeLimit;
            })
            .catch((err) => {
                console.log(err);
                PopError("Something went wrong");
            })
            .finally(() => {
                isLoading = false;
            });
    }
}
