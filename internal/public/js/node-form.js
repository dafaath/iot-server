//Clone the hidden element and shows it
console.log("node-form.js");

function generateID() {
  const length = 32;
  let result = "";
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

function addSensor(idHardwareSensor = undefined, fieldSensor = undefined) {
  const random_id = generateID();
  const mainForm = $("#sensor-form-main");
  // Clone the main form
  const clonned = mainForm.first().clone();

  // Set select for clonned object
  // Because clone does not copy the value of select, we need to do it manually
  // https://stackoverflow.com/questions/742810/clone-isnt-cloning-select-values
  const selects = mainForm.find("select");
  $(selects).each(function (i) {
    const select = this;
    $(clonned).find("select").eq(i).val($(select).val());
  });

  clonned.attr("id", random_id);
  clonned.appendTo("#sensor-form-holder").show();
  clonned.find("#add-sensor").remove();
  clonned.find(".remove-sensor").show();

  if (idHardwareSensor) {
    clonned.find(".sensor-hardware-id").val(idHardwareSensor).change();
  }

  if (fieldSensor) {
    clonned.find(".sensor-field").val(fieldSensor).change();
  }

  // Add delete functionality to the new button
  clonned.find(".remove-sensor").click(function () {
    $("#" + random_id).remove();
  });

  // Reset main form
  mainForm.find(".sensor-hardware-id").val("default").change();
  mainForm.find(".sensor-field").val("").change();
}

$("#add-sensor").click(addSensor);

const isEdit = window.location.href.includes("edit");
const separated = window.location.href.split("/");
const id = separated[separated.length - 2];
let editOptions = {};
if (isEdit) {
  $("#id_hardware_node").val(ID_HARDWARE_NODE);

  for (let index = 0; index < ID_HARDWARE_SENSOR.length; index++) {
    const idHardwareSensor = ID_HARDWARE_SENSOR[index];
    const fieldSensor = FIELD_SENSOR[index];
    addSensor(idHardwareSensor, fieldSensor);
  }

  editOptions = {
    url: `/node/${id}`,
    method: "PUT",
  };
}

handleFormSubmit({
  url: "/node/",
  ...editOptions,
  handleResponse: (res) => {
    setTimeout(() => {
      window.location.href = "/node";
    }, 1000);
  },
  alterData: (data) => {
    if (!Array.isArray(data.id_hardware_sensor)) {
      data.id_hardware_sensor = [data.id_hardware_sensor];
    }
    if (!Array.isArray(data.field_sensor)) {
      data.field_sensor = [data.field_sensor];
    }

    data.id_hardware_node = parseInt(data.id_hardware_node);

    data.id_hardware_sensor.forEach((idHardwareSensor, i) => {
      data.id_hardware_sensor[i] = parseInt(idHardwareSensor);
    });

    return data;
  },
});
