const form = document.querySelector("form");
const nameInput = document.querySelector("#name");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");

form.addEventListener("submit", function (event) {
  event.preventDefault();

  const formData = {
    name: nameInput.value,
    username: usernameInput.value,
    password: passwordInput.value,
  };

  console.log(formData);

  fetch("/api/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  })
    .then((response) => {
      if (response.status === 201) {
        alert("Вы успешно зарегистрировались!");
        window.location.href = "/login.html";
      } else {
        alert("Ошибка регистрации");
        throw new Error("Ошибка регистрации");
      }
      return response.json();
    })
});
