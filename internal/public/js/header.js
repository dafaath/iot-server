const loginRegisterSection = document.querySelector("#login-register-section");
const logoutSection = document.querySelector("#logout-section");

const authorizationCookie = Cookies.get("authorization");
console.log("🚀 ~ file: util.js:5 ~ authorizationCookie:", authorizationCookie);
if (authorizationCookie) {
  loginRegisterSection.style.display = "none";
  const jwt = authorizationCookie.split(" ")[1];
  const decoded = jwt_decode(jwt);
  if (decoded.isAdmin === true) {
    document.querySelector(
      "#head-username"
    ).innerHTML = `${decoded.username} <span class="badge bg-primary">Admin</span>`;
  } else {
    document.querySelector(
      "#head-username"
    ).innerHTML = `${decoded.username} <span class="badge bg-primary">User</span>`;
  }
  document.querySelector("#head-email").innerHTML = decoded.email;
} else {
  logoutSection.style.display = "none";
}

const logoutButton = document.querySelector("#logout-button");

logoutButton?.addEventListener("click", (e) => {
  e.preventDefault();
  Cookies.remove("authorization");
  window.location.href = "/";
});
