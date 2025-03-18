// Renders a list of post objects into #postsContainer.
// if clearFirst we first remove content inside it.
function RenderPosts(posts, offset, truncate = 300) {
    let postsContainer = document.getElementById("postsContainer");
    if (!postsContainer) {
        postsContainer = document.getElementById('dynamicContent');
        if (!postsContainer) return
    }

    // If offset = 0 clear the container
    if (offset == 0) {
        postsContainer.innerHTML = "";
    }

    if (!posts || posts.length === 0) {
        return;
    }

    posts.forEach((post) => {
        const postDiv = document.createElement("div");
        postDiv.classList.add("post-card");
        postDiv.setAttribute("data-id", post.id);

        post.content = truncateContent(post.content, truncate)

        RenderPost(post, postDiv)
        postsContainer.appendChild(postDiv);

        // e.g. single-post click
        postDiv.addEventListener("click", (e) => {
            if (e.target.closest(".reaction-buttons") || e.target.closest(".post-title")) return;
            // Then go to single post
            window.location.href = `/post?post_id=${post.id}`;
        });

        FetchReactions(post.id, postDiv, "post")
    });
    longTagNames()
}

// Build post html (user in getPosts and singlePost)
function RenderPost(post, postDiv, single = "") {
    const profilePic = post.profilePic || "../img/avatar.webp";

    // Build categories HTML
    const categoriesHTML = (post.categories || [])
        .map(
            (cat) => `
                <div class="post-tags">
                    <img src="../img/tag-icon.svg" alt="Tag icon" class="tag-icon">
                    <span>${cat.name}</span>
                </div>`
        )
        .join("");

    postDiv.innerHTML = `
        <!-- Post Header -->
        <div class="post-header">
            <img src="${profilePic}" alt="User Avatar" class="user-avatar">
            <div class="user-info">
                <div class="username">
                    ${post.username}
                    <span class="time-ago" data-timestamp="${post.created_at}">&nbsp• ${timeAgo(post.created_at)}</span>
                </div>
            </div>
            <div class="post-tags-container">${categoriesHTML}</div>
        </div>

        <!-- Post Content -->
        <div class="post-content" id="${single}postContent">
            <h1 id="${single}postTitle" class="post-title"><a href="/post?post_id=${post.id}" class="post-header-link">${post.title}</a></h1>
            <p style="white-space: pre-wrap">${post.content}</p>
        </div>

        <!-- Reaction Buttons -->
        <div class="reaction-buttons">
            <div class="reaction-item">
                <button class="reaction-button post-like-button">
                    <img class="post-like-icon" src="../img/like.svg" alt="like">
                </button>
                <span class="post-like-count">0</span>
            </div>

            <div class="reaction-item">
                <button class="reaction-button post-dislike-button">
                    <img class="post-dislike-icon" src="../img/dislike.svg" alt="dislike">
                </button>
                <span class="post-dislike-count">0</span>
            </div>
        </div>


        <!-- Comments Counts -->
        <div class="reaction-item" id="${single}comments-counts">
            <button class="reaction-button comment-button">
                <img id="comment-icon" src="../img/comments.png" alt="Comments">
            </button>
            <span class="comments-count">0</span>
        </div>
    `;

    AttachReactionListeners(post.id, postDiv, "post")
    FetchCommentsCount(post.id, postDiv)
    updateTagIcons()
}


// Fetch comments count
async function FetchCommentsCount(postId, postDiv) {
    try {
        const res = await fetch(`/api/comments-count?post_id=${postId}`);
        if (!res.ok) {
            const errData = await res.json();
            console.error(errData.msg);
            PopError("Something went wrong")
            return;
        }
        const data = await res.json();
        const commentsCount = postDiv.querySelector(".comments-count");
        commentsCount.textContent = data.count || 0;
    } catch (err) {
        console.error(err);
    }
}

// Make sure very longs tags don't break the style
function longTagNames() {
    document.querySelectorAll(".post-tags-container").forEach(tagsContainer => {
        // Check if the width of post-tags-container is greater than 300px
        if (tagsContainer.offsetWidth > 500) {
            tagsContainer.style.position = "relative";
            tagsContainer.style.top = "auto";
            tagsContainer.style.right = "auto";
            tagsContainer.style.marginLeft = "auto";
            tagsContainer.style.display = "flex";
            tagsContainer.style.justifyContent = "right";
            tagsContainer.style.maxWidth = "500px";
        }
    });
}