// ===== SLIDER ROLLER (carousel jeux) =====
(function () {
    const track   = document.getElementById('carousel-track');
    const prevBtn = document.getElementById('carousel-prev');
    const nextBtn = document.getElementById('carousel-next');
    if (!track || !prevBtn || !nextBtn) return;

    let currentIndex = 0;
    const slides     = track.querySelectorAll('.card--game');
    const visible    = 7;                          // cartes visibles à la fois
    const maxIndex   = Math.max(0, slides.length - visible);

    function slideWidth() {
        if (!slides[0]) return 0;
        const style = getComputedStyle(slides[0]);
        return slides[0].offsetWidth + parseFloat(style.marginRight || 8);
    }

    function update() {
        track.style.transform = `translateX(-${currentIndex * slideWidth()}px)`;
        prevBtn.style.opacity        = currentIndex === 0        ? '0' : '1';
        prevBtn.style.pointerEvents  = currentIndex === 0        ? 'none' : 'auto';
        nextBtn.style.opacity        = currentIndex >= maxIndex  ? '0' : '1';
        nextBtn.style.pointerEvents  = currentIndex >= maxIndex  ? 'none' : 'auto';
    }

    prevBtn.addEventListener('click', () => { if (currentIndex > 0)         { currentIndex--; update(); } });
    nextBtn.addEventListener('click', () => { if (currentIndex < maxIndex)  { currentIndex++; update(); } });
    window.addEventListener('resize', update);
    update();
})();

// ===== THEME TOGGLE =====
(function () {
    const btn  = document.getElementById('theme-toggle');
    const icon = document.getElementById('theme-icon');
    if (!btn) return;

    const MOON = '<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>';
    const SUN  = '<circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>';

    function apply(theme) {
        document.body.classList.toggle('theme-light', theme === 'light');
        if (icon) icon.innerHTML = theme === 'light' ? MOON : SUN;
    }

    const saved = localStorage.getItem('gf-theme') || 'dark';
    apply(saved);

    btn.addEventListener('click', () => {
        const next = document.body.classList.contains('theme-light') ? 'dark' : 'light';
        localStorage.setItem('gf-theme', next);
        apply(next);
    });
})();
