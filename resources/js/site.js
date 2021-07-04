
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

function expandIntroduction() {
  $('.introductionExtended').show();
  $('#expandIntroductionLink').hide();
  $('#expandIntroductionEllipse').hide();
}

$(document).ready(function (){
  $(window).resize(function (){
    console.log("Resize...")
    if($(window).width() >= 600 && $('#navLinksDiv')[0].hasAttribute("style")){
      $('#navLinksDiv').removeAttr('style')
    }
  })

  function topFunction(element) {
    element.scrollTop(0);
  }
  //Get the button
  var mybutton = $('#topButton')
  let cElem = null

  // When the user scrolls down 20px from the top of the document, show the button
  $('#mainContent').on('scroll', () => scrollFunction($('#mainContent')));
  $(window).on('scroll', () => scrollFunction($(window)));
  mybutton.on('click', ()=>topFunction(cElem));

  function scrollFunction(element) {
    console.log("Scroll even triggered")
    if (element.scrollTop() > 20 || element.scrollTop() > 20) {
      mybutton.fadeIn()
      cElem = element
    } else {
      mybutton.fadeOut()
      mybutton.hover(false)
    }
  }
});