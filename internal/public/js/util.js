async function handleFormSubmit({
  url,
  method = "post",
  showSuccess = true,
  successMessage = "",
  handleResponse = null,
  alterData = null,
}) {
  console.log("handleFormSubmit called");
  const form = document.querySelector("#submit-form");
  if (!form) {
    throw new Error("No form found with id #submit-form");
  }

  form?.addEventListener("submit", (e) => {
    e.preventDefault();
    console.log("form submitted");
    //Get the entire form fields
    let form = e.currentTarget;

    //Get URL for api endpoint
    const formData = new FormData(form);
    let data = {};
    formData.forEach((value, key) => {
      // Removing array in name
      key = key.replace("[]", "");

      // Reflect.has in favor of: object.hasOwnProperty(key)
      if (!Reflect.has(data, key)) {
        data[key] = value;
        return;
      }

      if (!Array.isArray(data[key])) {
        data[key] = [data[key]];
      }

      data[key].push(value);
    });

    console.log("Before alter", data);
    if (alterData) {
      data = alterData(data);
    }
    console.log("After alter", data);

    // Show loading
    showLoading(true);

    return axios({
      method: method,
      url: url,
      data: data,
    })
      .then((res) => {
        console.log("🚀 ~ file: util.js:23 ~ .then ~ res:", res);
        if (showSuccess) {
          const swalOptions = {
            position: "top",
            icon: "success",
            title: successMessage ? successMessage : res.data,
            showConfirmButton: false,
            toast: true,
            timer: 5000,
          };
          Swal.fire(swalOptions);
        }

        if (handleResponse) {
          handleResponse(res);
        }
      })
      .catch((err) => {
        if (err.response) {
          const swalOptions = {
            position: "top",
            icon: "error",
            title: err.response.data,
            showConfirmButton: false,
            toast: true,
            timer: 5000,
          };
          Swal.fire(swalOptions);
        }
        console.log(
          "🚀 ~ file: util.js:29 ~ form?.addEventListener ~ err:",
          err
        );
      })
      .finally(() => {
        showLoading(false);
      });
  });
}

function showLoading(show) {
  const loading = document.querySelector("#loading");
  if (show) {
    loading.style.visibility = "visible";
  } else {
    loading.style.visibility = "hidden";
  }
}

function deleteItem(object, id, identifier) {
  console.log("deleteItem called");
  const swalOptions = {
    title: `Are you sure you want to delete ${object} ${identifier} with id ${id}?`,
    text: "You won't be able to revert this!",
    icon: "warning",
    showCancelButton: true,
    confirmButtonColor: "#3085d6",
    cancelButtonColor: "#d33",
    confirmButtonText: "Yes, delete it!",
  };
  Swal.fire(swalOptions).then((result) => {
    if (result.isConfirmed) {
      axios
        .delete(`/${object}/${id}`)
        .then((res) => {
          console.log("🚀 ~ file: util.js:29 ~ .then ~ res:", res);
          const swalOptions = {
            position: "top",
            icon: "success",
            title: res.data,
            showConfirmButton: false,
            toast: true,
            timer: 5000,
          };
          Swal.fire(swalOptions);

          window.location.reload();
        })
        .catch((err) => {
          if (err.response) {
            const swalOptions = {
              position: "top",
              icon: "error",
              title: err.response.data,
              showConfirmButton: false,
              toast: true,
              timer: 5000,
            };
            Swal.fire(swalOptions);
          }
          console.log(
            "🚀 ~ file: util.js:29 ~ form?.addEventListener ~ err:",
            err
          );
        });
    }
  });
}
