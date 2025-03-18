
// Fetch reactions from backend
async function FetchReactions(ID, postOrComm, type) {
    try {
        const res = await fetch(`/api/get-reactions?${type}_id=${ID}`)
        if (!res.ok) return

        const { likes, dislikes, userReaction } = await res.json()

        // Update counts
        postOrComm.querySelector(`.${type}-like-count`).textContent = likes
        postOrComm.querySelector(`.${type}-dislike-count`).textContent = dislikes

        // Update the icon states
        UpdateReactionIcons(postOrComm, userReaction, type)

    } catch (err) {
        console.log(err);
    }
}


// Attach event listeners for the like/dislike buttons
function AttachReactionListeners(ID, postOrComm, type) {
    const likeBtn = postOrComm.querySelector(`.${type}-like-button`)
    const dislikeBtn = postOrComm.querySelector(`.${type}-dislike-button`)

    likeBtn.addEventListener("click", (e) => {
        e.stopPropagation()
        AddReaction(ID, "like", postOrComm, type)
    });

    dislikeBtn.addEventListener("click", (e) => {
        e.stopPropagation()
        AddReaction(ID, "dislike", postOrComm, type)
    })
}


// Send reaction updates to server
async function AddReaction(ID, reactionType, postOrComm, type) {
    try {
        const res = await fetch(`/api/add-reaction?type=${type}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                id: ID,
                reaction_type: reactionType,
            }),
        });

        if (res.status == 401) {
            // Not logged in => show auth modal
            document.getElementById("authModal")?.classList.remove("hidden")
            return
        }

        if (!res.ok) {
            return
        }

        // Re-fetch the new counts + userReaction
        FetchReactions(ID, postOrComm, type)
        if (window.location.href.includes("profile") && currentTab == "liked") {
            offset = 0;
            fetchLikedPosts(offset)
                .then(() => {
                    offset += ProfileLimit;
                })
        }
        // Set to refresh when back to profile/liked tab if there is a reaction
        localStorage.setItem("reactionUpdated", "true");

    } catch (err) {
        console.log(err);
        PopError("Something went wrong")
    }
}

// Update the icons based on userReaction.
function UpdateReactionIcons(postOrComm, userReaction, type) {
    const likeIcon = postOrComm.querySelector(`.${type}-like-icon`);
    const dislikeIcon = postOrComm.querySelector(`.${type}-dislike-icon`);

    if (userReaction === "like") {
        likeIcon.src = "../img/likeactive.svg";
        dislikeIcon.src = "../img/dislike.svg";
    } else if (userReaction === "dislike") {
        likeIcon.src = "../img/like.svg";
        dislikeIcon.src = "../img/dislikeactive.svg";
    } else {
        // not active user
        likeIcon.src = "../img/like.svg";
        dislikeIcon.src = "../img/dislike.svg";
    }
}