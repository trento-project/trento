/* eslint-disable no-undef */
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
      const tagSearchParam = 'tags';

      function refreshFilters(action, actionedTag) {
        fetch(
          '/api/tags?' +
            new URLSearchParams({
              resource_type: resourceType,
            })
        )
          .then((res) => res.json())
          .then((tags) => {
            var href = new URL(window.location.href);
            href.searchParams.delete(tagSearchParam);

            const oldOptions = Array.from(tagsFilter.options);
            const selectedOldOptions = oldOptions
              .filter((opt) => opt.selected)
              .map((opt) => opt.value);
            oldOptions.forEach((opt) => opt.remove());

            tags.forEach((tag) => {
              tagsFilter.add(new Option(tag, tag));
              if (selectedOldOptions.includes(tag)) {
                href.searchParams.append(tagSearchParam, tag);
              }
            });

            const selectedOptions = href.searchParams.getAll(tagSearchParam);
            $(tagsFilter).selectpicker('val', selectedOptions);
            $(tagsFilter).selectpicker('refresh');

            // Only reload the table on case of removal and tag doesn't exist anymore
            if (
              action == 'remove' &&
              selectedOldOptions.includes(actionedTag)
            ) {
              reloadTable(href.pathname + href.search);
            }
            history.pushState(undefined, '', href);
          });
      }

      tagify.on('remove', (e) => {
        const tag = e.detail.data.value;
        if (!tag) {
          return;
        }

        fetch('/api/' + resourceType + '/' + resourceId + '/tags/' + tag, {
          method: 'DELETE',
        }).then(function () {
          refreshFilters('remove', tag);
        });
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
        }).then(function () {
          refreshFilters('updated', tag);
        });
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
