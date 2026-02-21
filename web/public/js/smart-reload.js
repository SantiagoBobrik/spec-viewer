// Smart reload module — WebSocket-based live content refresh with scroll preservation
(function () {
  "use strict";

  document.addEventListener("alpine:init", function () {
    Alpine.data("smartReload", function () {
      return {
        init() {
          var scrollContainer = this.$el.closest(".overflow-y-auto");
          connect(scrollContainer);
        },
      };
    });
  });

  function connect(scrollContainer) {
    var ws = new WebSocket("ws://" + window.location.host + "/ws");

    ws.onmessage = function (event) {
      if (event.data !== "reload") return;

      var file = new URLSearchParams(window.location.search).get("file");
      if (!file) {
        window.location.reload();
        return;
      }

      var scrollTop = scrollContainer ? scrollContainer.scrollTop : 0;

      fetch("/api/view?file=" + encodeURIComponent(file))
        .then(function (resp) {
          if (!resp.ok) throw new Error("Failed to fetch content");
          return resp.text();
        })
        .then(function (html) {
          var el = document.getElementById("spec-content");
          if (!el) {
            window.location.reload();
            return;
          }

          el.innerHTML = html;

          if (scrollContainer) {
            scrollContainer.scrollTop = scrollTop;
          }
          if (window.reconcileComments) window.reconcileComments();
          if (window.applyCommentMarkers) window.applyCommentMarkers();
        })
        .catch(function () {
          window.location.reload();
        });
    };

    ws.onclose = function () {
      setTimeout(function () {
        connect(scrollContainer);
      }, 1000);
    };
  }
})();
