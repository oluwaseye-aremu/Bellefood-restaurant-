// Menu Dynamic Script with Cart Functionality
const API_URL = `${window.location.origin}/api`;

let cart = JSON.parse(localStorage.getItem('cart') || '[]');
let allMenuItems = []; // Store all items globally for search
let isFetching = false; // Anti-hammering safety latch flag

// --- 1. HELPER FUNCTIONS ---

function updateSection(sectionId, items) {
    const container = document.getElementById(`${sectionId}-grid`);
    if (!container) return;

    if (items.length === 0) {
        container.innerHTML = '<p class="text-center" style="grid-column: 1/-1; color: #666;">No items match your search.</p>';
        return;
    }

    container.innerHTML = items.map(item => `
        <div class="card food-item" data-item='${JSON.stringify(item).replace(/'/g, "&#39;")}'>
            <div class="card-img">
                <img src="${API_URL.replace('/api', '')}${item.image_url}" 
                     alt="${item.title}" 
                     onerror="this.src='https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=600'">
            </div>
            <div class="card-body">
                <h3>${item.title}</h3>
                <p class="price">₦${item.price.toLocaleString()}</p>
                <p style="font-size: 0.9rem; color: #666; margin: 0.5rem 0;">${item.description ? item.description.substring(0, 80) + '...' : ''}</p>
                <button class="btn btn-primary btn-sm mt-1" onclick="addToCart(${item.id})">
                    Add to Cart
                </button>
                <button class="btn btn-secondary btn-sm mt-1" onclick="viewDetails(${item.id})" style="margin-left: 0.5rem;">
                    View Details
                </button>
            </div>
        </div>
    `).join('');
}

function displayMenuItems(items) {
    const categories = {
        starters: [],
        main: [],
        specials: [],
        desserts: [],
        drinks: [],
        wine: [],
        kids: [],
        dietary: []
    };

    items.forEach(item => {
        if (categories[item.category]) {
            categories[item.category].push(item);
        }
    });

    updateSection('starters', categories.starters);
    updateSection('main', categories.main);
    updateSection('specials', categories.specials);
    updateSection('desserts', categories.desserts);
    updateSection('drinks', categories.drinks);
    updateSection('wine', categories.wine);
    updateSection('kids', categories.kids);
    updateSection('dietary', categories.dietary);
}

// --- 2. MAIN LOAD LOGIC ---

async function loadMenuItems(category = '') {
    if (isFetching) return; // Block overlapping requests from firing simultaneously
    isFetching = true;

    try {
        const url = category ? `${API_URL}/menu?category=${category}` : `${API_URL}/menu`;
        const response = await fetch(url);

        if (!response.ok) {
            if (response.status === 429) {
                throw new Error('Rate limit exceeded. Please wait a moment before refreshing.');
            }
            throw new Error(`Server returned status: ${response.status}`);
        }

        const items = await response.json();
        allMenuItems = items; // Save to global variable
        displayMenuItems(items);
    } catch (error) {
        console.error('Error loading menu:', error);
        showError(error.message || 'Failed to load menu items. Is the backend running?');
    } finally {
        isFetching = false; // Open the latch back up
    }
}

// --- 3. SEARCH LOGIC ---

function handleSearch(e) {
    const searchTerm = e.target.value.toLowerCase();

    const filteredItems = allMenuItems.filter(item => {
        return item.title.toLowerCase().includes(searchTerm) ||
            item.description.toLowerCase().includes(searchTerm) ||
            item.category.toLowerCase().includes(searchTerm);
    });

    displayMenuItems(filteredItems);
}

// --- 4. CART & UI HELPERS ---

function addToCart(itemId) {
    const itemElement = document.querySelector(`[data-item*='"id":${itemId}']`);
    if (!itemElement) return;

    const item = JSON.parse(itemElement.dataset.item);
    const existingItem = cart.find(i => i.id === item.id);

    if (existingItem) {
        existingItem.quantity++;
    } else {
        cart.push({
            id: item.id,
            title: item.title,
            price: item.price,
            image_url: item.image_url,
            quantity: 1
        });
    }

    localStorage.setItem('cart', JSON.stringify(cart));
    updateCartCount();
    showNotification(`${item.title} added to cart!`);
}

function viewDetails(itemId) {
    const itemElement = document.querySelector(`[data-item*='"id":${itemId}']`);
    if (!itemElement) return;

    const item = JSON.parse(itemElement.dataset.item);
    const modal = document.getElementById('foodModal');

    if (document.getElementById('modalImg')) document.getElementById('modalImg').src = `${API_URL.replace('/api', '')}${item.image_url}`;
    if (document.getElementById('modalTitle')) document.getElementById('modalTitle').textContent = item.title;
    if (document.getElementById('modalDesc')) document.getElementById('modalDesc').textContent = item.description;
    if (document.getElementById('modalPrice')) document.getElementById('modalPrice').textContent = `₦${item.price.toLocaleString()}`;
    if (document.getElementById('modalIngredients')) document.getElementById('modalIngredients').textContent = item.ingredients || 'Not specified';

    const orderBtn = document.querySelector('.modal-details .btn-primary');
    if (orderBtn) {
        orderBtn.onclick = () => {
            addToCart(item.id);
            hideModal();
        };
    }

    if (modal) {
        modal.style.display = 'flex';
        setTimeout(() => modal.classList.add('show'), 10);
        document.body.style.overflow = 'hidden';
    }
}

function hideModal() {
    const modal = document.getElementById('foodModal');
    if (modal) {
        modal.classList.remove('show');
        setTimeout(() => {
            modal.style.display = 'none';
            document.body.style.overflow = 'auto';
        }, 300);
    }
}

function updateCartCount() {
    const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);
    let badge = document.getElementById('cart-count');

    if (!badge) {
        const navLinks = document.querySelector('.nav-links');
        if (navLinks) {
            const cartItem = document.createElement('li');
            cartItem.id = 'dynamic-cart-node';
            cartItem.innerHTML = `
                <a href="cart.html" style="position: relative; display: inline-block; padding: 0.5rem;">
                    <i class="fas fa-shopping-cart"></i>
                    <span id="cart-count" style="position: absolute; top: -4px; right: -4px; background: #c0392b; color: white; border-radius: 50%; width: 18px; height: 18px; display: flex; align-items: center; justify-content: center; font-size: 0.65rem;">${totalItems}</span>
                </a>
            `;

            const dynamicAuth = document.getElementById('dynamic-auth-node');
            if (dynamicAuth) {
                navLinks.insertBefore(cartItem, dynamicAuth);
            } else {
                navLinks.appendChild(cartItem);
            }
            badge = document.getElementById('cart-count');
        }
    }

    if (badge) {
        badge.textContent = totalItems;
        badge.style.display = totalItems > 0 ? 'flex' : 'none';
    }
}

function showNotification(message) {
    // Remove any existing notification blocks to avoid filling up the layout view space
    const oldNotification = document.getElementById('dynamic-toast-popup');
    if (oldNotification) oldNotification.remove();

    const notification = document.createElement('div');
    notification.id = 'dynamic-toast-popup';
    notification.style.cssText = `
        position: fixed;
        top: 100px;
        right: 20px;
        background: #d4af37;
        color: #1a1a1a;
        padding: 1rem 1.5rem;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        z-index: 10000;
        animation: slideIn 0.3s ease;
        font-weight: 600;
    `;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 2000);
}

function showError(message) {
    // Remove duplication blocks
    const oldError = document.getElementById('dynamic-error-toast');
    if (oldError) oldError.remove();

    const error = document.createElement('div');
    error.id = 'dynamic-error-toast';
    error.style.cssText = `
        position: fixed;
        top: 100px;
        left: 50%;
        transform: translateX(-50%);
        background: #c0392b;
        color: white;
        padding: 1rem 2rem;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        z-index: 10000;
    `;
    error.textContent = message;
    document.body.appendChild(error);
    setTimeout(() => error.remove(), 4000);
}

function filterByCategory(category) {
    loadMenuItems(category);
}

function createSectionGrids() {
    const sections = ['starters', 'main', 'specials', 'desserts', 'drinks', 'wine', 'kids', 'dietary'];
    sections.forEach(section => {
        if (!document.getElementById(`${section}-grid`)) {
            const container = document.querySelector(`.${section} .grid-3`);
            if (container) {
                container.id = `${section}-grid`;
            }
        }
    });
}

function addCategoryFilters() {
    const menuHeader = document.querySelector('.page-header');
    if (menuHeader && !document.getElementById('category-filters')) {
        const filterDiv = document.createElement('div');
        filterDiv.id = 'category-filters';
        filterDiv.style.cssText = 'margin-top: 1rem; display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap;';
        filterDiv.innerHTML = `
            <button class="btn btn-sm btn-primary" onclick="filterByCategory('')">All</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('starters')">Starters</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('main')">Main</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('specials')">Specials</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('desserts')">Desserts</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('drinks')">Drinks</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('wine')">Wine</button>
            <button class="btn btn-sm btn-outline" onclick="filterByCategory('kids')">Kids</button>
        `;
        menuHeader.appendChild(filterDiv);
    }
}

// --- 5. INITIALIZATION RUN ENGINE ---

document.addEventListener('DOMContentLoaded', () => {
    // Only run the menu loader if we are actually on the specific menu page view containing the starters section block element
    if (document.getElementById('starters-grid') || window.location.pathname.includes('menu.html')) {
        createSectionGrids();
        loadMenuItems();
        addCategoryFilters();

        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', handleSearch);
        }
    }
    updateCartCount();
});

// Animations styling rules injection
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from { transform: translateX(400px); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }
    @keyframes slideOut {
        from { transform: translateX(0); opacity: 1; }
        to { transform: translateX(400px); opacity: 0; }
    }
`;
document.head.appendChild(style);