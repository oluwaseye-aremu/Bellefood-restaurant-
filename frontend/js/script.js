document.addEventListener('DOMContentLoaded', () => {
    // Hero Slider - FIXED VERSION
    const slides = document.querySelectorAll('.slide');
    let currentSlide = 0;
    const slideInterval = 5000; // 5 seconds

    function nextSlide() {
        if (!slides.length) return; // Guard clause

        console.log('Changing slide from', currentSlide, 'to', (currentSlide + 1) % slides.length); // Debug log

        slides[currentSlide].classList.remove('active');
        currentSlide = (currentSlide + 1) % slides.length;
        slides[currentSlide].classList.add('active');
    }

    if (slides.length > 0) {
        console.log('Found', slides.length, 'slides. Starting slider...'); // Debug log
        setInterval(nextSlide, slideInterval);
    } else {
        console.error('No slides found!'); // Debug log
    }

    // Mobile Navigation
    const navToggle = document.querySelector('.nav-toggle');
    const navLinks = document.querySelector('.nav-links');

    if (navToggle) {
        navToggle.addEventListener('click', () => {
            navLinks.classList.toggle('active');
        });
    }

    // Scroll Effect for Navbar
    const navbar = document.querySelector('.navbar');
    if (navbar) {
        window.addEventListener('scroll', () => {
            if (window.scrollY > 50) {
                navbar.style.padding = '0.5rem 0';
                navbar.style.backgroundColor = 'rgba(26, 26, 26, 1)';
            } else {
                navbar.style.padding = '1rem 0';
                navbar.style.backgroundColor = 'rgba(26, 26, 26, 0.95)';
            }
        });
    }

    // --- SHARED MODAL CLOSING LOGIC ---
    // Note: The opening logic is handled by menu-dynamic.js for the menu page.
    // script.js handles closing the modal and outside clicks.

    const modal = document.getElementById('foodModal');
    const closeModal = document.querySelector('.close-modal');

    function hideModalGlobal() {
        if (modal) {
            modal.classList.remove('show');
            setTimeout(() => {
                modal.style.display = 'none';
            }, 300); // Wait for transition
            document.body.style.overflow = 'auto';
        }
    }

    if (closeModal) {
        closeModal.addEventListener('click', hideModalGlobal);
    }

    // Close modal when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target === modal) {
            hideModalGlobal();
        }
    });

    // Smooth Scroll for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth'
                });
                if (navLinks && navLinks.classList.contains('active')) {
                    navLinks.classList.remove('active');
                }
            }
        });
    });

    // --- DYNAMIC NAVBAR AUTH SYNCHRONIZATION ---
    syncNavbarAuth();
});

// Functions exposed globally to handle state renders and user events cleanly
function syncNavbarAuth() {
    const navLinks = document.querySelector('.nav-links');
    if (!navLinks) return;

    const token = localStorage.getItem('bf_token');
    const userString = localStorage.getItem('bf_user');

    // Remove any stale dynamic login elements to prevent layout duplicates
    const oldAuthNode = document.getElementById('dynamic-auth-node');
    if (oldAuthNode) oldAuthNode.remove();

    const li = document.createElement('li');
    li.id = 'dynamic-auth-node';

    if (token && userString) {
        try {
            const user = JSON.parse(userString);
            const shortName = user.name ? user.name.split(' ')[0] : 'User';

            // Logged-in styling block matching your layout palette
            li.innerHTML = `
                <div style="display: flex; align-items: center; gap: 12px; background: rgba(255,255,255,0.05); padding: 0.4rem 1rem; border-radius: 50px; border: 1px solid rgba(214,175,55,0.3); margin-left: 10px;">
                    <img src="${user.avatar_url || 'https://via.placeholder.com/32'}" alt="Profile Image" style="width: 28px; height: 28px; border-radius: 50%; object-fit: cover; border: 1px solid #d4af37;">
                    <span style="color: #ffffff; font-weight: bold; font-size: 0.9rem;">Hi, ${shortName}</span>
                    <button onclick="triggerLogout()" style="background: none; border: none; color: #ec5b5b; font-weight: bold; cursor: pointer; font-size: 0.85rem; margin-left: 5px; transition: opacity 0.2s;" onmouseover="this.style.opacity=0.7" onmouseout="this.style.opacity=1">Logout</button>
                </div>
            `;
        } catch (e) {
            console.error("Error unpacking customer session profile payload stream:", e);
            renderLoggedOutLink(li);
        }
    } else {
        renderLoggedOutLink(li);
    }

    navLinks.appendChild(li);
}

function renderLoggedOutLink(listItemElement) {
    listItemElement.innerHTML = `<a href="auth.html" class="btn btn-primary" style="padding: 0.5rem 1.2rem; display: inline-block; text-decoration: none; margin-left: 10px;">Login / Sign Up</a>`;
}

function triggerLogout() {
    localStorage.removeItem('bf_token');
    localStorage.removeItem('bf_user');
    window.location.href = 'index.html';
}