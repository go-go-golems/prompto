// Utility functions for managing favorites
function initFavorites() {
    if (!localStorage.getItem('favorites')) {
        localStorage.setItem('favorites', JSON.stringify([]));
    }
}

function getFavorites() {
    return JSON.parse(localStorage.getItem('favorites') || '[]');
}

function copyToClipboard(text) {
    fetch("/prompts/" + text, {
        headers: {
            'Accept': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            navigator.clipboard.writeText(data.content).then(() => {
                const toastEl = document.getElementById('copyToast');
                const toastBody = toastEl.querySelector('.toast-body');
                toastBody.innerHTML = `
                    <i class="bi bi-clipboard-check me-2"></i>Copied to clipboard!<br>
                    <small class="text-white-50">
                        ${data.stats.tokens} tokens • ${data.stats.lines} lines • ${data.stats.size} bytes
                    </small>
                `;
                const toast = new bootstrap.Toast(toastEl);
                toast.show();
            });
        });
}

function addToFavorites(name) {
    const favorites = getFavorites();
    if (!favorites.includes(name)) {
        favorites.push(name);
        localStorage.setItem('favorites', JSON.stringify(favorites));
        renderFavorites();
        
        const toastEl = document.getElementById('favToast');
        const toast = new bootstrap.Toast(toastEl);
        toast.show();
    }
}

function removeFromFavorites(name) {
    const favorites = getFavorites();
    const newFavorites = favorites.filter(fav => fav !== name);
    localStorage.setItem('favorites', JSON.stringify(newFavorites));
    renderFavorites();
}

function renderFavorites() {
    const favorites = getFavorites();
    const favoritesList = document.getElementById('favorites-list');
    if (!favoritesList) return;
    
    favoritesList.innerHTML = favorites.length === 0 
        ? '<p class="text-muted mb-0">No favorites yet</p>'
        : favorites.map(fav => `
            <div class="d-flex justify-content-between align-items-center mb-2">
                <a href="/prompts/${fav}" class="text-decoration-none">${fav}</a>
                <div>
                    <button class="btn btn-sm btn-outline-secondary me-2" onclick="copyToClipboard('${fav}')">
                        <i class="bi bi-clipboard"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="removeFromFavorites('${fav}')">
                        <i class="bi bi-x-lg"></i>
                    </button>
                </div>
            </div>
        `).join('');
}

// Initialize favorites when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    initFavorites();
    renderFavorites();
});