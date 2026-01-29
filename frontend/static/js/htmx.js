const errorDiv = document.getElementById('memberError');

function handleError(event) {
    if (event.detail.xhr.status === 404) {
        errorDiv.innerHTML = "User not found"
        return
    }
}

function checkDuplicate(event) {
    const select = document.getElementById('members');
    const errorDiv = document.getElementById('memberError');

    const lastOption = select.lastElementChild;
    if (!lastOption) return;

    const duplicate = Array.from(select.options).slice(0, -1)
        .some(opt => opt.value === lastOption.value);

    if (duplicate) {
        select.removeChild(lastOption);
        errorDiv.innerText = `User "${lastOption.text}" is already added!`;
        errorDiv.style.display = "block";
        setTimeout(() => errorDiv.style.display = "none", 3000);
    } else {
        errorDiv.style.display = "none";
    }
}