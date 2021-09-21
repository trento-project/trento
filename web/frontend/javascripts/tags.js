$(() => {
  function initTags() {
    const inputs = document.querySelectorAll('.tags-input');

    inputs.forEach((elm) => {
      const tagify = new Tagify(elm, {
        whitelist: [],
        editTags: false,
        pattern: /^[0-9A-Za-z\s\-_]+$/,
        dropdown: {
          maxItems: 20,
          enabled: 1,
          closeOnSelect: false,
          placeAbove: false,
          classname: 'tags-look',
        },
      });

      // Add tags by clicking on the input
      tagify.DOM.scope.addEventListener('click', (_e) => tagify.addEmptyTag());

      const resourceType = elm.getAttribute('data-resource-type');
      const resourceId = elm.getAttribute('data-resource-id');
      const tagsFilter = document.getElementById('tags_filter');

      function refreshFilters() {
        fetch(
          '/api/tags?' +
            new URLSearchParams({
              resource_type: resourceType,
            })
        )
          .then((res) => res.json())
          .then((tags) => {
            const oldOptions = Array.from(tagsFilter.options);
            oldOptions.forEach((o) => o.remove());
            tags.forEach((tag) => tagsFilter.add(new Option(tag, tag)));

            $(tagsFilter).selectpicker('refresh');
          });
      }

      tagify.on('remove', (e) => {
        const tag = e.detail.data.value;
        if (!tag) {
          return;
        }

        fetch(
          '/api/' +
            resourceType +
            '/' +
            resourceId +
            '/tags/' +
            e.detail.data.value,
          {
            method: 'DELETE',
          }
        ).then(refreshFilters);
      });

      tagify.on('edit:updated', (e) => {
        const tag = e.detail.data.value;
        if (!tag) {
          return;
        }

        const object = { tag: tag };
        fetch('/api/' + resourceType + '/' + resourceId + '/tags', {
          method: 'POST',
          body: JSON.stringify(object),
        }).then(refreshFilters);
      });

      // Get tags suggestions
      tagify.on('edit:start', (_e) => {
        tagify.whitelist = null;

        fetch('/api/tags')
          .then((res) => res.json())
          .then(function (newWhitelist) {
            tagify.whitelist = newWhitelist;
          });
      });
    });
  }

  initTags();

  // Re-init Tagify if the table is reloaded
  window.addEventListener('table:reloaded', initTags);
});
