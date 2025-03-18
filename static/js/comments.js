let commentOffset = 0;

async function FetchComments(postID, isLoadMore = false) {
    try {
        const res = await fetch(`/api/get-comments?post_id=${postID}&offset=${commentOffset}`);
        if (!res.ok) {
            return
        }
        const comments = await res.json()

        RenderComments(comments, isLoadMore)

        if (!comments || comments.length < 20) {
            document.getElementById("showMoreBtn").style.display = "none";
        } else {
            document.getElementById("showMoreBtn").style.display = "inline-block";
        }

    } catch (err) {
        console.error(err)
        PopError("Something went wrong")
    }
}

function RenderComments(comments, isLoadMore) {
    const commentsList = document.getElementById("commentsList")

    // Only clear the list on initial load
    if (!isLoadMore) {
        commentsList.innerHTML = ""
    }

    if (!comments || comments.length === 0) {
        if (!isLoadMore) {
            commentsList.innerHTML = `<p>No comments yet. Be the first to comment!</p>`
        }
        return
    }

    comments.forEach((comment) => {
        const profilePic = comment.profilePic || "../img/avatar.webp";
        const commentEl = document.createElement("div")
        commentEl.classList.add("comment-item")

        commentEl.innerHTML = `
            <p class="comment-meta">
                <div class="username">
                    <img src="${profilePic}" alt="User Avatar" class="comment-user-avatar">
                    ${comment.username}
                    <span class="time-ago" data-timestamp="${comment.created_at}">&nbsp• ${timeAgo(comment.created_at)}</span>
                </div>
            </p>
            <p class="comment-content">${comment.content}</p>

            <!-- Reaction Section for each comment -->
            <div class="comment-reaction-buttons">
                <button class="reaction-button comment-like-button">
                    <img class="comment-like-icon" src="../img/like.svg" alt="like">
                </button>
                <span class="comment-like-count">0</span>

                <button class="reaction-button comment-dislike-button">
                    <img class="comment-dislike-icon" src="../img/dislike.svg" alt="dislike">
                </button>
                <span class="comment-dislike-count">0</span>
            </div>
        `

        // Always append to the bottom of the list
        commentsList.appendChild(commentEl)

        FetchReactions(comment.id, commentEl, "comment")
        AttachReactionListeners(comment.id, commentEl, "comment")
    })
}

async function AddComment(postID, content) {
    const commContainer = document.getElementById("addCommentContainer")

    try {
        const res = await fetch("/api/add-comment", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                id: parseInt(postID, 10),
                content,
            })
        })

        if (res.status === 401) {
            document.getElementById("authModal").classList.remove("hidden")
            return
        }

        if (!res.ok) {
            RemoveError("commentErrMsg", commContainer)
            const errData = await res.json()
            DisplayError("commentErrMsg", commContainer, errData.msg)
            return
        }

        // Clear input and any error message
        RemoveError("commentErrMsg", commContainer)
        document.getElementById("commentContent").value = ""

        // Reset offset and fetch fresh comments
        commentOffset = 0
        await FetchComments(postID, false)

        // Scroll up to the comments section
        document.querySelector(".reaction-buttons").scrollIntoView({ behavior: "smooth", block: "start" });

    } catch (err) {
        console.log(err);
        PopError("Something went wrong")
    }
}

// Update all time-ago spans every minute
document.addEventListener('DOMContentLoaded', () => {
    setInterval(updateTimeAgo, 60000);
});