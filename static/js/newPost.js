// Handle new Post creation.
async function HandleNewPost(e, categories, clearTags) {
    e.preventDefault();
    RemoveError("postErrorMsg", e.target)

    const newPostModal = document.getElementById("newPostModal");

    const titleInput = document.getElementById("formPostTitle");
    const contentInput = document.getElementById("formPostContent");

    let title = titleInput.value.trim();
    let content = contentInput.value.trim();

    if (CheckPostSize(content.length, title.length, e.target)) {
        return
    }

    try {
        const formData = new FormData();
        formData.append("title", title)
        formData.append("content", content)
        formData.append("categories", JSON.stringify(categories))

        const res = await fetch("/api/create-post", {
            method: "POST",
            body: formData,
        });

        if (!res.ok) {
            const errData = await res.json()
            DisplayError("postErrorMsg", e.target, errData.msg);
        } else {
            // Hide modal
            newPostModal.classList.add("hidden");

            // Clear form fields
            titleInput.value = "";
            contentInput.value = "";
            clearTags;

            if (window.location.href.includes("profile") && currentTab == "liked") {               
                offset = 0;
                fetchLikedPosts(offset)
                .then(() => {
                    offset += ProfileLimit;
                })
            } else if (window.location.href.includes("profile") && currentTab == "posts") {
                offset = 0;
                fetchUserPosts(offset)
                .then(() => {
                    offset += ProfileLimit; // increment offset
                })
            } else if (window.location.pathname === "/") { // Don't work in post page
                offset = 0;
                LoadPosts(offset, selectedTags.join(",")).then(() => {
                    offset += HomeLimit;
                });               
            }
        }
        localStorage.setItem("postCreated", "true");
    } catch (err) {
        console.log(err);
        PopError("Something went wrong")
    }
}

// Implement front-end limits on title/content size.
function CheckPostSize(contentLength, titleLength, elem) {
    if (contentLength == 0 || titleLength == 0) {
        DisplayError("postErrorMsg", elem, "Title and content are required");
        return true
    }
    if (contentLength < 6) {
        DisplayError("postErrorMsg", elem, "Post content is too short");
        return true;
    }
    if (titleLength < 4) {
        DisplayError("postErrorMsg", elem, "Post title is too short");
        return true;
    }
    if (contentLength > 8000) {
        DisplayError("postErrorMsg", elem, "Post content is too large");
        return true;
    }
    if (titleLength > 500) {
        DisplayError("postErrorMsg", elem, "Post title is too large");
        return true;
    }

    return false;
}
