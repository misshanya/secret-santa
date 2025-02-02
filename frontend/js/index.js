function fetchUserData() {
    fetch("/api/me")
    .then(response => {
        if (response.ok) {
            return response.json();
        }
        window.location.href = "/login.html"
    })
    .then(data => {
        localStorage.setItem("user", JSON.stringify(data));
    });

    updateHeader();
}

function updateHeader() {
    const user = JSON.parse(localStorage.getItem("user"));
    const userInfoContainer = document.getElementById("user-info");

  if (user) {
    userInfoContainer.innerHTML = `
            <span>Привет, ${user.name}!</span>
        `;
  } else {
    userInfoContainer.innerHTML = `
        <button onclick="window.location.href='/login.html'">Войти</button>
        `;
    }
}

window.onload = fetchUserData;