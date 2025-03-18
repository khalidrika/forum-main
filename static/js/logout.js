async function HandleLogout() {
    try {
        const res = await fetch("/api/logout", { method: "POST" });
        if (!res.ok) {
            const errData = await res.json()
            console.log(errData.msg);
            PopError("Something went wrong.")
        } else {
            if (window.location.href.includes("profile")) {
                window.location.href = "/";
            } else {
                window.location.reload()
            }
        }
    } catch (err) {
        console.log(err);
        PopError("Something went wrong.")
    }
}
