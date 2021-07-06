// ajax operation to reload the table on pagination operations
// query can be `page` or `per_page`
function reloadTable(path) {
  $.get(path, function(response){
    var table = $(response).find('.table-responsive');
    $('.table-responsive').replaceWith(table);
    var nav = $(response).find('.pagination-wrap');
    $('.pagination-wrap').replaceWith(nav);
  });
}

$(document).ready(function() {
  // enable bootstrap tooltips
  $('[data-toggle="tooltip"]').tooltip();

  let now = new Date();
  $("#last_update").html(now.toLocaleString());

  // pagination events
  $('body').on('click', '.page-item', function() {
    var href = new URL(window.location.href);
    href.searchParams.set('page', this.value);
    path = href.pathname + href.search;
    reloadTable(path)
    history.pushState(undefined, '', href);
  });

  $('body').on('click', '.pagination-wrap .dropdown-item', function() {
    var href = new URL(window.location.href);
    href.searchParams.set('per_page', this.textContent);
    path = href.pathname + href.search;
    reloadTable(path)
    history.pushState(undefined, '', href);
  });
});
