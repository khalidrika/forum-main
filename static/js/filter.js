async function FilterCategories() {
    const tagFilterSection = document.getElementById("tagFilterSection");
    const selectedTagsDiv = document.getElementById("selectedTags");
    const tagInput = document.getElementById("tagFilterInput");
    const suggestionsBox = document.getElementById("tagSuggestions");

    const TAGS = await FetchCategories(tagFilterSection);

    // Show all suggestions on focus if input is empty
    tagInput.addEventListener("focus", () => {
        if (!tagInput.value.trim()) {
            const availableTags = TAGS.filter((tag) => !selectedTags.includes(tag));
            showSuggestions(availableTags);
        }
    });

    // Filter suggestions as user types
    tagInput.addEventListener("input", (e) => {
        const query = e.target.value.trim().toLowerCase();
        let matched = TAGS.filter((tag) =>
            tag.toLowerCase().includes(query) && !selectedTags.includes(tag)
        );

        matched = sortByQuery(matched, query);
        matched.length === 0 ? hideSuggestions() : showSuggestions(matched);
    });

    document.addEventListener("click", (e) => {
        if (!suggestionsBox.contains(e.target) && e.target !== tagInput) {
            hideSuggestions();
        }
    });

    function showSuggestions(tagArray) {
        suggestionsBox.classList.remove("hidden");
        suggestionsBox.innerHTML = "";
        tagArray.forEach((tag) => {
            const item = document.createElement("div");
            item.classList.add("tag-suggestion-item");
            item.textContent = tag;

            item.addEventListener("click", () => {
                if (!selectedTags.includes(tag)) {
                    selectedTags.push(tag);
                    renderSelectedTags();
                    tagInput.value = "";
                    hideSuggestions();

                    offset = 0;
                    // Load the first page with the new tags
                    LoadPosts(offset, selectedTags.join(",")).then(() => {
                        offset += HomeLimit;
                    });
                }
            });
            suggestionsBox.appendChild(item);
        });
    }

    function hideSuggestions() {
        suggestionsBox.classList.add("hidden");
        suggestionsBox.innerHTML = "";
    }

    function renderSelectedTags() {
        selectedTagsDiv.innerHTML = "";
        selectedTags.forEach((tag) => {
            const pill = document.createElement("div");
            pill.className = "tag-pill";
            pill.textContent = tag;

            const removeBtn = document.createElement("button");
            removeBtn.className = "remove-tag-btn";
            removeBtn.innerHTML = "&times;";
            removeBtn.addEventListener("click", () => {
                // Remove this tag from array
                selectedTags = selectedTags.filter((t) => t !== tag);
                renderSelectedTags();

                // When removing a tag, also reset offset=0
                offset = 0;
                // Reload posts from scratch with updated tags
                LoadPosts(offset, selectedTags.join(",")).then(() => {
                    offset += HomeLimit;
                });
            });

            pill.appendChild(removeBtn);
            selectedTagsDiv.appendChild(pill);
        });

        // Disable input if maximum tags are selected
        if (selectedTags.length >= 3) {
            tagInput.disabled = true;
            tagInput.placeholder = "3 tags selected";
        } else {
            tagInput.disabled = false;
            tagInput.placeholder = "Type a tag...";
        }
    }
}
// Call the function handling Filter
FilterCategories()