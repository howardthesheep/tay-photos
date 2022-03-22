import ApiManager from "./apiManager.js";

const api = new ApiManager();

async function login() {
    const userStr = document.querySelector("#username").value;
    const passwordStr = document.querySelector("#password").value;

    const result = await api.login({
        username: userStr,
        password: passwordStr
    });

    // TODO: Save the returned apiToken and other stuff in session storage
    console.log(result);
}

window.onload = () => {
    const htmlButton = document.querySelector("#login-form button");
    htmlButton.addEventListener("click", login)
}