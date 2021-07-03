
function toggleMenu(){
  var x = document.getElementById("navLinksDiv");
  $('#navLinksDiv').slideToggle("slow");

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