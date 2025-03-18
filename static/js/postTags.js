// Handle new post tags (categories)
// Prevents accidental modifications
async function HandleTags(formElem) {

    const MAX_TAGS = 3; // Up to three tags
    let chosenTags = []; // Kept as private

    const TAGS = await FetchCategories(formElem);

    const tagListEl = document.getElementById("tagList");
    const tagInputEl = document.getElementById("tagInput");
    const suggestionsBox = document.getElementById("suggestionsBox");

    // Hide suggestions
    function hideSuggestions() {
        suggestionsBox.classList.add("hidden");
        suggestionsBox.innerHTML = "";
    }

    // Show suggestions (suggestions drop-down in UI)
    function showSuggestions(suggestions) {
        suggestionsBox.innerHTML = ""
        suggestionsBox.classList.remove("hidden")

        // Exclude already chosen tags
        const filteredSuggestions = suggestions.filter((suggest) => !chosenTags.includes(suggest));

        filteredSuggestions.forEach((suggest) => {
            const item = document.createElement("div")
            item.className = "suggestion-item"
            item.textContent = suggest
            item.addEventListener("click", () => {
                addTag(suggest)
                tagInputEl.value = "";
                hideSuggestions()
            })
            suggestionsBox.appendChild(item);
        })
    }

    // --- Render chosen tags to the UI ---
    function renderChosenTags() {
        tagListEl.innerHTML = "";
        chosenTags.forEach((tag) => {
            const pill = document.createElement("div");
            pill.className = "tag-pill";
            pill.textContent = tag;

            // remove button
            const removeBtn = document.createElement("button");
            removeBtn.className = "remove-tag";
            removeBtn.innerHTML = "&times;";
            removeBtn.addEventListener("click", () => removeTag(tag));
            pill.appendChild(removeBtn);

            tagListEl.appendChild(pill);
        });
    }

    // --- Add a new tag ---
    function addTag(tagName) {
        RemoveError("postErrorMsg", formElem);

        // If we've already reached the limit
        if (chosenTags.length >= MAX_TAGS) {
            DisplayError("postErrorMsg", formElem, "You can only select up to 3 tags.");
            return;
        }

        if (!chosenTags.includes(tagName)) {
            chosenTags.push(tagName);
            renderChosenTags();
        }
    }

    // --- Remove a tag ---
    function removeTag(tagName) {
        chosenTags = chosenTags.filter((t) => t !== tagName);
        renderChosenTags();
    }

    /*****************************************************
    *                  Event Listeners                   *
    ******************************************************/

    // Show all suggestions on focus/click if input is empty
    tagInputEl?.addEventListener("focus", () => {
        if (!tagInputEl.value.trim()) {
            showSuggestions(TAGS);
        }
    });

    // Event Listeners for the input
    tagInputEl?.addEventListener("input", (e) => {
        const query = e.target.value.trim().toLowerCase();
        if (!query) {
            hideSuggestions();
            return;
        }
        let filtered = TAGS.filter((tag) => tag.toLowerCase().includes(query));
        filtered = sortByQuery(filtered, query)

        if (filtered.length === 0) {
            hideSuggestions();
        } else {
            showSuggestions(filtered);
        }
    });

    // Hide suggestions if user clicks outside box/input ---
    document.addEventListener("click", (e) => {
        if (!suggestionsBox.contains(e.target) && e.target !== tagInputEl) {
            hideSuggestions();
        }
    });

    // ------ Expose a method to retrieve chosenTags ------
    // aka using enclosure to get the last updated chosenTags
    return {
        getChosenTags: () => chosenTags,
        clearTags: () => {
            tagInputEl.value = "" // Input
            chosenTags = []; // The list
            renderChosenTags(); // UI.
        }
    }
}

/****************************************************
*        Fetch categories from the server           *
****************************************************/
async function FetchCategories(formElem) {
    let tags = []
    try {
        const res = await fetch("/api/get-categories");
        if (!res.ok) {
            const errData = await res.json();
            DisplayError("postErrorMsg", formElem, errData.msg || "Error fetching categories.");
            return tags; // Stop execution here
        }

        tags = await res.json();
        tags.sort((a, b) => a.localeCompare(b))

    } catch (err) {
        console.log(err);
    }
    return tags
}

/******************************************************
*                   Post Categories                   *
* NO UI for this, it's only for admin/devlopper usage *
******************************************************/

// async function AddCategories() {
//     try {
//         const resp = await fetch("http:127.0.0.1:8080/api/add-categories", {
//             method: "POST",
//             headers: {
//                 "content-type": "application/json",
//                 cookie: "session_token:KinoTan_n session token"
//             },
//             body: JSON.stringify({
//                 "categories": ["cat1", "cat2", "cat3"]
//             })
//         })

//         if (!resp.ok) {
//             console.log("failed");
//             const respData = await resp.json()
//             console.log(respData.msg);         
//         } else {
//             console.log("success");            
//         }
//     } catch (err) {
//         console.log(err);
//     }
// }
// AddCategories()
