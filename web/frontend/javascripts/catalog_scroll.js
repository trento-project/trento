const anchor = window.location.hash.slice(1);
if (anchor) {
  const anchorEl = '#' + anchor;
  $(anchorEl).closest('.card').find('.collapse-toggle').click();
  const checkElem = $(anchorEl).closest('tr').find('a');
  checkElem.click();
  setTimeout(() => {
    checkElem[0].scrollIntoView({ behavior: 'smooth', block: 'end' });
  }, 500);
}
