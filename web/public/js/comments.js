// Comments module — client-side annotation system for spec blocks
(function () {
  "use strict";

  // --- localStorage helpers ---

  const COMMENT_SVG =
    '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">' +
    '<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>' +
    "</svg>";

  function currentFilePath() {
    return new URLSearchParams(window.location.search).get("file") || "";
  }

  function storageKey() {
    return "specComments:" + currentFilePath();
  }

  function parseStoredComments(key) {
    try {
      const raw = localStorage.getItem(key);
      if (!raw) return [];
      const data = JSON.parse(raw);
      if (data && data.version === 1 && Array.isArray(data.comments)) {
        return data.comments;
      }
    } catch (_) {
      // Corrupted data; fall through to empty array
    }
    return [];
  }

  function loadComments() {
    return parseStoredComments(storageKey());
  }

  function saveComments(comments) {
    if (comments.length === 0) {
      localStorage.removeItem(storageKey());
    } else {
      localStorage.setItem(
        storageKey(),
        JSON.stringify({ version: 1, comments })
      );
    }
    window.dispatchEvent(new CustomEvent("comments-changed"));
  }

  function purgeComments() {
    localStorage.removeItem(storageKey());
    window.dispatchEvent(new CustomEvent("comments-changed"));
  }

  // --- Block helpers ---

  function getBlocks() {
    const el = document.getElementById("spec-content");
    if (!el) return [];
    return Array.from(el.children);
  }

  function blockPreview(block) {
    return (block.textContent || "").trim().substring(0, 80);
  }

  // --- Reconciliation ---

  function reconcileComments() {
    const comments = loadComments();
    if (comments.length === 0) return;

    const blocks = getBlocks();
    const previews = blocks.map(blockPreview);

    for (const c of comments) {
      // Try exact index first
      if (
        blocks[c.blockIndex] &&
        blockPreview(blocks[c.blockIndex]) === c.blockTextPreview
      ) {
        continue;
      }

      // Search for best match
      let bestIdx = -1;
      let bestScore = 0;
      previews.forEach((p, i) => {
        if (p === c.blockTextPreview) {
          bestIdx = i;
          bestScore = 999;
        } else if (
          bestScore < 1 &&
          p &&
          c.blockTextPreview &&
          p.substring(0, 40) === c.blockTextPreview.substring(0, 40)
        ) {
          bestIdx = i;
          bestScore = 1;
        }
      });

      if (bestIdx >= 0) {
        c.blockIndex = bestIdx;
        c.blockTextPreview = previews[bestIdx];
      }
    }

    saveComments(comments);
  }

  // --- Comment markers ---

  function applyCommentMarkers() {
    const blocks = getBlocks();
    const comments = loadComments();

    const countMap = {};
    for (const c of comments) {
      countMap[c.blockIndex] = (countMap[c.blockIndex] || 0) + 1;
    }

    blocks.forEach((block, idx) => {
      if (getComputedStyle(block).position === "static") {
        block.style.position = "relative";
      }

      const old = block.querySelector(".comment-indicator");
      if (old) old.remove();

      block.classList.toggle("has-comment", !!countMap[idx]);

      const count = countMap[idx];
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className = "comment-indicator" + (count ? " has-comments" : "");
      btn.setAttribute("aria-label", "Add comment");
      btn.setAttribute("data-block-index", idx);
      btn.innerHTML = count
        ? `${COMMENT_SVG}<span class="comment-count">${count}</span>`
        : COMMENT_SVG;

      block.appendChild(btn);
    });
  }

  // --- LLM prompt formatting ---

  function formatCommentsForLLM() {
    const comments = loadComments();
    if (comments.length === 0) return "";

    const filePath = currentFilePath();

    const grouped = {};
    for (const c of comments) {
      if (!grouped[c.blockIndex]) {
        grouped[c.blockIndex] = { preview: c.blockTextPreview, items: [] };
      }
      grouped[c.blockIndex].items.push(c.text);
    }

    const lines = [
      "You are reviewing a technical specification document.",
      "Apply each requested change to the specification file: " + filePath,
      "Preserve the existing markdown formatting and style.",
      "Only modify the sections mentioned. Do not change anything else.",
      "",
      "---",
    ];

    const sortedKeys = Object.keys(grouped).sort((a, b) => Number(a) - Number(b));
    for (const idx of sortedKeys) {
      const g = grouped[idx];
      lines.push("");
      lines.push(`### Section starting with: "${g.preview}"`);
      for (const text of g.items) {
        lines.push("- " + text);
      }
    }

    return lines.join("\n");
  }

  // --- Alpine.js components ---

  document.addEventListener("alpine:init", () => {
    Alpine.data("commentPopover", function () {
      return {
        open: false,
        blockIndex: -1,
        blockPreviewText: "",
        comments: [],
        newComment: "",

        show(blockIdx) {
          const blocks = getBlocks();
          this.blockIndex = blockIdx;
          this.blockPreviewText = blockPreview(blocks[blockIdx] || {});
          this.comments = loadComments().filter((c) => c.blockIndex === blockIdx);
          this.newComment = "";
          this.open = true;

          this.$nextTick(() => {
            const block = blocks[blockIdx];
            const popover = document.getElementById("comment-popover");
            if (!block || !popover) return;

            const rect = block.getBoundingClientRect();
            const popoverHeight = popover.offsetHeight || 300;
            let top = rect.top;

            if (top + popoverHeight > window.innerHeight) {
              top = window.innerHeight - popoverHeight - 8;
            }
            if (top < 8) top = 8;

            popover.style.top = top + "px";
            popover.style.left = Math.max(8, rect.left - 16) + "px";
          });
        },

        addComment() {
          const text = this.newComment.trim();
          if (!text) return;

          const all = loadComments();
          all.push({
            id: Date.now().toString(36) + Math.random().toString(36).substring(2, 7),
            blockIndex: this.blockIndex,
            blockTextPreview: this.blockPreviewText,
            text,
            createdAt: new Date().toISOString(),
          });
          saveComments(all);

          this.newComment = "";
          this.comments = all.filter((c) => c.blockIndex === this.blockIndex);
          applyCommentMarkers();
        },

        deleteComment(id) {
          const all = loadComments().filter((c) => c.id !== id);
          saveComments(all);

          this.comments = all.filter((c) => c.blockIndex === this.blockIndex);
          applyCommentMarkers();

          if (this.comments.length === 0 && !this.newComment) {
            this.open = false;
          }
        },
      };
    });

    Alpine.data("copyComments", function () {
      return {
        hasComments: false,
        copied: false,

        init() {
          this.checkComments();
          window.addEventListener("comments-changed", () => {
            this.checkComments();
          });
        },

        checkComments() {
          this.hasComments = loadComments().length > 0;
        },

        copyAndPurge() {
          const prompt = formatCommentsForLLM();
          if (!prompt) return;

          navigator.clipboard.writeText(prompt).then(() => {
            this.copied = true;
            setTimeout(() => {
              purgeComments();
              applyCommentMarkers();
              this.copied = false;
            }, 1500);
          });
        },
      };
    });
  });

  // --- Sidebar badges ---

  function updateSidebarBadges() {
    const links = document.querySelectorAll("[data-spec-name]");
    for (const link of links) {
      const old = link.querySelector(".sidebar-comment-badge");
      if (old) old.remove();

      const href = link.getAttribute("href") || "";
      const match = href.match(/[?&]file=([^&]+)/);
      if (!match) continue;

      const filePath = decodeURIComponent(match[1]);
      const count = parseStoredComments("specComments:" + filePath).length;
      if (count > 0) {
        const badge = document.createElement("span");
        badge.className = "sidebar-comment-badge";
        badge.textContent = count;
        link.appendChild(badge);
      }
    }
  }

  // --- Event delegation for indicator clicks ---

  document.addEventListener("click", (e) => {
    const btn = e.target.closest(".comment-indicator");
    if (!btn) return;

    const idx = parseInt(btn.getAttribute("data-block-index"), 10);
    if (isNaN(idx)) return;

    window.dispatchEvent(
      new CustomEvent("open-comment-popover", { detail: { blockIndex: idx } })
    );
  });

  // --- Init on page load ---

  document.addEventListener("DOMContentLoaded", () => {
    updateSidebarBadges();
    if (!currentFilePath()) return;
    reconcileComments();
    applyCommentMarkers();
  });

  window.addEventListener("comments-changed", updateSidebarBadges);

  // --- Expose globally for smart reload ---

  window.applyCommentMarkers = applyCommentMarkers;
  window.reconcileComments = reconcileComments;
})();
