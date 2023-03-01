console.log("hardware-form.js loaded");
switch (type) {
  case "microcontroller unit":
    document.getElementById("microcontroller unit").selected = true;
    break;
  case "single-board computer":
    document.getElementById("single-board computer").selected = true;
    break;
  case "sensor":
    document.getElementById("sensor").selected = true;
    break;
}

const isEdit = window.location.href.includes("edit");
const separated = window.location.href.split("/");
const id = separated[separated.length - 2];
let editOptions = {};
if (isEdit) {
  editOptions = {
    url: `/hardware/${id}`,
    method: "PUT",
  };
}

handleFormSubmit({
  url: "/hardware/",
  ...editOptions,
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/hardware";
    }, 1000);
  },
});
