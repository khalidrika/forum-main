/***********************************************************
*          LoggedIn NavBar (Login-button listener)         *
************************************************************/
function NavBarListener() {
    // Open Auth Modal when login is clicked.
    const openAuthBtn = document.getElementById("openAuthBtn");
    openAuthBtn?.addEventListener("click", () => {
        document.getElementById("authModal").classList.remove("hidden");
        document.getElementById("loginContainer").classList.remove("hidden");
        document.getElementById("signUpContainer").classList.add("hidden");
    });

    // Listener on profile link to reset cooldown counter
    document.querySelector(".profile-link")?.addEventListener("click", () => {
        // Reset cooldown counter
        localStorage.setItem("reloadCount", 0);
    })

    // Create Post and Logout listeners.
    AddPostListener()
    document.getElementById("logoutBtn")?.addEventListener("click", (e) => {
        e.preventDefault()
        HandleLogout()
    })
}

/************************************
*      Authentification Modals      *
*************************************/
function AuthListener() {
    // Close Modals (when we click on closeModalBtn)
    const authModal = document.getElementById("authModal");
    const closeModalBtn = document.getElementById("closeModalBtn");
    closeModalBtn?.addEventListener("click", () => {
        authModal.classList.add("hidden");
        ResetModal();
    });

    // Close when we click outside (inside authModal).
    window.addEventListener("click", (e) => {
        if (e.target === authModal) {
            authModal.classList.add("hidden");
            ResetModal();
        }
    });

    // Switch between forms.
    const login = document.getElementById("loginContainer");
    const signUp = document.getElementById("signUpContainer");
    document.getElementById("showSignUpLink")?.addEventListener("click", (e) => {
        e.preventDefault();
        login.classList.add("hidden");
        signUp.classList.remove("hidden");
    });
    document.getElementById("showLoginLink")?.addEventListener("click", (e) => {
        e.preventDefault();
        signUp.classList.add("hidden");
        login.classList.remove("hidden");
    });
}

// Listen for inputs in Log In form and trigger login operation.
function LoginFormListener() {
    const loginForm = document.getElementById("loginForm");
    const loginEmail = document.getElementById("loginEmail");
    const loginPassword = document.getElementById("loginPassword");
    const loginSubmit = document.getElementById("loginSubmit");

    loginEmail?.addEventListener("input", validateLogin);
    loginPassword?.addEventListener("input", validateLogin);
    loginForm?.addEventListener("submit", HandleLogin);

    function validateLogin() {
        const valid = loginEmail.value && loginPassword.value;
        loginSubmit.disabled = !valid; // JS
        loginSubmit.classList.toggle("disabled", !valid); // CSS
    }
}

// Listen for inputs in Sign Up form and trigger login operation.
function SignUpFormListener() {
    const signUpForm = document.getElementById("signUpForm");
    const signUpEmail = document.getElementById("signUpEmail");
    const signUpUsername = document.getElementById("signUpUsername");
    const signUpPassword = document.getElementById("signUpPassword");
    const signUpSubmit = document.getElementById("signUpSubmit");

    signUpEmail?.addEventListener("input", validateSignUp)
    signUpUsername?.addEventListener("input", validateSignUp)
    signUpPassword?.addEventListener("input", validateSignUp)
    signUpForm?.addEventListener("submit", HandleSignUp)

    function validateSignUp() {
        const valid =
            signUpEmail.value &&
            signUpUsername.value &&
            signUpPassword.value;
        signUpSubmit.disabled = !valid;
        signUpSubmit.classList.toggle("disabled", !valid);
    }
}

/*************************************************************
*     Add-post link and add-post + button(fab) listeners     *     
**************************************************************/
function AddPostListener() {
    // Floating action button (add post button)
    const fabAddPost = document.getElementById("fabAddPost");
    fabAddPost?.addEventListener("click", () => {
        document.getElementById("newPostModal").classList.remove("hidden");
    });

    // Create post in dropDown menu.
    document.getElementById("createPost")?.addEventListener("click", (e) => {
        e.preventDefault()
        document.getElementById("newPostModal")?.classList.remove("hidden");
    });
}

/**********************************
*     New Post Form listeners     *     
***********************************/
function NewPostListener() {
    const newPostModal = document.getElementById("newPostModal");
    const closeBtn = document.getElementById("closeNewPostModal");
    const newPostForm = document.getElementById("newPostForm");

    closeBtn?.addEventListener("click", () => {
        newPostModal.classList.add("hidden");
    });

    window.addEventListener("click", (e) => {
        if (e.target === newPostModal) {
            newPostModal.classList.add("hidden");
        }
    });

    HandleTags(newPostForm).then(tagHandler => {
        newPostForm.addEventListener("submit", (e) => {
            HandleNewPost(e, tagHandler.getChosenTags(), tagHandler.clearTags());
        });
    }).catch((err) => {
        console.error("HandleTags failed", err);
        RemoveError("postErrorMsg", newPostForm);
        DisplayError("postErrorMsg", newPostForm, "Something went wrong initializing tags");
    });
}

