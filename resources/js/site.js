
function toggleMenu(){
  var x = document.getElementById("navLinks");
  if (x.style.display == "none"){
    x.style.display="block";
  }else{
    x.style.display="none";
  }
}

function displayPopup(){
  if(localStorage.getItem("wipnotification") != "confirmed"){
    $('.popupBackground').fadeIn();
  }
}

function closePopup(){
  localStorage.setItem("wipnotification", "confirmed");
  $('.popupBackground').fadeOut();
}