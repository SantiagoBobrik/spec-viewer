/**
 * Mermaid diagram rendering support.
 *
 * Goldmark renders ```mermaid blocks as <pre><code class="language-mermaid">...</code></pre>.
 * This script finds those blocks, replaces them with <div class="mermaid"> containers,
 * and calls mermaid.run() to render them as SVG diagrams.
 *
 * Theme integration: reads the current dark/light mode from the <html> element's class
 * list and configures Mermaid's theme accordingly. A MutationObserver watches for theme
 * changes so diagrams are re-rendered when the user toggles dark mode.
 */
(function () {
  "use strict";

  function isDarkMode() {
    return document.documentElement.classList.contains("dark");
  }

  function getMermaidTheme() {
    return isDarkMode() ? "dark" : "default";
  }

  /**
   * Find all <pre><code class="language-mermaid"> blocks, extract their text content,
   * and replace the <pre> with a <div class="mermaid"> that Mermaid can process.
   * Returns the list of newly created mermaid divs.
   */
  function prepareMermaidBlocks() {
    var codeBlocks = document.querySelectorAll("pre > code.language-mermaid");
    var containers = [];

    codeBlocks.forEach(function (codeEl) {
      var preEl = codeEl.parentElement;
      var diagramSource = codeEl.textContent;

      var div = document.createElement("div");
      div.classList.add("mermaid");
      div.textContent = diagramSource;

      preEl.parentNode.replaceChild(div, preEl);
      containers.push(div);
    });

    return containers;
  }

  /**
   * Initialize Mermaid with the current theme and render all diagram blocks.
   */
  function initAndRender() {
    if (typeof mermaid === "undefined") {
      return;
    }

    mermaid.initialize({
      startOnLoad: false,
      theme: getMermaidTheme(),
    });

    var containers = prepareMermaidBlocks();
    if (containers.length > 0) {
      mermaid.run({ nodes: containers });
    }
  }

  /**
   * Re-render all mermaid diagrams with the updated theme. Because Mermaid
   * has already processed the divs (replacing text with SVG), we need to
   * retrieve the original source from the data attribute Mermaid stores,
   * reset each container, and re-run.
   */
  function reRenderWithTheme() {
    if (typeof mermaid === "undefined") {
      return;
    }

    mermaid.initialize({
      startOnLoad: false,
      theme: getMermaidTheme(),
    });

    var containers = document.querySelectorAll("div.mermaid");
    if (containers.length === 0) {
      return;
    }

    // Mermaid stores the original source in a data-mermaid-src attribute (v11+)
    // or we can read it from the [id].mermaidAPI internal. The simplest reliable
    // approach: Mermaid v11 keeps the original graph definition accessible via
    // data attributes. We reset each container and re-run.
    containers.forEach(function (div) {
      // Mermaid v11 stores original text in data-mermaid-original-text or
      // we stored it ourselves below as a fallback.
      var original = div.getAttribute("data-original-source");
      if (original) {
        div.removeAttribute("data-processed");
        div.innerHTML = "";
        div.textContent = original;
      }
    });

    mermaid.run({ nodes: Array.from(containers) });
  }

  // On first prepareMermaidBlocks, also stash the original source so we can
  // re-render on theme change. Override prepareMermaidBlocks to do this.
  var _originalPrepare = prepareMermaidBlocks;
  prepareMermaidBlocks = function () {
    var containers = _originalPrepare();
    containers.forEach(function (div) {
      div.setAttribute("data-original-source", div.textContent);
    });
    return containers;
  };

  // Wait for DOM ready, then initialize.
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initAndRender);
  } else {
    initAndRender();
  }

  // Watch for theme changes on <html> class list.
  var observer = new MutationObserver(function (mutations) {
    mutations.forEach(function (mutation) {
      if (
        mutation.type === "attributes" &&
        mutation.attributeName === "class"
      ) {
        reRenderWithTheme();
      }
    });
  });

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });
})();
