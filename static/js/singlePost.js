document.addEventListener("DOMContentLoaded", async () => {
    const params = new URLSearchParams(window.location.search)
    const postID = params.get("post_id")
    const postDiv = document.getElementById("postDiv")

    // Fetchers
    if (postID) {
        await FetchFullPost(postID)
        FetchReactions(postID, postDiv, "post")
        FetchComments(postID)
    } else {
        console.log("No post found");
        PopError("Something went wrong")
        return
    }

    // Comment Listener
    const commentForm = document.getElementById("commentForm");
    const commentSubmit = document.getElementById("commentSubmit");
    commentForm.addEventListener("submit", (e) => {
        e.preventDefault();
        commentSubmit.disabled = true;

        setTimeout(() => {
            commentSubmit.disabled = false;
        }, 1000);

        const content = document.getElementById("commentContent").value;
        if (!content) return;
        AddComment(postID, content);
    });

    // Listen for "Show more" button click:
    document.getElementById("showMoreBtn").addEventListener("click", () => {
        commentOffset += 20;
        FetchComments(postID, true); // pass your postID variable
    });
})

// Send fetch to backEnd with post ID
async function FetchFullPost(id) {
    try {
        const res = await fetch(`/api/get-singlePost?post_id=${encodeURIComponent(id)}`);
        if (!res.ok) {
            const errData = await res.json();
            console.error("Failed to fetch post:", errData.msg);
            PopError("Something went wrong")
            return;
        }
        const post = await res.json()

        RenderPost(post, postDiv, "single")
        longTagNames()

    } catch (err) {
        console.log(err);
        PopError("Something went wrong")
    }
}

document.getElementById('backLink').addEventListener('click', (e) => {
    e.preventDefault(); 
    history.back();
});
