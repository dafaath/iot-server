console.log("node.js testing");

function changeLimit() {
  const limitInput = document.getElementById("limit");
  const limit = limitInput.value;
  const searchParams = new URLSearchParams(window.location.search);
  searchParams.set("limit", limit);
  window.location.search = searchParams.toString();
}

const searchParams = new URLSearchParams(window.location.search);
const limit = searchParams.get("limit");
if (limit) {
  const limitInput = document.getElementById("limit");
  limitInput.value = limit;
}
