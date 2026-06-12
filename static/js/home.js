/**
 * home.js — Forum Gaming Homepage
 * Handles: burger menu, lazy images, sticky navbar shadow
 */

(function () {
    'use strict';

    // ---- Burger menu toggle ----
    const burger = document.getElementById('navbar-burger');
    const nav    = document.getElementById('main-nav');

    if (burger && nav) {
        burger.addEventListener('click', function () {
            const expanded = this.getAttribute('aria-expanded') === 'true';
            this.setAttribute('aria-expanded', String(!expanded));
            nav.classList.toggle('navbar__nav--open');
        });
    }

    // ---- Search drop toggle (icône loupe dans la navbar) ----
    const searchToggle = document.getElementById('navbar-search-toggle');
    const searchDrop   = document.getElementById('jv-search-drop');
    const searchInput  = document.getElementById('navbar-search-input');

    if (searchToggle && searchDrop) {
        searchToggle.addEventListener('click', function () {
            const isOpen = searchDrop.classList.toggle('is-open');
            searchDrop.setAttribute('aria-hidden', String(!isOpen));
            if (isOpen && searchInput) {
                // Petit délai pour laisser l'animation se déclencher
                setTimeout(function () { searchInput.focus(); }, 50);
            }
        });

        // Fermer en cliquant en dehors
        document.addEventListener('click', function (e) {
            if (!searchDrop.contains(e.target) && !searchToggle.contains(e.target)) {
                searchDrop.classList.remove('is-open');
                searchDrop.setAttribute('aria-hidden', 'true');
            }
        });

        // Fermer avec Échap
        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape' && searchDrop.classList.contains('is-open')) {
                searchDrop.classList.remove('is-open');
                searchDrop.setAttribute('aria-hidden', 'true');
                searchToggle.focus();
            }
        });
    }

    // ---- Sticky navbar shadow on scroll ----
    const navbar = document.getElementById('main-navbar');
    if (navbar) {
        window.addEventListener('scroll', function () {
            if (window.scrollY > 8) {
                navbar.style.boxShadow = '0 4px 24px rgba(0,0,0,0.55)';
            } else {
                navbar.style.boxShadow = '0 2px 12px rgba(0,0,0,0.4)';
            }
        }, { passive: true });
    }

    // ---- Lazy-load sidebar thumbnail images ----
    // Images are set as CSS background-image via inline style, so they load
    // automatically. This script adds a fade-in effect when they become visible.
    function observeSideThumbs() {
        if (!('IntersectionObserver' in window)) return;

        const thumbs = document.querySelectorAll('.side-news-item__thumb');
        const io = new IntersectionObserver(function (entries) {
            entries.forEach(function (entry) {
                if (entry.isIntersecting) {
                    entry.target.classList.add('side-news-item__thumb--loaded');
                    io.unobserve(entry.target);
                }
            });
        }, { rootMargin: '100px' });

        thumbs.forEach(function (thumb) { io.observe(thumb); });
    }

    observeSideThumbs();

    // ---- Hero search — prevent empty submit ----
    const heroForm = document.getElementById('hero-search-form');
    if (heroForm) {
        heroForm.addEventListener('submit', function (e) {
            const input = document.getElementById('hero-search-input');
            if (input && !input.value.trim()) {
                e.preventDefault();
                input.focus();
            }
        });
    }

    // ---- Navbar search — same guard ----
    const navForm = document.getElementById('navbar-search-form');
    if (navForm) {
        navForm.addEventListener('submit', function (e) {
            const input = document.getElementById('navbar-search-input');
            if (input && !input.value.trim()) {
                e.preventDefault();
                input.focus();
            }
        });
    }

    // ---- Subtle hover animations on category cards ----
    const cards = document.querySelectorAll('.category-card');
    cards.forEach(function (card) {
        card.addEventListener('mouseenter', function () {
            this.style.transition = 'background 0.15s';
        });
    });

    // ---- Carrousel TOP FORUMS ----
    (function () {
        const track    = document.getElementById('carousel-track');
        const viewport = document.getElementById('carousel-viewport');
        const btnPrev  = document.getElementById('carousel-prev');
        const btnNext  = document.getElementById('carousel-next');

        if (!track || !viewport || !btnPrev || !btnNext) return;

        let currentIndex = 0;

        function getItemsVisible() {
            // On lit la largeur calculée d'un item
            const item = track.querySelector('.jv-carousel__item');
            if (!item) return 4;
            return Math.round(viewport.offsetWidth / item.offsetWidth);
        }

        function getTotal() {
            return track.querySelectorAll('.jv-carousel__item').length;
        }

        function slideTo(index) {
            const total   = getTotal();
            const visible = getItemsVisible();
            const maxIdx  = Math.max(0, total - visible);
            currentIndex  = Math.max(0, Math.min(index, maxIdx));

            const item = track.querySelector('.jv-carousel__item');
            if (!item) return;
            const itemW = item.offsetWidth;
            track.style.transform = 'translateX(-' + (currentIndex * itemW) + 'px)';

            btnPrev.disabled = currentIndex === 0;
            btnNext.disabled = currentIndex >= maxIdx;
            btnPrev.style.opacity = btnPrev.disabled ? '0.35' : '1';
            btnNext.style.opacity = btnNext.disabled ? '0.35' : '1';
        }

        btnPrev.addEventListener('click', function () { slideTo(currentIndex - 1); });
        btnNext.addEventListener('click', function () { slideTo(currentIndex + 1); });

        // Recalcul au resize
        window.addEventListener('resize', function () { slideTo(currentIndex); }, { passive: true });

        // Init
        slideTo(0);
    }());

})();
