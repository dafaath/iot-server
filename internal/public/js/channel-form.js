handleFormSubmit({
  url: "/channel/",
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/sensor";
    }, 1000);
  },
  alterData: (data) => {
    data.value = parseFloat(data.value);
    data.id_sensor = parseInt(data.id_sensor);
    return data;
  },
});
