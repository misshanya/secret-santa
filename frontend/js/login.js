const form = document.querySelector("form");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");

form.addEventListener("submit", function (event) {
  event.preventDefault();

  const formData = {
    username: usernameInput.value,
    password: passwordInput.value,
  };

  console.log(formData);

  fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((response) => {
      if (response.status === 200) {
        alert("Успешный вход!");
        window.location.href = "/index.html";
      } else {
        alert("Ошибка входа");
        throw new Error("Ошибка входа");
      }
      return response.json();
    })
});
