let popThrottle = false; 

// Append an error to the form
function DisplayError(id, formElem, text) {    
    const errorEl = document.createElement("div");
    errorEl.id = id;
    errorEl.className = "error-message";
    errorEl.textContent = text;
    formElem.appendChild(errorEl);
}

// Remove error (so they don't accumulate when using DisplayError)
function RemoveError(id, formElem) {
    const existing = document.getElementById(id);
    if (existing && formElem.contains(existing)) {
        existing.remove();
    }
}

// Pop up a fading notifation error.
function PopError(message) {
    // If a notification is already being shown, exit early
    if (popThrottle) {
        return;
    }

    popThrottle = true;

    const errorNotification = document.createElement('div');
    errorNotification.className = "error-notification";
    errorNotification.textContent = message;

    document.body.appendChild(errorNotification);
    setTimeout(() => {
        errorNotification.style.opacity = '0';
        errorNotification.style.transform = 'translateY(-20px)';
        setTimeout(() => {
            errorNotification.remove();
            popThrottle = false;
        }, 500);
    }, 3000);
}

// Reset modal's forms and error messages after closing
function ResetModal() {
    document.getElementById("signUpForm")?.reset();
    document.getElementById("loginForm")?.reset();
    const signUpError = document.getElementById("signUpErrorMsg");
    if (signUpError) signUpError.textContent = "";
    const loginError = document.getElementById("loginErrorMsg");
    if (loginError) loginError.textContent = "";
}

// Focus on field that triggered the error
function FocusOnField(errMsg) {
    const fieldMap = {
        email: "signUpEmail",
        username: "signUpUsername",
        password: "signUpPassword"
    };

    for (const [key, fieldId] of Object.entries(fieldMap)) {
        if (errMsg.includes(key)) {
            document.getElementById(fieldId)?.focus();
        }
    }
}
