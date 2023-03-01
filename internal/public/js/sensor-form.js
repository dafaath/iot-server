const isEdit = window.location.href.includes("edit");
const separated = window.location.href.split("/");
const id = separated[separated.length - 2];
let editOptions = {};
if (isEdit) {
  editOptions = {
    url: `/sensor/${id}`,
    method: "PUT",
  };
}
handleFormSubmit({
  url: "/sensor/",
  ...editOptions,
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/sensor";
    }, 1000);
  },
  alterData: (data) => {
    data.id_node = parseInt(data.id_node);
    data.id_hardware = parseInt(data.id_hardware);
    return data;
  },
});
