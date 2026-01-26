(function () {
    try {
        var stored = localStorage.getItem("themeMode");
        var isDark = stored
            ? stored === "dark"
            : window.matchMedia("(prefers-color-scheme: dark)").matches;
        if (isDark) {
            document.documentElement.classList.add("dark");
        } else {
            document.documentElement.classList.remove("dark");
        }
    } catch (e) { }
})();
