const cooldownTime = 7000; // 7 seconds
const lastReloadTime = parseInt(localStorage.getItem("lastReloadTime") || "0", 10);

function updateCountdown() {
    const now = Date.now();
    const timeLeft = cooldownTime - (now - lastReloadTime);
    const remainingSeconds = Math.max(0, Math.ceil(timeLeft / 1000));
    document.getElementById("countdown").textContent = remainingSeconds + " seconds";

    if (remainingSeconds <= 0) {
        // Reset counter and clear start time
        localStorage.setItem("reloadCount", "0");
        localStorage.removeItem("lastReloadTime");

        // Retrieve lastPage, clear it, and redirect back
        let lastPage = sessionStorage.getItem("lastPage");
        sessionStorage.removeItem("lastPage");
        window.location.href = lastPage || "/";
    } else {
        setTimeout(updateCountdown, 1000); // Recheck each second
    }
}

updateCountdown();