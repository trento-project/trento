$(() => {

  function initTags() {
    const tagsInputs = $(".tags-input")

    tagsInputs.select2({
      tags: true,
      width: '300px',
      minimumInputLength: 1,
      tokenSeparators: [','],
      matcher: (_params, _data) => {
        return null
      },
      createTag: params => {
        const regex = new RegExp(/^[0-9A-Za-z\s\-_]+$/);

        if (regex.test(params.term)) {
          return {
            id: params.term,
            text: params.term
          }
        }
        return null
      }
    }).on('select2:select', e => {
      const tag = e.params.data.id
      if (tag == null) {
        return
      }

      const url = $(e.target).attr('data-url')
      $.ajax({
        url: url,
        type: 'POST',
        data: JSON.stringify({tag: tag}),
        dataType: 'json',
        success: _data => {
          let o = new Option(tag, tag);
          $(o).html(tag);

          const tagsFilter = $('#tags_filter')
          tagsFilter.append(o)
          tagsFilter.selectpicker('refresh')
        }
      })

    }).on('select2:unselect', e => {
      const url = $(e.target).attr('data-url') + "/" + e.params.data.id
      $.ajax({
        url: url,
        type: 'DELETE',
        success: _data => {
          const tagsFilter = $('#tags_filter')
          $("#tags_filter option[value='" + e.params.data.id + "']").remove()
          tagsFilter.selectpicker('refresh')
        }
      })
    }).on('select2:open', _e => {
      $('.select2-container--open .select2-dropdown--below').css('display', 'none')
    })
  }

  initTags()

  $(window).on('table:reloaded', _e => {
    initTags()
  })
})
