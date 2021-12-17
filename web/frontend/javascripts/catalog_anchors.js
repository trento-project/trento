$(function () {
  const anchor = window.location.hash.slice(1);
  if (anchor) {
    const anchorEl = '#collapse-' + anchor;
    $(anchorEl).collapse('show');
  }
});
