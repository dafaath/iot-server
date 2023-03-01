handleFormSubmit({
  url: "/channel/",
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/node";
    }, 1000);
  },
  alterData: (data) => {
    data.id_node = parseInt(data.id_node);
    return data;
  },
});
