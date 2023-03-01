const isEdit = window.location.href.includes("edit");
const separated = window.location.href.split("/");
const id = separated[separated.length - 2];
let editOptions = {};
if (isEdit) {
  editOptions = {
    url: `/node/${id}`,
    method: "PUT",
  };
}
console.log("test");
handleFormSubmit({
  url: "/node/",
  ...editOptions,
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/node";
    }, 1000);
  },
  alterData: (data) => {
    data.id_hardware_node = parseInt(data.id_hardware_node);
    return data;
  },
});
