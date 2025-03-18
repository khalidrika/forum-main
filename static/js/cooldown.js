window.addEventListener("beforeunload", () => {
    const currentCount = parseInt(localStorage.getItem("reloadCount") || "0", 10);
    localStorage.setItem("reloadCount", currentCount + 1); // increment on reload
});

window.addEventListener("DOMContentLoaded", () => {
    const currentCount = parseInt(localStorage.getItem("reloadCount") || "0", 10);
    const lastCooldown = parseInt(localStorage.getItem("lastReloadTime") || "0", 10);
    const now = Date.now();
    const cooldownTime = 7000; // 7 seconds

    // console.log("Number of reloads:", currentCount);

    // Make sure cooldown is applied even on a different tab 
    // (check if we are in the middle of a cooldown)
    if (lastCooldown && now - lastCooldown < cooldownTime) {
        window.location.href = "/cooldown";
        return;
    }

    if (currentCount > 7) {
        // Reset counter and store start time
        localStorage.setItem("reloadCount", 0);
        localStorage.setItem("lastReloadTime", now.toString());

        // Redirect only if not already on cooldown page
        if (!window.location.pathname.includes("cooldown")) {
            // Store last visited page before redirecting
            sessionStorage.setItem("lastPage", window.location.href);
            window.location.href = "/cooldown";
        }
    }
});

// Reset if reload isn't within 7 seconds
setInterval(() => {
    localStorage.setItem("reloadCount", 0);
}, 7000);