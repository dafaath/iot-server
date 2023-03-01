handleFormSubmit({
  url: "/user/login",
  handleResponse: (res) => {
    const token = res.data;
    Cookies.set("authorization", `Bearer ${token}`, { expires: 365 });
    window.location.href = "/hardware";
  },
  successMessage: "Login Successful",
});
