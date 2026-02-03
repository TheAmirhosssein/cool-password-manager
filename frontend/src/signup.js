import { OpaqueClientWrapper } from "./opaque.js"
import { base64ToBytes, uint8ArrayToBase64 } from "./utils.js"

const form = document.getElementById("signupForm");
const errBox = document.getElementById("errorBox");

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const password = form.password.value;
    const username = form.username.value;
    const email = form.email.value;
    const firstName = form.firstName.value;
    const lastName = form.lastName.value;

    const opaque = new OpaqueClientWrapper("cool-password-manager");

    try {
        const registrationRequest = await opaque.registerInit(password);

        const res1 = await fetch("/account/auth/sign-up/init/", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username,
                email,
                firstName,
                lastName,
                registrationRequest: uint8ArrayToBase64(registrationRequest),
            }),
        });

        const res1Data = await res1.json();
        if (!res1.ok) {
            errBox.innerHTML = res1Data.message;
            return;
        }

        const record = await opaque.registerFinish(base64ToBytes(res1Data.record), res1Data.registrationID);

        htmx.ajax("POST", "/account/auth/sign-up/final/", {
            target: "#signup-container",
            swap: "outerHTML",
            values: {
                registrationID: res1Data.registrationID,
                registrationRecord: uint8ArrayToBase64(record),
            },
        });

    } catch (err) {
        console.error(err);
        errBox.innerHTML = "Registration failed. See console for details.";
    }
});