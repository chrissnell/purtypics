package editor

const editorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Purtypics Metadata Editor</title>
    <link rel="stylesheet" href="/static/editor.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>Purtypics Metadata Editor</h1>
            <div class="actions">
                <button id="saveBtn" class="btn btn-primary">Save All Changes</button>
                <span id="saveStatus"></span>
            </div>
        </header>

        <div class="tabs">
            <button class="tab-btn active" data-tab="gallery">Gallery</button>
            <button class="tab-btn" data-tab="albums">Albums</button>
            <button class="tab-btn" data-tab="photos">Photos</button>
        </div>

        <div class="tab-content">
            <!-- Gallery Tab -->
            <div id="gallery-tab" class="tab-pane active">
                <h2>Gallery Settings</h2>
                <form id="gallery-form">
                    <div class="form-group">
                        <label for="gallery-title">Gallery Title</label>
                        <input type="text" id="gallery-title" class="form-control">
                    </div>
                    <div class="form-group">
                        <label for="gallery-description">Description</label>
                        <textarea id="gallery-description" class="form-control" rows="3"></textarea>
                    </div>
                    <div class="form-group">
                        <label for="gallery-author">Author</label>
                        <input type="text" id="gallery-author" class="form-control">
                    </div>
                    <div class="form-group">
                        <label for="gallery-copyright">Copyright</label>
                        <input type="text" id="gallery-copyright" class="form-control">
                    </div>
                </form>
            </div>

            <!-- Albums Tab -->
            <div id="albums-tab" class="tab-pane">
                <h2>Albums</h2>
                <div id="albums-list" class="albums-grid"></div>
            </div>

            <!-- Photos Tab -->
            <div id="photos-tab" class="tab-pane">
                <h2>Photos</h2>
                <div class="album-selector">
                    <label for="photo-album-select">Select Album:</label>
                    <select id="photo-album-select" class="form-control">
                        <option value="">Choose an album...</option>
                    </select>
                </div>
                <div id="photos-list" class="photos-grid"></div>
            </div>
        </div>
    </div>

    <!-- Album Edit Modal -->
    <div id="album-modal" class="modal">
        <div class="modal-content">
            <h3>Edit Album</h3>
            <form id="album-form">
                <input type="hidden" id="album-path">
                <div class="form-group">
                    <label for="album-title">Title</label>
                    <input type="text" id="album-title" class="form-control">
                </div>
                <div class="form-group">
                    <label for="album-description">Description</label>
                    <textarea id="album-description" class="form-control" rows="3"></textarea>
                </div>

                <div class="form-group">
                    <label for="album-cover">Cover Photo</label>
                    <input type="hidden" id="album-cover">
                    <div id="cover-photo-selector" class="photo-grid-selector">
                        <!-- Photos will be loaded here -->
                    </div>
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="closeAlbumModal()">Cancel</button>
                    <button type="submit" class="btn btn-primary">Save</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Photo Edit Modal -->
    <div id="photo-modal" class="modal">
        <div class="modal-content">
            <h3>Edit Photo</h3>
            <div class="photo-preview">
                <img id="photo-preview-img" src="" alt="">
            </div>
            <form id="photo-form">
                <input type="hidden" id="photo-path">
                <div class="form-group">
                    <label for="photo-title">Title</label>
                    <input type="text" id="photo-title" class="form-control">
                </div>
                <div class="form-group">
                    <label for="photo-description">Description</label>
                    <textarea id="photo-description" class="form-control" rows="3"></textarea>
                </div>
                <div class="form-group">
                    <label for="photo-hidden">
                        <input type="checkbox" id="photo-hidden">
                        Hide this photo
                    </label>
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="closePhotoModal()">Cancel</button>
                    <button type="submit" class="btn btn-primary">Save</button>
                </div>
            </form>
        </div>
    </div>

    <script src="/static/editor.js"></script>
</body>
</html>`

const editorCSS = `
@import url('https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;700&display=swap');

* {
    box-sizing: border-box;
}

body {
    font-family: 'Inconsolata', monospace;
    margin: 0;
    padding: 0;
    background: #f8f9fa;
    color: #212529;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    padding-bottom: 20px;
    border-bottom: 1px dotted #17a2b8;
}

h1 {
    margin: 0;
    color: #17a2b8;
    font-weight: 700;
    text-transform: uppercase;
}

.actions {
    display: flex;
    align-items: center;
    gap: 15px;
}

#saveStatus {
    color: #28a745;
    font-size: 14px;
    font-weight: 700;
}

.btn {
    padding: 10px 20px;
    border: 1px dotted transparent;
    border-radius: 0;
    cursor: pointer;
    font-size: 14px;
    font-family: 'Inconsolata', monospace;
    font-weight: 700;
    text-transform: uppercase;
    transition: all 0.2s;
}

.btn-primary {
    background: #215e21;
    color: #FFFFFF;
    border-color: #215e21;
}

.btn-primary:hover {
    background: #1a4a1a;
    border-color: #1a4a1a;
}

.btn-secondary {
    background: #6c757d;
    color: #FFFFFF;
    border-color: #6c757d;
}

.btn-secondary:hover {
    background: #5a6268;
    border-color: #5a6268;
}

.tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 30px;
}

.tab-btn {
    padding: 10px 20px;
    background: #FFFFFF;
    border: 1px dotted #dee2e6;
    border-radius: 0;
    cursor: pointer;
    transition: all 0.2s;
    font-family: 'Inconsolata', monospace;
    font-weight: 700;
    text-transform: uppercase;
    color: #6c757d;
}

.tab-btn:hover {
    background: #f8f9fa;
    border-color: #17a2b8;
    color: #17a2b8;
}

.tab-btn.active {
    background: #ffaa00;
    color: #FFFFFF;
    border-color: #ffaa00;
}

.tab-pane {
    display: none;
    background: white;
    padding: 30px;
    border-radius: 0;
    border: 1px dotted #dee2e6;
}

.tab-pane.active {
    display: block;
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 5px;
    font-weight: 500;
}

.form-control {
    width: 100%;
    padding: 8px 12px;
    border: 1px dotted #ced4da;
    border-radius: 0;
    font-size: 14px;
    font-family: 'Inconsolata', monospace;
    background: #FFFFFF;
}

.form-control:focus {
    outline: none;
    border-color: #17a2b8;
    background: #f0f8ff;
}

textarea.form-control {
    resize: vertical;
}

.albums-grid, .photos-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.album-card, .photo-card {
    background: #FFFFFF;
    border: 1px dotted #dee2e6;
    border-radius: 0;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.2s;
}

.album-card:hover, .photo-card:hover {
    transform: translateY(-2px);
    border-color: #17a2b8;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.album-card img, .photo-card img {
    width: 100%;
    height: 200px;
    object-fit: cover;
}

.album-card-info, .photo-card-info {
    padding: 15px;
}

.album-card h3, .photo-card h3 {
    margin: 0 0 5px 0;
    font-size: 16px;
}

.album-card p, .photo-card p {
    margin: 0;
    color: #666;
    font-size: 14px;
}

.album-card p:first-of-type {
    margin-bottom: 8px;
}

.modal {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.5);
    z-index: 1000;
}

.modal-content {
    position: relative;
    background: white;
    max-width: 600px;
    margin: 50px auto;
    padding: 30px;
    border-radius: 0;
    border: 2px solid #17a2b8;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
}

.photo-preview {
    margin-bottom: 20px;
    text-align: center;
}

.photo-preview img {
    max-width: 100%;
    max-height: 300px;
    border-radius: 4px;
}

.album-selector {
    margin-bottom: 20px;
}

.album-selector label {
    margin-right: 10px;
}

.album-selector select {
    width: auto;
    min-width: 300px;
}

.hidden-badge {
    display: inline-block;
    background: #dc3545;
    color: white;
    padding: 2px 8px;
    border-radius: 0;
    font-size: 12px;
    margin-left: 10px;
    font-weight: 700;
}

.photo-grid-selector {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 10px;
    max-height: 300px;
    overflow-y: auto;
    padding: 10px;
    border: 1px dotted #ced4da;
    border-radius: 0;
    background: #f8f9fa;
}

.photo-grid-selector .photo-option {
    position: relative;
    cursor: pointer;
    border: 1px dotted transparent;
    border-radius: 0;
    overflow: hidden;
    transition: all 0.2s;
}

.photo-grid-selector .photo-option:hover {
    transform: scale(1.05);
    border-color: #17a2b8;
}

.photo-grid-selector .photo-option.selected {
    border-color: #ffaa00;
    box-shadow: 0 0 0 2px rgba(255, 170, 0, 0.3);
}

.photo-grid-selector .photo-option img {
    width: 100%;
    height: 80px;
    object-fit: contain;
    background: #000000;
}

.photo-grid-selector .photo-option .photo-name {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    font-size: 10px;
    padding: 2px 4px;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
`

const editorJS = `
let metadata = {
    title: '',
    description: '',
    author: '',
    copyright: '',
    albums: {},
    photos: {}
};

let albums = [];
let currentAlbum = null;
let autoSaveTimer = null;
let hasUnsavedChanges = false;

// Load initial data
async function loadData() {
    try {
        const [metaResponse, albumsResponse] = await Promise.all([
            fetch('/api/metadata'),
            fetch('/api/albums')
        ]);
        
        metadata = await metaResponse.json();
        albums = await albumsResponse.json();
        
        updateGalleryForm();
        renderAlbums();
        populateAlbumSelect();
    } catch (error) {
        console.error('Error loading data:', error);
    }
}

// Update gallery form with metadata
function updateGalleryForm() {
    document.getElementById('gallery-title').value = metadata.title || '';
    document.getElementById('gallery-description').value = metadata.description || '';
    document.getElementById('gallery-author').value = metadata.author || '';
    document.getElementById('gallery-copyright').value = metadata.copyright || '';
}

// Render albums grid
function renderAlbums() {
    const container = document.getElementById('albums-list');
    container.innerHTML = '';
    
    albums.forEach(album => {
        const card = document.createElement('div');
        card.className = 'album-card';
        card.onclick = () => editAlbum(album);
        
        const albumName = album.path.split('/').pop();
        const coverImage = album.coverPhoto || 'placeholder.jpg';
        
        card.innerHTML = ` + "`" + `
            <img src="/images/${albumName}/${coverImage}" alt="${album.title}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22250%22 height=%22200%22 viewBox=%220 0 250 200%22><rect fill=%22%23ddd%22 width=%22250%22 height=%22200%22/><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%22.3em%22>No Image</text></svg>'">
            <div class="album-card-info">
                <h3>${album.title}${album.hidden ? '<span class="hidden-badge">Hidden</span>' : ''}</h3>
                <p>${album.photoCount} photos</p>
                ${album.description ? '<p>' + album.description + '</p>' : ''}
            </div>
        ` + "`" + `;
        
        container.appendChild(card);
    });
}

// Populate album select for photos tab
function populateAlbumSelect() {
    const select = document.getElementById('photo-album-select');
    select.innerHTML = '<option value="">Choose an album...</option>';
    
    albums.forEach(album => {
        const option = document.createElement('option');
        option.value = album.path;
        option.textContent = album.title;
        select.appendChild(option);
    });
}

// Load photos for selected album
async function loadPhotos(albumPath) {
    if (!albumPath) {
        document.getElementById('photos-list').innerHTML = '';
        return;
    }
    
    try {
        const albumName = albumPath.split('/').pop();
        const response = await fetch('/api/photos/' + albumName);
        const photos = await response.json();
        
        renderPhotos(photos, albumPath);
    } catch (error) {
        console.error('Error loading photos:', error);
    }
}

// Render photos grid
function renderPhotos(photos, albumPath) {
    const container = document.getElementById('photos-list');
    container.innerHTML = '';
    
    const albumName = albumPath.split('/').pop();
    
    photos.forEach(photo => {
        const card = document.createElement('div');
        card.className = 'photo-card';
        card.onclick = () => editPhoto(photo, albumPath);
        
        const imageUrl = ` + "`" + `/images/${albumName}/${photo.filename}` + "`" + `;
        
        card.innerHTML = ` + "`" + `
            <img src="${imageUrl}" alt="${photo.title}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22250%22 height=%22200%22 viewBox=%220 0 250 200%22><rect fill=%22%23ddd%22 width=%22250%22 height=%22200%22/><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%22.3em%22>No Image</text></svg>'">
            <div class="photo-card-info">
                <h3>${photo.title}${photo.hidden ? '<span class="hidden-badge">Hidden</span>' : ''}</h3>
                ${photo.description ? '<p>' + photo.description + '</p>' : ''}
                ${photo.isVideo ? '<p>Video</p>' : ''}
            </div>
        ` + "`" + `;
        
        container.appendChild(card);
    });
}

// Edit album
async function editAlbum(album) {
    document.getElementById('album-path').value = album.path;
    document.getElementById('album-title').value = album.title;
    document.getElementById('album-description').value = album.description || '';

    document.getElementById('album-cover').value = album.coverPhoto || '';
    
    // Load photos for cover photo selection
    await loadCoverPhotoOptions(album);
    
    document.getElementById('album-modal').style.display = 'block';
}

// Load cover photo options for album
async function loadCoverPhotoOptions(album) {
    const selector = document.getElementById('cover-photo-selector');
    selector.innerHTML = '<div style="text-align: center; padding: 20px;">Loading photos...</div>';
    
    try {
        const albumName = album.path.split('/').pop();
        const response = await fetch('/api/photos/' + albumName);
        const photos = await response.json();
        
        selector.innerHTML = '';
        
        photos.forEach(photo => {
            const option = document.createElement('div');
            option.className = 'photo-option';
            if (photo.filename === album.coverPhoto) {
                option.classList.add('selected');
            }
            
            // Always use dynamic thumbnail for cover photo selector
            const thumbUrl = ` + "`" + `/thumbs/small/${albumName}/${photo.filename}` + "`" + `;
            
            option.innerHTML = ` + "`" + `
                <img src="${thumbUrl}" alt="${photo.filename}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22100%22 height=%2280%22 viewBox=%220 0 100 80%22><rect fill=%22%23ddd%22 width=%22100%22 height=%2280%22/><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%22.3em%22 font-size=%2210%22>No Image</text></svg>'">
                <div class="photo-name">${photo.filename}</div>
            ` + "`" + `;
            
            option.addEventListener('click', () => selectCoverPhoto(photo.filename, albumName));
            
            selector.appendChild(option);
        });
        
        if (photos.length === 0) {
            selector.innerHTML = '<div style="text-align: center; padding: 20px; color: #666;">No photos in this album</div>';
        }
    } catch (error) {
        console.error('Error loading cover photo options:', error);
        selector.innerHTML = '<div style="text-align: center; padding: 20px; color: #e74c3c;">Error loading photos</div>';
    }
}

// Select cover photo
function selectCoverPhoto(filename, albumName) {
    // Update hidden input
    document.getElementById('album-cover').value = filename;
    
    // Update visual selection
    document.querySelectorAll('#cover-photo-selector .photo-option').forEach(option => {
        option.classList.remove('selected');
    });
    
    event.currentTarget.classList.add('selected');
}

// Edit photo
function editPhoto(photo, albumPath) {
    const albumName = albumPath.split('/').pop();
    document.getElementById('photo-path').value = photo.path;
    document.getElementById('photo-title').value = photo.title;
    document.getElementById('photo-description').value = photo.description || '';
    document.getElementById('photo-hidden').checked = photo.hidden || false;
    
    // Set preview image
    const previewImg = document.getElementById('photo-preview-img');
    previewImg.src = ` + "`" + `/images/${albumName}/${photo.filename}` + "`" + `;
    
    document.getElementById('photo-modal').style.display = 'block';
}

// Close modals
function closeAlbumModal() {
    document.getElementById('album-modal').style.display = 'none';
}

function closePhotoModal() {
    document.getElementById('photo-modal').style.display = 'none';
}

// Save all changes
async function saveAll() {
    const saveBtn = document.getElementById('saveBtn');
    const saveStatus = document.getElementById('saveStatus');
    
    saveBtn.disabled = true;
    saveStatus.textContent = 'Saving...';
    
    try {
        const response = await fetch('/api/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(metadata)
        });
        
        if (response.ok) {
            hasUnsavedChanges = false;
            saveStatus.textContent = 'Saved successfully!';
            saveStatus.style.color = '#27ae60';
            
            // Restore Save button to green when saved
            const saveBtn = document.getElementById('saveBtn');
            saveBtn.style.background = '#215e21';
            saveBtn.style.borderColor = '#215e21';
            
            setTimeout(() => {
                saveStatus.textContent = '';
            }, 3000);
        } else {
            throw new Error('Save failed');
        }
    } catch (error) {
        console.error('Error saving:', error);
        saveStatus.textContent = 'Error saving changes';
        saveStatus.style.color = '#e74c3c';
        
        // Keep button red on error
        const saveBtn = document.getElementById('saveBtn');
        saveBtn.style.background = '#AA0000';
        saveBtn.style.borderColor = '#AA0000';
    } finally {
        saveBtn.disabled = false;
    }
}

// Schedule auto-save
function scheduleAutoSave() {
    hasUnsavedChanges = true;
    
    // Clear existing timer
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
    }
    
    // Update status to show unsaved changes
    const saveStatus = document.getElementById('saveStatus');
    saveStatus.textContent = 'Unsaved changes';
    saveStatus.style.color = '#e67e22';
    
    // Change Save button to red when there are unsaved changes
    const saveBtn = document.getElementById('saveBtn');
    saveBtn.style.background = '#AA0000';
    saveBtn.style.borderColor = '#AA0000';
    
    // Schedule save after 2 seconds of inactivity
    autoSaveTimer = setTimeout(() => {
        saveAll();
    }, 2000);
}

// Warn before leaving with unsaved changes
window.addEventListener('beforeunload', (e) => {
    if (hasUnsavedChanges) {
        e.preventDefault();
        e.returnValue = '';
    }
});

// Initialize event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadData();
    
    // Tab switching
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const tabName = btn.dataset.tab;
            
            // Update buttons
            document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            
            // Update panes
            document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
            document.getElementById(tabName + '-tab').classList.add('active');
        });
    });
    
    // Gallery form updates
    document.getElementById('gallery-form').addEventListener('input', (e) => {
        const field = e.target.id.replace('gallery-', '');
        metadata[field] = e.target.value;
        scheduleAutoSave();
    });
    
    // Album form submit
    document.getElementById('album-form').addEventListener('submit', (e) => {
        e.preventDefault();
        
        const path = document.getElementById('album-path').value;
        
        if (!metadata.albums) metadata.albums = {};
        
        metadata.albums[path] = {
            title: document.getElementById('album-title').value,
            description: document.getElementById('album-description').value,

            cover_photo: document.getElementById('album-cover').value
        };
        
        // Update local albums data
        const album = albums.find(a => a.path === path);
        if (album) {
            album.title = metadata.albums[path].title;
            album.description = metadata.albums[path].description;

            album.coverPhoto = metadata.albums[path].cover_photo;
        }
        
        closeAlbumModal();
        renderAlbums();
        scheduleAutoSave();
    });
    
    // Photo form submit
    document.getElementById('photo-form').addEventListener('submit', (e) => {
        e.preventDefault();
        
        const path = document.getElementById('photo-path').value;
        
        if (!metadata.photos) metadata.photos = {};
        
        metadata.photos[path] = {
            title: document.getElementById('photo-title').value,
            description: document.getElementById('photo-description').value,
            hidden: document.getElementById('photo-hidden').checked
        };
        
        closePhotoModal();
        
        // Refresh photos if we're viewing the same album
        const currentAlbumPath = document.getElementById('photo-album-select').value;
        if (currentAlbumPath && path.startsWith(currentAlbumPath)) {
            loadPhotos(currentAlbumPath);
        }
        
        scheduleAutoSave();
    });
    
    // Photo album select
    document.getElementById('photo-album-select').addEventListener('change', (e) => {
        loadPhotos(e.target.value);
    });
    
    // Save button
    document.getElementById('saveBtn').addEventListener('click', saveAll);
    
    // Modal close on background click
    document.getElementById('album-modal').addEventListener('click', (e) => {
        if (e.target.id === 'album-modal') closeAlbumModal();
    });
    
    document.getElementById('photo-modal').addEventListener('click', (e) => {
        if (e.target.id === 'photo-modal') closePhotoModal();
    });
});
`