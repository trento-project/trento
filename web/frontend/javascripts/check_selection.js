$(document).ready(function() {
  $('body').on('click', '.parent', function() {
    $('.parent-' + $(this)[0].id).prop('checked', $(this)[0].checked);
  });

  $('body').on('click', 'input:checkbox[class^="parent-"]', function() {
    var temp = $(this)[0].id.split('-');
    var parentId = temp[0];

    if ($(this)[0].checked) {
      $('#' + parentId).prop('checked', true);
      return
    }

    var atLeastOneEnabled = false;
    $('input:checkbox[class="' + $(this)[0].className + '"]').each(function(index, item) {
      if (item.checked) {
        atLeastOneEnabled = true;
      }
    });
    if (!atLeastOneEnabled) {
      $('#' + parentId).prop('checked', false);
    }
  });
});
