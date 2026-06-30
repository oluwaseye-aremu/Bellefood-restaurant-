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
// Functions exposed globally to handle state renders and user events cleanly
function syncNavbarAuth() {
    const navLinks = document.querySelector('.nav-links');
    if (!navLinks) return;

    const token = localStorage.getItem('bf_token');
    const userString = localStorage.getItem('bf_user');

    // Remove any stale dynamic elements to prevent layout duplicates
    const oldAuthNode = document.getElementById('dynamic-auth-node');
    if (oldAuthNode) oldAuthNode.remove();

    const oldOrdersNode = document.getElementById('customer-orders-nav-node');
    if (oldOrdersNode) oldOrdersNode.remove();

    if (token && userString) {
        try {
            const user = JSON.parse(userString);
            const shortName = user.name ? user.name.split(' ')[0] : 'User';

            // 1. DYNAMIC REDIRECT LINK CHECK BASED ON USER ROLE
            let dashboardUrl = 'index.html';
            if (user.role === 'rider') {
                dashboardUrl = 'rider-dashboard.html';
            } else if (user.role === 'admin') {
                dashboardUrl = 'admin/dashboard.html';
            }

            // 2. IF CUSTOMER: Inject the "My Orders" link right into the main navbar sequence
            if (user.role !== 'rider' && user.role !== 'admin') {
                const ordersLi = document.createElement('li');
                ordersLi.id = 'customer-orders-nav-node';
                ordersLi.innerHTML = `
                    <a href="orders.html" style="display: flex; align-items: center; gap: 6px;">
                        <i class="fas fa-receipt" style="color: #d4af37; font-size: 0.9rem;"></i> My Orders
                    </a>
                `;
                navLinks.appendChild(ordersLi);
            }

            // 3. BUILD PROFILE CONTROL CAPSULE
            const li = document.createElement('li');
            li.id = 'dynamic-auth-node';
            li.innerHTML = `
                <div style="display: flex; align-items: center; gap: 12px; background: rgba(255,255,255,0.05); padding: 0.4rem 1rem; border-radius: 50px; border: 1px solid rgba(214,175,55,0.3); margin-left: 10px;">
                    <a href="${dashboardUrl}" style="display: flex; align-items: center; gap: 8px; text-decoration: none; transition: opacity 0.2s;" onmouseover="this.style.opacity=0.8" onmouseout="this.style.opacity=1">
                        <img src="https://ui-avatars.com/api/?name=${encodeURIComponent(user.name || 'User')}&background=d4af37&color=fff&bold=true&rounded=true" alt="Profile Image" style="width: 28px; height: 28px; border-radius: 50%; object-fit: cover; border: 1px solid #d4af37;">
                        <span style="color: #ffffff; font-weight: bold; font-size: 0.9rem;">Hi, ${shortName}</span>
                    </a>
                    <span style="color: rgba(255,255,255,0.2);">|</span>
                    <button onclick="triggerLogout()" style="background: none; border: none; color: #ec5b5b; font-weight: bold; cursor: pointer; font-size: 0.85rem; transition: opacity 0.2s;" onmouseover="this.style.opacity=0.7" onmouseout="this.style.opacity=1">Logout</button>
                </div>
            `;
            navLinks.appendChild(li);

        } catch (e) {
            console.error("Error unpacking customer session profile payload stream:", e);
            fallbackLoggedOutRender(navLinks);
        }
    } else {
        fallbackLoggedOutRender(navLinks);
    }
}

// Helper block to clean up the logged out rendering routing
function fallbackLoggedOutRender(navLinks) {
    const li = document.createElement('li');
    li.id = 'dynamic-auth-node';
    if (typeof renderLoggedOutLink === 'function') {
        renderLoggedOutLink(li);
    } else {
        li.innerHTML = `<a href="auth.html">Login</a>`;
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