function openClose(id, parent) {
  var el = document.getElementById(id);
  if (el.classList.contains("fa-plus")) {
    el.classList.remove("fa-plus");
    el.classList.add("fa-minus");
    document.getElementById(parent).classList.add("show");
    if (parent === "permit") {
      let aux = document.getElementById("permit");
      if (aux.classList.contains("show")) {
        aux.classList.remove("show");
        aux.click();
      }
    }
  } else {
    el.classList.remove("fa-minus");
    el.classList.add("fa-plus");
    document.getElementById(parent).classList.remove("show");
  }
}
