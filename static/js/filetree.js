var toggler = document.getElementsByClassName("caret");
var i;

for (i = 0; i < toggler.length; i++) {
  toggler[i].addEventListener("click", function() {
    this.parentElement.querySelector(".nested").classList.toggle("active");
    this.classList.toggle("caret-down");
  });
} 

function w3_open_sidebar() {
  document.getElementById("filetree").style.width = "100%";
  document.getElementById("filetree").style.display = "block";
}

function w3_close_sidebar() {
  document.getElementById("filetree").style.display = "none";
}
