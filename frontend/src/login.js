import { OpaqueClientWrapper } from "./opaque.js"
import { base64ToBytes, uint8ArrayToBase64 } from "./utils.js"

const form = document.getElementById("loginForm");
const errBox = document.getElementById("errorBox");

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const password = form.password.value;
    const username = form.username.value;

    const opaque = new OpaqueClientWrapper("cool-password-manager");

    try {
        const ke1 = await opaque.loginInit(password);

        const res1 = await fetch("/account/auth/login/init/", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username,
                ke1: uint8ArrayToBase64(ke1),
            }),
        });

        const res1Data = await res1.json();
        if (!res1.ok) {
            errBox.innerHTML = res1Data.message;
            return;
        }

    } catch (err) {
        console.error(err);
        errBox.innerHTML = "Registration failed. See console for details.";
    }
});