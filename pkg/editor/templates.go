package editor

const editorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PurtyPics Gallery Editor</title>
    <link rel="stylesheet" href="/static/editor.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>PurtyPics Gallery Editor</h1>
            <div class="actions">
                <button id="saveBtn" class="btn btn-primary">Save All Changes</button>
                <button id="generateBtn" class="btn btn-secondary">Generate Gallery</button>
                <button id="viewBtn" class="btn btn-secondary">View Gallery (Local)</button>
            </div>
        </header>

        <div class="tabs">
            <button class="tab-btn active" data-tab="gallery">Gallery</button>
            <button class="tab-btn" data-tab="albums">Albums</button>
            <button class="tab-btn" data-tab="photos">Photos</button>
            <button class="tab-btn" data-tab="deploy">Deploy</button>
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
                    <div class="form-group">
                        <label>
                            <input type="checkbox" id="gallery-show-locations">
                            Show Photo Locations
                        </label>
                        <p style="margin-top: 5px; font-size: 12px; color: var(--text-secondary);">Display a map with photo locations at the bottom of gallery pages</p>
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

            <!-- Deploy Tab -->
            <div id="deploy-tab" class="tab-pane">
                <h2>Deployment Settings</h2>
                <div class="deploy-tabs">
                    <button class="deploy-tab-btn active" data-deploy-tab="rsync">Rsync over SSH</button>
                    <button class="deploy-tab-btn" data-deploy-tab="s3">AWS S3</button>
                    <button class="deploy-tab-btn" data-deploy-tab="cloudflare">Cloudflare Pages</button>
                </div>
                
                <!-- Rsync Configuration -->
                <div id="rsync-deploy-pane" class="deploy-pane active">
                    <div class="deploy-section">
                        <h3>Rsync Configuration</h3>
                        <form id="rsync-form">
                            <div class="form-group">
                                <label for="rsync-host">Host (user@hostname)</label>
                                <input type="text" id="rsync-host" class="form-control" placeholder="user@example.com">
                            </div>
                            <div class="form-group">
                                <label for="rsync-path">Remote Path</label>
                                <input type="text" id="rsync-path" class="form-control" placeholder="/var/www/html/gallery">
                            </div>
                            <div class="form-group">
                                <label for="rsync-port">SSH Port</label>
                                <input type="number" id="rsync-port" class="form-control" value="22">
                            </div>
                            <div class="deploy-actions">
                                <button type="button" class="btn btn-primary deploy-save-btn">Save Configuration</button>
                                <button type="button" class="btn btn-secondary deploy-test-btn">Test Connection (Dry Run)</button>
                                <button type="button" class="btn btn-deploy deploy-deploy-btn"><div class="progress-wipe"></div><span>Deploy Now</span></button>
                            </div>
                        </form>
                    </div>
                </div>
                
                <!-- S3 Configuration -->
                <div id="s3-deploy-pane" class="deploy-pane">
                    <div class="deploy-section">
                        <h3>AWS S3 Configuration</h3>
                        <form id="s3-form">
                            <div class="form-group">
                                <label for="s3-bucket">Bucket Name</label>
                                <input type="text" id="s3-bucket" class="form-control" placeholder="my-gallery-bucket">
                            </div>
                            <div class="form-group">
                                <label for="s3-region">Region</label>
                                <input type="text" id="s3-region" class="form-control" placeholder="us-east-1">
                            </div>
                            <div class="form-group">
                                <p style="margin: 5px 0; font-size: 14px; color: var(--text-secondary);">Set AWS credentials via environment variables:<br>
                                <code>AWS_ACCESS_KEY_ID</code> and <code>AWS_SECRET_ACCESS_KEY</code></p>
                            </div>
                            <div class="deploy-actions">
                                <button type="button" class="btn btn-primary deploy-save-btn">Save Configuration</button>
                                <button type="button" class="btn btn-secondary deploy-test-btn">Test Connection</button>
                                <button type="button" class="btn btn-deploy deploy-deploy-btn"><div class="progress-wipe"></div><span>Deploy Now</span></button>
                            </div>
                        </form>
                    </div>
                </div>
                
                <!-- Cloudflare Pages Configuration -->
                <div id="cloudflare-deploy-pane" class="deploy-pane">
                    <div class="deploy-section">
                        <h3>Cloudflare Pages Configuration</h3>
                        <form id="cloudflare-form">
                            <div class="form-group">
                                <label for="cf-project">Project Name</label>
                                <input type="text" id="cf-project" class="form-control" placeholder="my-gallery">
                            </div>
                            <div class="form-group">
                                <label for="cf-account">Account ID</label>
                                <input type="text" id="cf-account" class="form-control" placeholder="023e105f4ecef8ad9ca31a8372d0c353">
                            </div>
                            <div class="form-group">
                                <p style="margin: 5px 0; font-size: 14px; color: var(--text-secondary);">Set Cloudflare API token via environment variable:<br>
                                <code>CLOUDFLARE_API_TOKEN</code></p>
                            </div>
                            <div class="deploy-actions">
                                <button type="button" class="btn btn-primary deploy-save-btn">Save Configuration</button>
                                <button type="button" class="btn btn-secondary deploy-test-btn">Test Connection</button>
                                <button type="button" class="btn btn-deploy deploy-deploy-btn"><div class="progress-wipe"></div><span>Deploy Now</span></button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Album Edit Modal -->
    <div id="album-modal" class="modal">
        <div class="modal-content">
            <h3>Edit Album</h3>
            <form id="album-form">
                <input type="hidden" id="album-path">
                <input type="hidden" id="album-relative-path">
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

    <!-- Error Overlay -->
    <div id="error-overlay" class="error-overlay">
        <div class="error-message">
            <h2>Output Directory Not Found</h2>
            <p>Please generate the gallery first before attempting to deploy. Click the "Generate Gallery" button to create the output directory and build your gallery.</p>
            <button onclick="hideErrorOverlay()">OK</button>
        </div>
    </div>

    <!-- Friendly Deploy Notice -->
    <div id="deploy-notice" class="deploy-notice-overlay">
        <div class="deploy-notice-content">
            <h2>Hold On There, Friend!</h2>
            <p>Looks like you haven't generated your gallery yet. You'll need to build it first before we can deploy it to the world!</p>
            <p class="deploy-notice-hint">Just click the <strong>"Generate Gallery"</strong> button in the top-right corner, wait for it to finish, and then you can deploy.</p>
            <div class="deploy-notice-actions">
                <button class="btn btn-primary" onclick="hideDeployNotice(); focusGenerateButton()">Got it, let's generate first!</button>
                <button class="btn btn-secondary" onclick="hideDeployNotice()">I'll do it later</button>
            </div>
        </div>
    </div>

    <script src="/static/editor.js"></script>
</body>
</html>`

const editorCSS = `
@import url('https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;700&display=swap');

:root {
    /* Primary Colors */
    --primary-orange: #ffaa00;
    --primary-orange-light: #ffcc44;
    --primary-orange-dark: #cc8800;
    --primary-orange-alpha: rgba(255, 170, 0, 0.3);
    
    /* Accent Colors */
    --accent-teal: #00a8cc;
    --accent-teal-light: #00c9f0;
    --accent-teal-dark: #0087a3;
    
    /* Success/Error/Warning */
    --success-green: #00cc88;
    --success-green-dark: #00a66e;
    --error-red: #cc3333;
    --error-red-dark: #aa1111;
    --warning-yellow: #ffcc00;
    
    /* Neutral Colors */
    --neutral-100: #ffffff;
    --neutral-200: #f8f9fa;
    --neutral-300: #e9ecef;
    --neutral-400: #dee2e6;
    --neutral-500: #adb5bd;
    --neutral-600: #6c757d;
    --neutral-700: #495057;
    --neutral-800: #343a40;
    --neutral-900: #212529;
    --neutral-1000: #000000;
    
    /* Semantic Colors */
    --bg-primary: var(--neutral-200);
    --bg-secondary: var(--neutral-100);
    --bg-overlay: rgba(0, 0, 0, 0.5);
    --bg-overlay-dark: rgba(0, 0, 0, 0.8);
    
    --text-primary: var(--neutral-900);
    --text-secondary: var(--neutral-600);
    --text-light: var(--neutral-100);
    
    --border-light: var(--neutral-400);
    --border-medium: var(--neutral-500);
    --border-focus: var(--accent-teal);
    
    /* Component Specific */
    --shadow-light: rgba(0, 0, 0, 0.1);
    --shadow-medium: rgba(0, 0, 0, 0.2);
    --shadow-glow-orange: rgba(255, 170, 0, 0.5);
    --shadow-glow-red: rgba(255, 0, 0, 0.5);
}

* {
    box-sizing: border-box;
}

body {
    font-family: 'Inconsolata', monospace;
    margin: 0;
    padding: 0;
    background: var(--bg-primary);
    color: var(--text-primary);
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
    border-bottom: 1px dotted var(--accent-teal);
}

h1 {
    margin: 0;
    color: var(--accent-teal);
    font-weight: 700;
    text-transform: uppercase;
}

.actions {
    display: flex;
    align-items: center;
    gap: 15px;
}

#saveStatus {
    color: var(--accent-teal);
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
    min-width: 160px;
}

.btn-primary {
    background: var(--accent-teal);
    color: var(--text-light);
    border-color: var(--accent-teal);
}

.btn-primary:hover {
    background: var(--accent-teal-dark);
    border-color: var(--accent-teal-dark);
}

.btn-secondary {
    background: var(--neutral-600);
    color: var(--text-light);
    border-color: var(--neutral-600);
}

.btn-secondary:hover {
    background: var(--neutral-700);
    border-color: var(--neutral-700);
}
.btn-success {
    background: var(--success-green);
    color: var(--text-light);
    border-color: var(--success-green);
}
.btn-success:hover {
    background: var(--success-green-dark);
    border-color: var(--success-green-dark);
}
.btn-deploy {
    background: var(--error-red);
    color: var(--text-light);
    border-color: var(--error-red);
    position: relative;
    overflow: hidden;
    transition: all 0.3s;
}
.btn-deploy:hover {
    background: var(--error-red-dark);
    border-color: var(--error-red-dark);
}
.btn-deploy.deploying {
    background: var(--warning-yellow);
    border-color: var(--warning-yellow);
    color: var(--neutral-900);
}
.btn-deploy.deployed {
    background: var(--success-green);
    border-color: var(--success-green);
    color: var(--text-light);
}
.btn-deploy .progress-wipe {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: var(--warning-yellow);
    transition: left 0.5s ease-in-out;
    z-index: 1;
}
.btn-deploy.deploying .progress-wipe {
    left: 0;
}
.btn-deploy span {
    position: relative;
    z-index: 2;
}
.btn-deploy.error {
    background: var(--error-red);
    border-color: var(--error-red);
    color: var(--text-light);
}
.deploy-section {
    max-width: 600px;
}
.deploy-actions {
    margin-top: 30px;
    display: flex;
    gap: 10px;
    align-items: center;
}
.deploy-tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    border-bottom: 2px solid var(--border-light);
}
.deploy-tab-btn {
    padding: 8px 16px;
    background: none;
    border: none;
    border-bottom: 3px solid transparent;
    cursor: pointer;
    font-family: 'Inconsolata', monospace;
    font-weight: 700;
    text-transform: uppercase;
    color: var(--text-secondary);
    transition: all 0.2s;
}
.deploy-tab-btn:hover {
    color: var(--accent-teal);
}
.deploy-tab-btn.active {
    color: var(--primary-orange);
    border-bottom-color: var(--primary-orange);
}
.deploy-pane {
    display: none;
}
.deploy-pane.active {
    display: block;
}

.tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 30px;
}

.tab-btn {
    padding: 10px 20px;
    background: var(--bg-secondary);
    border: 1px dotted var(--border-light);
    border-radius: 0;
    cursor: pointer;
    transition: all 0.2s;
    font-family: 'Inconsolata', monospace;
    font-weight: 700;
    text-transform: uppercase;
    color: var(--text-secondary);
}

.tab-btn:hover {
    background: var(--bg-primary);
    border-color: var(--accent-teal);
    color: var(--accent-teal);
}

.tab-btn.active {
    background: var(--primary-orange);
    color: var(--text-light);
    border-color: var(--primary-orange);
}

.tab-pane {
    display: none;
    background: var(--bg-secondary);
    padding: 30px;
    border-radius: 0;
    border: 1px dotted var(--border-light);
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
    border: 1px dotted var(--border-light);
    border-radius: 0;
    font-size: 14px;
    font-family: 'Inconsolata', monospace;
    background: var(--bg-secondary);
}

.form-control:focus {
    outline: none;
    border-color: var(--border-focus);
    background: var(--neutral-100);
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
    background: var(--bg-secondary);
    border: 1px dotted var(--border-light);
    border-radius: 0;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.2s;
}

.album-card:hover, .photo-card:hover {
    transform: translateY(-2px);
    border-color: var(--accent-teal);
    box-shadow: 0 2px 4px var(--shadow-light);
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
    color: var(--text-secondary);
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
    background: var(--bg-overlay);
    z-index: 1000;
}

.error-overlay {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: var(--bg-overlay-dark);
    z-index: 2000;
    justify-content: center;
    align-items: center;
}

.error-message {
    background: var(--error-red-dark);
    color: var(--text-light);
    padding: 40px 60px;
    border-radius: 0;
    border: 3px solid var(--error-red);
    max-width: 600px;
    text-align: center;
    font-family: 'Inconsolata', monospace;
    box-shadow: 0 0 30px var(--shadow-glow-red);
}

.error-message h2 {
    margin: 0 0 20px 0;
    font-size: 28px;
    font-weight: 700;
    text-transform: uppercase;
}

.error-message p {
    margin: 0 0 30px 0;
    font-size: 18px;
    line-height: 1.5;
}

.error-message button {
    padding: 12px 30px;
    background: var(--neutral-100);
    color: var(--error-red-dark);
    border: none;
    font-size: 16px;
    font-weight: 700;
    text-transform: uppercase;
    cursor: pointer;
    transition: all 0.2s;
}

.error-message button:hover {
    background: var(--primary-orange);
    color: var(--text-light);
    transform: scale(1.05);
}

/* Friendly Deploy Notice */
.deploy-notice-overlay {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.85);
    z-index: 3000;
    justify-content: center;
    align-items: center;
    backdrop-filter: blur(5px);
}

.deploy-notice-content {
    background: var(--bg-secondary);
    padding: 40px;
    max-width: 550px;
    text-align: center;
    font-family: 'Inconsolata', monospace;
    border: 3px solid var(--primary-orange);
    box-shadow: 0 0 50px var(--shadow-glow-orange),
                0 10px 40px rgba(0, 0, 0, 0.3);
    animation: deployNoticeAppear 0.3s ease-out;
}

@keyframes deployNoticeAppear {
    from {
        transform: scale(0.9) translateY(-20px);
        opacity: 0;
    }
    to {
        transform: scale(1) translateY(0);
        opacity: 1;
    }
}


.deploy-notice-content h2 {
    margin: 0 0 20px 0;
    font-size: 28px;
    font-weight: 700;
    color: var(--primary-orange);
    text-transform: uppercase;
}

.deploy-notice-content p {
    margin: 0 0 20px 0;
    font-size: 16px;
    line-height: 1.6;
    color: var(--text-primary);
}

.deploy-notice-hint {
    background: var(--neutral-200);
    padding: 15px;
    border-left: 4px solid var(--primary-orange);
    font-size: 14px !important;
}

.deploy-notice-hint strong {
    color: var(--primary-orange);
}

.deploy-notice-actions {
    display: flex;
    gap: 15px;
    justify-content: center;
    margin-top: 30px;
}

.deploy-notice-actions .btn {
    min-width: 180px;
    padding: 12px 24px;
    font-size: 14px;
}

.modal-content {
    position: relative;
    background: var(--bg-secondary);
    max-width: 600px;
    margin: 50px auto;
    padding: 30px;
    border-radius: 0;
    border: 2px solid var(--accent-teal);
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 6px var(--shadow-light);
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
    background: var(--error-red);
    color: var(--text-light);
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
    border: 1px dotted var(--border-light);
    border-radius: 0;
    background: var(--bg-primary);
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
    border-color: var(--accent-teal);
}

.photo-grid-selector .photo-option.selected {
    border-color: var(--primary-orange);
    box-shadow: 0 0 0 2px var(--primary-orange-alpha);
}

.photo-grid-selector .photo-option img {
    width: 100%;
    height: 80px;
    object-fit: contain;
    background: var(--neutral-1000);
}

.photo-grid-selector .photo-option .photo-name {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    background: rgba(0, 0, 0, 0.7);
    color: var(--text-light);
    font-size: 10px;
    padding: 2px 4px;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

#saveBtn, #generateBtn {
    min-width: 180px;
}

#generateBtn {
    position: relative;
    overflow: hidden;
}

#generateBtn.generating {
    background: var(--accent-teal);
    color: var(--text-light);
    border-color: var(--accent-teal);
}

#generateBtn .progress-bar {
    position: absolute;
    top: 0;
    left: 0;
    width: 0%;
    height: 100%;
    background: var(--accent-teal-light);
    z-index: 0;
}

#generateBtn .btn-text {
    position: relative;
    z-index: 1;
}
`

const editorJS = `
let metadata = {
    title: '',
    description: '',
    author: '',
    copyright: '',
    show_locations: false,
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
    document.getElementById('gallery-show-locations').checked = metadata.show_locations || false;
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
        // Use first photo as default cover if none set
        const coverImage = album.coverPhoto || (album.photos && album.photos.length > 0 ? album.photos[0] : 'placeholder.jpg');
        
        card.innerHTML = ` + "`" + `
            <img src="/images/${albumName}/${coverImage}" alt="${album.title}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22250%22 height=%22200%22 viewBox=%220 0 250 200%22><rect fill=%22%23ddd%22 width=%22250%22 height=%22200%22/><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%22.3em%22>No Image</text></svg>'">
            <div class="album-card-info">
                <h3>${album.title}${album.hidden ? '<span class="hidden-badge">Hidden</span>' : ''}</h3>
                <p>${album.photoCount} ${album.photoCount === 1 ? 'photo' : 'photos'}</p>
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
    
    let firstAlbumPath = null;
    
    albums.forEach((album, index) => {
        const option = document.createElement('option');
        option.value = album.path;
        option.textContent = album.title;
        select.appendChild(option);
        
        // Store the first album path
        if (index === 0) {
            firstAlbumPath = album.path;
        }
    });
    
    // Automatically select and load the first album
    if (firstAlbumPath && albums.length > 0) {
        select.value = firstAlbumPath;
        loadPhotos(firstAlbumPath);
    }
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
    document.getElementById('album-relative-path').value = album.relativePath;
    document.getElementById('album-title').value = album.title;
    document.getElementById('album-description').value = album.description || '';

    // Set cover photo - use first photo if none selected
    let coverPhoto = album.coverPhoto;
    if (!coverPhoto && album.photos && album.photos.length > 0) {
        coverPhoto = album.photos[0];
    }
    document.getElementById('album-cover').value = coverPhoto || '';
    
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
        
        // Use first photo as default if no cover selected
        let selectedCover = album.coverPhoto;
        if (!selectedCover && photos.length > 0) {
            selectedCover = photos[0].filename;
            // Update the hidden input with the default
            document.getElementById('album-cover').value = selectedCover;
        }
        
        photos.forEach(photo => {
            const option = document.createElement('div');
            option.className = 'photo-option';
            if (photo.filename === selectedCover) {
                option.classList.add('selected');
            }
            
            // Always use dynamic thumbnail for cover photo selector
            const thumbUrl = ` + "`" + `/thumbs/small/${albumName}/${photo.filename}` + "`" + `;
            
            option.innerHTML = ` + "`" + `
                <img src="${thumbUrl}" alt="${photo.filename}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22100%22 height=%2280%22 viewBox=%220 0 100 80%22><rect fill=%22%23ddd%22 width=%22100%22 height=%2280%22/><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%22.3em%22 font-size=%228%22>Click to select</text><text fill=%22%23999%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22 dy=%221.5em%22 font-size=%228%22>cover photo</text></svg>'">
                <div class="photo-name">${photo.filename}</div>
            ` + "`" + `;
            
            option.addEventListener('click', () => selectCoverPhoto(photo.filename, albumName));
            
            selector.appendChild(option);
        });
        
        if (photos.length === 0) {
            selector.innerHTML = '<div style="text-align: center; padding: 20px; color: var(--text-secondary);">No photos in this album</div>';
        }
    } catch (error) {
        console.error('Error loading cover photo options:', error);
        selector.innerHTML = '<div style="text-align: center; padding: 20px; color: var(--error-red);">Error loading photos</div>';
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

// Show error overlay
function showErrorOverlay() {
    const overlay = document.getElementById('error-overlay');
    overlay.style.display = 'flex';
}

// Hide error overlay
function hideErrorOverlay() {
    const overlay = document.getElementById('error-overlay');
    overlay.style.display = 'none';
}

// Show deploy notice
function showDeployNotice() {
    const overlay = document.getElementById('deploy-notice');
    overlay.style.display = 'flex';
}

// Hide deploy notice
function hideDeployNotice() {
    const overlay = document.getElementById('deploy-notice');
    overlay.style.display = 'none';
}

// Focus and highlight generate button
function focusGenerateButton() {
    const generateBtn = document.getElementById('generateBtn');
    if (generateBtn) {
        // Add a glow effect to draw attention
        generateBtn.style.boxShadow = '0 0 20px var(--primary-orange), 0 0 40px var(--primary-orange-alpha)';
        generateBtn.style.transform = 'scale(1.05)';
        
        // Scroll to the button if needed
        generateBtn.scrollIntoView({ behavior: 'smooth', block: 'center' });
        
        // Remove the glow after a few seconds
        setTimeout(() => {
            generateBtn.style.boxShadow = '';
            generateBtn.style.transform = '';
        }, 3000);
    }
}

// Save all changes
async function saveAll() {
    const saveBtn = document.getElementById('saveBtn');
    const originalText = saveBtn.textContent;
    
    saveBtn.disabled = true;
    saveBtn.textContent = 'Saving...';
    saveBtn.style.transform = 'scale(0.95)';
    
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
            saveBtn.textContent = 'Saved!';
            
            // Restore Save button to teal when saved
            saveBtn.style.background = 'var(--accent-teal)';
            saveBtn.style.borderColor = 'var(--accent-teal)';
            
            setTimeout(() => {
                saveBtn.textContent = originalText;
                saveBtn.style.transform = '';
            }, 2000);
        } else {
            throw new Error('Save failed');
        }
    } catch (error) {
        console.error('Error saving:', error);
        saveBtn.textContent = 'Error!';
        
        // Keep button red on error
        saveBtn.style.background = 'var(--error-red)';
        saveBtn.style.borderColor = 'var(--error-red)';
        
        setTimeout(() => {
            saveBtn.textContent = originalText;
            saveBtn.style.background = 'var(--error-red-dark)';
            saveBtn.style.borderColor = 'var(--error-red-dark)';
            saveBtn.style.transform = '';
        }, 2000);
    } finally {
        saveBtn.disabled = false;
        saveBtn.style.transform = '';
    }
}

// Schedule auto-save
function scheduleAutoSave() {
    hasUnsavedChanges = true;
    
    // Clear existing timer
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer);
    }
    
    // Change Save button to red when there are unsaved changes
    const saveBtn = document.getElementById('saveBtn');
    saveBtn.style.background = 'var(--error-red-dark)';
    saveBtn.style.borderColor = 'var(--error-red-dark)';
    
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
    
    // Load deployment config if deploy tab is initially active
    if (document.querySelector('.tab-btn.active[data-tab="deploy"]')) {
        loadDeployConfig();
    }
    
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
            
            // Load deployment config when switching to deploy tab
            if (tabName === 'deploy') {
                loadDeployConfig();
            }
        });
    });
    
    // Gallery form updates
    document.getElementById('gallery-form').addEventListener('input', (e) => {
        const field = e.target.id.replace('gallery-', '').replace(/-/g, '_');
        if (e.target.type === 'checkbox') {
            metadata[field] = e.target.checked;
        } else {
            metadata[field] = e.target.value;
        }
        scheduleAutoSave();
    });
    
    // Album form submit
    document.getElementById('album-form').addEventListener('submit', (e) => {
        e.preventDefault();
        
        const path = document.getElementById('album-path').value;
        const relativePath = document.getElementById('album-relative-path').value;
        
        if (!metadata.albums) metadata.albums = {};
        
        metadata.albums[relativePath] = {
            title: document.getElementById('album-title').value,
            description: document.getElementById('album-description').value,

            cover_photo: document.getElementById('album-cover').value
        };
        
        // Update local albums data
        const album = albums.find(a => a.path === path);
        if (album) {
            album.title = metadata.albums[relativePath].title;
            album.description = metadata.albums[relativePath].description;
            album.coverPhoto = metadata.albums[relativePath].cover_photo;
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
    
    // Generate button
    document.getElementById('generateBtn').addEventListener('click', generateGallery);
    
    // View Gallery button
    document.getElementById('viewBtn').addEventListener('click', viewGallery);
    
    // Deployment buttons - attach to all deployment save/test/deploy buttons
    document.querySelectorAll('.deploy-save-btn').forEach(btn => {
        btn.addEventListener('click', saveDeployConfig);
    });
    document.querySelectorAll('.deploy-test-btn').forEach(btn => {
        btn.addEventListener('click', () => deployGallery(true));
    });
    document.querySelectorAll('.deploy-deploy-btn').forEach(btn => {
        btn.addEventListener('click', () => deployGallery(false));
    });
    
    // Deploy tab switching
    document.querySelectorAll('.deploy-tab-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const tabName = e.target.getAttribute('data-deploy-tab');
            
            // Update active tab button
            document.querySelectorAll('.deploy-tab-btn').forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');
            
            // Update active pane
            document.querySelectorAll('.deploy-pane').forEach(pane => pane.classList.remove('active'));
            document.getElementById(tabName + '-deploy-pane').classList.add('active');
        });
    });
    
    // Modal close on background click
    document.getElementById('album-modal').addEventListener('click', (e) => {
        if (e.target.id === 'album-modal') closeAlbumModal();
    });
    
    document.getElementById('photo-modal').addEventListener('click', (e) => {
        if (e.target.id === 'photo-modal') closePhotoModal();
    });
});

// Generate gallery with progress tracking
async function generateGallery() {
    const generateBtn = document.getElementById('generateBtn');
    
    // Save any unsaved changes first
    if (hasUnsavedChanges) {
        await saveAll();
    }
    
    // Update button state
    generateBtn.disabled = true;
    generateBtn.classList.add('generating');
    generateBtn.innerHTML = '<div class="progress-bar"></div><span class="btn-text">Generating...</span>';
    
    try {
        // Start generation
        const response = await fetch('/api/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            throw new Error('Generation failed');
        }
        
        // Poll for progress
        const progressBar = generateBtn.querySelector('.progress-bar');
        let progress = 0;
        
        const pollProgress = async () => {
            try {
                const progressResponse = await fetch('/api/generate/progress');
                const data = await progressResponse.json();
                
                progress = data.progress || 0;
                progressBar.style.width = progress + '%';
                
                // Update button text with progress
                const btnText = generateBtn.querySelector('.btn-text');
                if (btnText) {
                    btnText.textContent = 'Generating... ' + progress + '%';
                }
                
                if (data.status !== 'completed' && data.status !== 'error' && data.status !== 'idle') {
                    setTimeout(pollProgress, 500);
                } else if (data.status === 'completed') {
                    progressBar.style.width = '100%';
                    setTimeout(() => {
                        generateBtn.classList.remove('generating');
                        generateBtn.innerHTML = 'Gallery Generated!';
                        generateBtn.style.background = 'var(--success-green-dark)';
                        generateBtn.style.borderColor = 'var(--success-green-dark)';
                        
                        setTimeout(() => {
                            generateBtn.innerHTML = 'Generate Gallery';
                            generateBtn.style.background = '';
                            generateBtn.style.borderColor = '';
                            generateBtn.disabled = false;
                        }, 3000);
                    }, 500);
                } else if (data.status === 'error') {
                    throw new Error(data.error || 'Generation failed');
                }
            } catch (error) {
                console.error('Error polling progress:', error);
                throw error;
            }
        };
        
        // Start polling after a short delay
        setTimeout(pollProgress, 500);
        
    } catch (error) {
        console.error('Error generating gallery:', error);
        generateBtn.classList.remove('generating');
        generateBtn.innerHTML = 'Generation Failed';
        generateBtn.style.background = 'var(--error-red)';
        generateBtn.style.borderColor = 'var(--error-red)';
        
        setTimeout(() => {
            generateBtn.innerHTML = 'Generate Gallery';
            generateBtn.style.background = '';
            generateBtn.style.borderColor = '';
            generateBtn.disabled = false;
        }, 3000);
    }
}

// View generated gallery
function viewGallery() {
    // Open the gallery served through the editor server
    window.open('/gallery/', '_blank');
}

// Load deployment configuration
async function loadDeployConfig() {
    try {
        const response = await fetch('/api/deploy-config');
        if (!response.ok) {
            console.error('Failed to load deploy config:', response.status, response.statusText);
            return;
        }
        
        const config = await response.json();
        
        // Load rsync configuration
        if (config.rsync) {
            document.getElementById('rsync-host').value = config.rsync.host || '';
            document.getElementById('rsync-path').value = config.rsync.path || '';
            document.getElementById('rsync-port').value = config.rsync.port || 22;
        }
        
        // Load S3 configuration
        if (config.s3) {
            document.getElementById('s3-bucket').value = config.s3.bucket || '';
            document.getElementById('s3-region').value = config.s3.region || '';
            // Don't load AWS credentials - they come from environment
        }
        
        // Load Cloudflare configuration
        if (config.cloudflare) {
            document.getElementById('cf-project').value = config.cloudflare.project || '';
            document.getElementById('cf-account').value = config.cloudflare.account_id || '';
            // Don't load API token - it comes from environment
        }
    } catch (error) {
        console.error('Error loading deployment config:', error);
    }
}

// Save deployment configuration
async function saveDeployConfig() {
    // Get the active deployment tab to determine which config to save
    const activeTab = document.querySelector('.deploy-tab-btn.active').getAttribute('data-deploy-tab');
    
    const config = {
        rsync: {
            host: document.getElementById('rsync-host').value,
            path: document.getElementById('rsync-path').value,
            port: parseInt(document.getElementById('rsync-port').value) || 22
        },
        s3: {
            bucket: document.getElementById('s3-bucket').value,
            region: document.getElementById('s3-region').value
            // AWS credentials should come from environment variables
        },
        cloudflare: {
            project: document.getElementById('cf-project').value,
            account_id: document.getElementById('cf-account').value
            // API token should come from environment variables
        }
    };
    
    // Get the save button for the active tab
    const saveBtn = document.querySelector('#' + activeTab + '-deploy-pane .btn-primary');
    const originalText = saveBtn.textContent;
    
    saveBtn.disabled = true;
    saveBtn.textContent = 'Saving...';
    
    try {
        const response = await fetch('/api/deploy-config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
        });
        
        if (response.ok) {
            saveBtn.textContent = 'Configuration Saved!';
            saveBtn.style.background = 'var(--success-green)';
            saveBtn.style.borderColor = 'var(--success-green)';
            
            setTimeout(() => {
                saveBtn.textContent = originalText;
                saveBtn.style.background = '';
                saveBtn.style.borderColor = '';
                saveBtn.disabled = false;
            }, 2000);
        } else {
            saveBtn.textContent = 'Save Failed';
            saveBtn.style.background = '#dc3545';
            saveBtn.style.borderColor = '#dc3545';
            
            setTimeout(() => {
                saveBtn.textContent = originalText;
                saveBtn.style.background = '';
                saveBtn.style.borderColor = '';
                saveBtn.disabled = false;
            }, 2000);
        }
    } catch (error) {
        console.error('Error saving deployment config:', error);
        saveBtn.textContent = 'Save Failed';
        saveBtn.style.background = 'var(--error-red)';
        saveBtn.style.borderColor = 'var(--error-red)';
        
        setTimeout(() => {
            saveBtn.textContent = originalText;
            saveBtn.style.background = '';
            saveBtn.style.borderColor = '';
            saveBtn.disabled = false;
        }, 2000);
    }
}

// Deploy gallery
async function deployGallery(dryRun) {
    // Get the active deployment tab
    const activeTab = document.querySelector('.deploy-tab-btn.active').getAttribute('data-deploy-tab');
    
    // Get the appropriate button based on the active tab and action
    const buttonClass = dryRun ? '.deploy-test-btn' : '.deploy-deploy-btn';
    const deployBtn = document.querySelector('#' + activeTab + '-deploy-pane ' + buttonClass);
    
    if (!deployBtn) {
        console.error('Deploy button not found');
        return;
    }
    
    const originalText = deployBtn.querySelector('span').textContent;
    const buttonSpan = deployBtn.querySelector('span');
    
    // Start the animation
    deployBtn.classList.add('deploying');
    buttonSpan.textContent = dryRun ? 'Testing...' : 'Deploying...';
    deployBtn.disabled = true;
    
    try {
        const response = await fetch('/api/deploy', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                target: activeTab,
                dry_run: dryRun
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // For dry runs, just show success
            if (dryRun) {
                deployBtn.classList.remove('deploying');
                deployBtn.classList.add('deployed');
                buttonSpan.textContent = 'Test Successful!';
                
                setTimeout(() => {
                    deployBtn.classList.remove('deployed');
                    buttonSpan.textContent = originalText;
                    deployBtn.disabled = false;
                }, 3000);
            } else {
                // For actual deployments, show progress
                // Poll for progress
                let progress = 0;
                
                const pollProgress = async () => {
                    try {
                        const progressResponse = await fetch('/api/deploy/progress');
                        const data = await progressResponse.json();
                        
                        progress = data.progress || 0;
                        // Update progress text
                        buttonSpan.textContent = 'Deploying... ' + Math.round(progress) + '%';
                        
                        if (data.status !== 'completed' && data.status !== 'error' && data.status !== 'idle') {
                            setTimeout(pollProgress, 500);
                        } else if (data.status === 'completed') {
                            // Transition to success state
                            deployBtn.classList.remove('deploying');
                            deployBtn.classList.add('deployed');
                            buttonSpan.textContent = 'Deploy Complete!';
                            
                            setTimeout(() => {
                                deployBtn.classList.remove('deployed');
                                buttonSpan.textContent = originalText;
                                deployBtn.disabled = false;
                            }, 3000);
                        } else if (data.status === 'error') {
                            throw new Error(data.error || 'Deployment failed');
                        }
                    } catch (error) {
                        console.error('Error polling progress:', error);
                        deployBtn.classList.remove('deploying');
                        deployBtn.classList.add('error');
                        buttonSpan.textContent = 'Deploy Failed';
                        
                        setTimeout(() => {
                            deployBtn.classList.remove('error');
                            buttonSpan.textContent = originalText;
                            deployBtn.disabled = false;
                        }, 3000);
                    }
                };
                
                // Start polling after a short delay
                setTimeout(pollProgress, 500);
            }
        } else {
            const error = await response.text();
            console.error('Deployment error:', error);
            
            // Check if it's the output directory error
            if (error.includes('Output directory not found')) {
                deployBtn.classList.remove('deploying');
                buttonSpan.textContent = originalText;
                deployBtn.disabled = false;
                showDeployNotice();
            } else {
                deployBtn.classList.remove('deploying');
                deployBtn.classList.add('error');
                buttonSpan.textContent = dryRun ? 'Test Failed' : 'Deploy Failed';
                
                setTimeout(() => {
                    deployBtn.classList.remove('error');
                    buttonSpan.textContent = originalText;
                    deployBtn.disabled = false;
                }, 3000);
            }
        }
    } catch (error) {
        console.error('Error deploying:', error);
        deployBtn.classList.remove('deploying');
        deployBtn.classList.add('error');
        buttonSpan.textContent = 'Error!';
        
        setTimeout(() => {
            deployBtn.classList.remove('error');
            buttonSpan.textContent = originalText;
            deployBtn.disabled = false;
        }, 3000);
    }
}
`
