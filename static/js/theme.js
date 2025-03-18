// Load theme and handle switching from dark to light mode.
function LoadTheme() {
    const themeToggleHTML = `
        <button
            class="theme-toggle"
            type="button"
            title="Dark/Light Theme"
            aria-label="Dark/Light Theme">
            <img
                src="../img/sun.svg" 
                alt="Switch to light mode" 
                class="theme-icon">
        </button>
    `;

    // Dynamically insert the toggle button into the document.
    document.body.insertAdjacentHTML('afterbegin', themeToggleHTML);

    // Get a reference to the newly inserted button and to its icon
    const themeToggleButton = document.querySelector('.theme-toggle');
    const themeIcon = themeToggleButton.querySelector('.theme-icon');

    // A helper function to apply a theme and save it in localStorage
    function setTheme(mode) {
        if (mode === 'dark') {
            document.documentElement.classList.add('dark-mode');
            localStorage.setItem('theme', 'dark');
            themeIcon.src = '../img/sun.svg'; // Replace icon
            themeIcon.alt = 'Switch to light mode';
        } else {
            document.documentElement.classList.remove('dark-mode');
            localStorage.setItem('theme', 'light');
            themeIcon.src = '../img/moon.svg'; // Replace icon
            themeIcon.alt = 'Switch to dark mode';
        }
        updateTagIcons(); // Update tags icons
    }

    // Check localStorage for a saved theme
    const savedTheme = localStorage.getItem('theme');

    // Also get user’s system preference
    const darkModeQuery = window.matchMedia('(prefers-color-scheme: dark)');

    // Decide which theme to apply on initial load
    if (savedTheme) {
        // If user has explicitly chosen a theme before, use it
        setTheme(savedTheme);
    } else {
        // Otherwise, use system preference
        if (darkModeQuery.matches) {
            setTheme('dark');
        } else {
            setTheme('light');
        }
    }

    // Toggle dark mode on button click and update localStorage
    themeToggleButton.addEventListener('click', () => {
        const isCurrentlyDark = document.documentElement.classList.contains('dark-mode');
        if (isCurrentlyDark) {
            setTheme('light');
        } else {
            setTheme('dark');
        }
    });

    // Theme to auto-update ONLY if the user never set a preference:
    darkModeQuery.addEventListener('change', (e) => {
        // Only auto-change if the user has no saved preference
        if (!localStorage.getItem('theme')) {
            if (e.matches) {
                setTheme('dark');
            } else {
                setTheme('light');
            }
        }
    });
}
