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

if (localStorage.getItem("user")) {
  const roomCreationContainer = document.getElementById("room-creation");
  roomCreationContainer.innerHTML = `
        <h3>Создание комнаты</h3>
        <form id="room-creation-form">
          <label for="name">Название комнаты</label>
          <input type="text" name="name" id="name" />

          <label for="description">Описание комнаты</label>
          <input type="text" name="description" id="description" />

          <button type="submit">Создать</button>
        </form>
  `;
}
