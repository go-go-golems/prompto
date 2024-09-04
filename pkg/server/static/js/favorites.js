// static/js/favorites.js
let favorites = [];

function addToFavorites(promptName) {
    if (!favorites.includes(promptName)) {
        favorites.push(promptName);
        renderFavorites();
    }
}

function removeFromFavorites(promptName) {
    favorites = favorites.filter(fav => fav !== promptName);
    renderFavorites();
}

function renderFavorites() {
    const favoritesList = document.getElementById('favorites-list');
    favoritesList.innerHTML = '';
    favorites.forEach(fav => {
        const li = document.createElement('li');
        li.innerHTML = `
            <a href="/prompts/${fav}">${fav}</a>
            <span class="clipboard-icon" onclick="copyToClipboard('/prompts/${fav}')">📋</span>
            <span class="remove-icon" onclick="removeFromFavorites('${fav}')">-</span>
        `;
        favoritesList.appendChild(li);
    });
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        alert('Copied to clipboard');
    }, function(err) {
        alert('Failed to copy: ', err);
    });
}