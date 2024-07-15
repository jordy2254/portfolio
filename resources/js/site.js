function toggleMenu() {
  const navLinksDiv = document.getElementById('navLinksDiv');
  navLinksDiv.classList.toggle("expanded")
}

function scrollFunction() {
  if (window.scrollY > 20) {
    topButton.style.display = 'block';
  } else {
    topButton.style.display = 'none';
  }
}

function topFunction(element) {
  element.scrollTop = 0;
}

function windowResizeListener() {
    const navLinksDiv = document.getElementById('navLinksDiv');
    if (window.innerWidth >= 600 && navLinksDiv.style.display !== 'none') {
      navLinksDiv.style.display = '';
    }
}

document.addEventListener('DOMContentLoaded', function () {
  window.addEventListener('resize', windowResizeListener);

  //Get the button
  const topButton = document.getElementById('topButton');

  // When the user scrolls down 20px from the top of the document, show the button
  document.getElementById('mainContent').addEventListener('scroll', function () {
    scrollFunction(document.getElementById('mainContent'));
  });
  window.addEventListener('scroll', function () {
    scrollFunction(window);
  });
  topButton.addEventListener('click', function () {
    topFunction(null);
  });
});
