/* -------------------------
   GLOBAL STATE
--------------------------*/
let selectedTags = [];
let currentTab = "posts"; // or "liked" or "info"
let offset = 0;
// Loading flag to prevent overlapping fetches (throttle)
let isLoading = false;

document.addEventListener("DOMContentLoaded", () => {
    LoadTheme()
    // Add new post modal.
    document.body.insertAdjacentHTML("beforeend", NewPostForm);
    NewPostListener();

    // Prepare dropdown
    const dropdown = document.getElementById("actionsDropdown");
    const toggleButton = document.getElementById("actionsButton");

    toggleButton.addEventListener("click", (e) => {
        e.stopPropagation(); // prevent click from bubbling
        dropdown.classList.toggle("open");
    });

    // If user clicks outside the dropdown, close it
    document.addEventListener("click", (event) => {
        if (!dropdown.contains(event.target)) {
            dropdown.classList.remove("open");
        }
    });

    // Also close dropdown if user clicks anywhere else
    document.addEventListener("click", function () {
        actionsDropdown.classList.remove("open");
    });


    // ========== Tab Buttons ==========
    const dynamicContent = document.getElementById("dynamicContent");
    const actionButtons = document.querySelectorAll(".action-btn");

    actionButtons.forEach((button) => {
        button.addEventListener("click", function () {
            // Remove .active from all
            actionButtons.forEach((btn) => btn.classList.remove("active"));
            // Add .active to the one we clicked
            this.classList.add("active");

            // Clear old content
            dynamicContent.innerHTML = "";

            // Reset global offset & loading for the new tab
            offset = 0;
            isLoading = false;

            // Decide which type we’re loading
            const type = this.getAttribute("data-type").toLowerCase();
            if (type === "liked") {
                currentTab = "liked";
                fetchLikedPosts(offset);
            } else if (type === "posts") {
                currentTab = "posts";
                fetchUserPosts(offset);
            } else if (type === "info") {
                currentTab = "info";
                // "info" won't have infinite scroll
                fetchUserInfo();
            }
        });
    });

    // By default, trigger the "Posts" button
    const defaultInfoBtn = document.querySelector('.action-btn[data-type="posts"]');
    if (defaultInfoBtn) {
        defaultInfoBtn.click();
    }

    // Listen for scroll => infinite loading
    window.addEventListener("scroll", handleProfileScroll, { passive: true });
});

/* -------------------
   SCROLL HANDLER
--------------------*/
function handleProfileScroll() {
    // If "info" tab, do nothing
    if (currentTab === "info") return;

    // Already loading => skip
    if (isLoading) return;

    const windowHeight = window.innerHeight;
    const scrollY = window.scrollY;
    const docHeight = document.documentElement.scrollHeight || document.body.scrollHeight;

    // If close to bottom, fetch more
    if (scrollY + windowHeight >= docHeight - 400) {
        isLoading = true;

        if (currentTab === "posts") {
            fetchUserPosts(offset)
                .then(() => {
                    offset += ProfileLimit; // increment offset
                })
                .finally(() => {
                    isLoading = false;
                });
        } else if (currentTab === "liked") {
            fetchLikedPosts(offset)
                .then(() => {
                    offset += ProfileLimit;
                })
                .finally(() => {
                    isLoading = false;
                });
        }
    }
}

/* -------------------
   FETCH USER INFO
--------------------*/
async function fetchUserInfo() {
    try {
        const res = await fetch("/api/user-info");
        if (!res.ok) {
            throw new Error("Failed to fetch user info");
        }
        const userInfo = await res.json();

        const dynamicContent = document.getElementById("dynamicContent");
        const infoDiv = document.createElement("div");
        infoDiv.classList.add("content-block", "profile-info");

        infoDiv.innerHTML = `
            <div class="profile-info-item">
                <i class="fas fa-user"></i>
                <span class="info-label">Username:</span>
                <span class="info-value">${userInfo.username}</span>
            </div>
            <div class="profile-info-item">
                <i class="fas fa-envelope"></i>
                <span class="info-label">Email:</span>
                <span class="info-value">${userInfo.email}</span>
            </div>
            <div class="profile-info-item">
                <i class="fas fa-calendar-day"></i>
                <span class="info-label">Joined:</span>
                <span class="info-value">${new Date(userInfo.created_at).toLocaleDateString()}</span>
            </div>
        `;

        dynamicContent.appendChild(infoDiv);
    } catch (err) {
        console.error(err);
        DisplayError("errMsg", dynamicContent, "Error loading user info.");
    }
}

/* -------------------
   FETCH USER POSTS
--------------------*/
async function fetchUserPosts(offset) {
    const dynamicContent = document.getElementById("dynamicContent");

    try {
        const res = await fetch(`/api/user-posts?offset=${offset}`);
        if (!res.ok) return;

        const posts = await res.json();
        if ((!posts || posts.length == 0) && offset == 0) {
            dynamicContent.innerHTML = "No posts yet!"
            return
        }
        RenderPosts(posts, offset, 100);
    } catch (err) {
        console.error(err);
        DisplayError("errMsg", dynamicContent, "Error loading user posts.");
    }
}

/* -------------------
   FETCH LIKED POSTS
--------------------*/
async function fetchLikedPosts(offset) {
    const dynamicContent = document.getElementById("dynamicContent");
    try {
        const res = await fetch(`/api/user-liked-posts?offset=${offset}`);
        if (!res.ok) return;

        const posts = await res.json();
        if ((!posts || posts.length == 0) && offset == 0) {
            dynamicContent.innerHTML = "No liked posts!"
            return
        }

        RenderPosts(posts, offset, 100);
    } catch (err) {
        console.error(err);
        DisplayError("errMsg", dynamicContent, "Error loading liked posts.");
    }
}


/* -----------------------
New Post Form listeners 
-------------------------*/
// Create post in dropDown menu.
document.getElementById("createPost")?.addEventListener("click", (e) => {
    e.preventDefault() // to prevent going to # (default)
    document.getElementById("newPostModal")?.classList.remove("hidden");
});


/* -----------------------
    Logout listeners 
-------------------------*/
document.getElementById("logoutBtn")?.addEventListener("click", () => {
    HandleLogout()
})
