package gallery

// CSS for the gallery
const galleryCSS = `/* Purtypics Gallery */
:root {
  --bg-color: #ffffff;
  --text-color: #333333;
  --border-color: #e0e0e0;
  --hover-color: #f5f5f5;
  --font-main: 'Roboto Slab', serif;
  --font-mono: 'Inconsolata', monospace;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  margin: 0;
  padding: 0;
}

body {
  margin: 0;
  padding: 0;
  font-family: var(--font-main);
  color: var(--text-color);
  background: var(--bg-color);
  line-height: 1.6;
}

a {
  color: inherit;
  text-decoration: none;
}

/* Header */
.header {
  padding: 2rem;
  border-bottom: 1px dotted var(--border-color);
}

.site-title {
  font-size: 2rem;
  font-weight: 300;
  color: var(--text-color);
}

/* Breadcrumbs */
.breadcrumbs {
  padding: 1rem 2rem;
  font-family: var(--font-mono);
  font-size: 0.9rem;
}

.breadcrumbs a {
  border-bottom: 1px dotted var(--border-color);
}

.breadcrumbs a:hover {
  border-bottom-style: solid;
}

/* Main content wrapper */
main {
  padding-bottom: 2rem;
}

/* Masonry Grid - Using Masonry.js */
.masonry-container {
  padding: 2rem;
  text-align: center; /* Centers the masonry grid */
}

.masonry-grid {
  margin: 0 auto;
  display: inline-block; /* Allows fitWidth to work properly */
}

/* Album grid container */
.album-grid {
  margin: 0 auto;
  display: inline-block;
}

/* Photo items for masonry.js - fixed width for automatic column adjustment */
.photo-item {
  width: 280px;
  margin-bottom: 20px;
}

@media (max-width: 600px) {
  .photo-item {
    width: calc(50% - 10px);
  }
}

@media (max-width: 400px) {
  .photo-item {
    width: 100%;
  }
}

/* Photo Items styling */
.photo-item {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--border-color);
  transition: all 0.2s ease;
  cursor: pointer;
  background: #f8f8f8;
  display: block;
}

.photo-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.photo-item img {
  width: 100%;
  height: auto;
  display: block;
}

.photo-title {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 0.4rem 0.6rem;
  background: rgba(0,0,0,0.75);
  color: #ffffff;
  font-family: var(--font-mono);
  font-size: 0.85rem;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.photo-item:hover .photo-title {
  opacity: 1;
}

/* Album Grid - Using Masonry.js */
.album-grid-container {
  padding: 2rem;
  text-align: center;
}

.album-item {
  width: 380px;
  margin-bottom: 20px;
  border: 1px dotted var(--border-color);
  padding: 1rem;
  transition: all 0.2s ease;
  display: block;
}

@media (max-width: 800px) {
  .album-item {
    width: calc(50% - 10px);
  }
}

@media (max-width: 400px) {
  .album-item {
    width: 100%;
  }
}

.album-item:hover {
  background-color: var(--hover-color);
  transform: translateY(-2px);
}

.album-thumbnail {
  width: 100%;
  height: auto;
  margin-bottom: 1rem;
  border: 1px solid var(--border-color);
  display: block;
}

.album-title {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
}

.album-count {
  font-family: var(--font-mono);
  font-size: 0.9rem;
  color: #666;
}

/* Lightbox - Fixed to fill whole window */
.lightbox {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.95);
  z-index: 9999;
  cursor: pointer;
}

.lightbox-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  box-sizing: border-box;
}

.lightbox-image-container {
  max-width: 90%;
  max-height: 90%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.lightbox img {
  max-width: 100%;
  max-height: calc(100vh - 8rem);
  object-fit: contain;
  display: block;
}

.lightbox-close {
  position: absolute;
  top: 2rem;
  right: 2rem;
  font-size: 2rem;
  color: white;
  font-family: var(--font-mono);
  z-index: 1001;
}

.lightbox-info {
  color: white;
  font-family: var(--font-mono);
  font-size: 1rem;
  text-align: center;
  margin-top: 1rem;
}

.lightbox-exif {
  color: rgba(255,255,255,0.6);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  text-align: center;
  margin-top: 0.5rem;
}

/* Video styling */
.video-container {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.video-poster,
.video-preview {
  width: 100%;
  height: auto;
  display: block;
}

.video-preview {
  position: absolute;
  top: 0;
  left: 0;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.photo-item.video-item:hover .video-preview {
  opacity: 1;
}

.play-button {
  position: absolute;
  bottom: 1rem;
  left: 1rem;
  pointer-events: none;
  transition: transform 0.2s ease;
}

.photo-item.video-item:hover .play-button {
  transform: scale(1.1);
  opacity: 0.8;
}

/* Map section */
.map-section {
  padding: 2rem;
  border-top: 1px dotted var(--border-color);
}

.map-section h2 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 300;
}

.photo-map {
  width: 100%;
  height: 500px;
  border: 1px solid var(--border-color);
  background: #f0f0f0;
}

/* EXIF Info */
.exif-info {
  margin-top: 2rem;
  padding: 1rem;
  border: 1px dotted var(--border-color);
  font-family: var(--font-mono);
  font-size: 0.85rem;
}

/* Loading */
.loading {
  text-align: center;
  padding: 2rem;
  font-family: var(--font-mono);
  color: #666;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .masonry-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    grid-gap: 1rem;
  }
  
  .album-grid {
    grid-template-columns: 1fr;
  }
  
  .header {
    padding: 1rem;
  }
  
  .breadcrumbs {
    padding: 0.5rem 1rem;
  }
  
  .masonry-container {
    padding: 1rem;
  }
}`

// JavaScript for the gallery
const galleryJS = `// Purtypics Gallery JavaScript

// Video hover functionality
function initVideoHover() {
  const videos = document.querySelectorAll('.video-item');
  
  videos.forEach(item => {
    const poster = item.querySelector('.video-poster');
    const video = item.querySelector('.video-preview');
    
    if (!video) return;
    
    item.addEventListener('mouseenter', () => {
      video.style.opacity = '1';
      video.play().catch(e => {
        console.error('Video play failed:', e);
      });
    });
    
    item.addEventListener('mouseleave', () => {
      video.style.opacity = '0';
      video.pause();
      video.currentTime = 0;
    });
  });
}

// Lightbox functionality
function initLightbox() {
  const photos = document.querySelectorAll('.photo-item');
  const lightbox = document.createElement('div');
  lightbox.className = 'lightbox';
  lightbox.innerHTML = ` + "`" + `
    <div class="lightbox-content">
      <span class="lightbox-close">&times;</span>
      <div class="lightbox-image-container">
        <img src="" alt="" style="display:none;">
        <video controls style="display:none; max-width: 100%; max-height: calc(100vh - 8rem);" preload="metadata"></video>
        <div class="lightbox-info"></div>
        <div class="lightbox-exif"></div>
      </div>
    </div>
  ` + "`" + `;
  document.body.appendChild(lightbox);
  
  const lightboxImg = lightbox.querySelector('img');
  const lightboxVideo = lightbox.querySelector('video');
  const lightboxInfo = lightbox.querySelector('.lightbox-info');
  const lightboxExif = lightbox.querySelector('.lightbox-exif');
  const closeBtn = lightbox.querySelector('.lightbox-close');
  
  photos.forEach(photo => {
    photo.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      
      const isVideo = photo.dataset.video === 'true';
      const title = photo.querySelector('img').alt;
      
      // Get EXIF data from photo element
      const camera = photo.dataset.camera;
      const datetime = photo.dataset.datetime;
      
      lightboxInfo.textContent = title;
      
      if (isVideo) {
        // Show video
        const videoSrc = photo.dataset.videoSrc;
        lightboxVideo.src = videoSrc;
        lightboxVideo.style.display = 'block';
        lightboxImg.style.display = 'none';
        
        // Pause any playing preview videos
        const previewVideo = photo.querySelector('.video-preview');
        if (previewVideo) {
          previewVideo.pause();
        }
      } else {
        // Show image
        const img = photo.querySelector('img');
        const fullSrc = img.dataset.full || img.src;
        lightboxImg.src = fullSrc;
        lightboxImg.style.display = 'block';
        lightboxVideo.style.display = 'none';
      }
      
      // Show EXIF info if available
      if (camera || datetime) {
        let exifText = '';
        
        // Format camera info
        if (camera && datetime) {
          // Clean up camera string - remove extra quotes
          const cleanCamera = camera.replace(/"/g, '');
          exifText = ` + "`" + `Taken with ${cleanCamera} on ${datetime}` + "`" + `;
        } else if (camera) {
          const cleanCamera = camera.replace(/"/g, '');
          exifText = ` + "`" + `Taken with ${cleanCamera}` + "`" + `;
        } else if (datetime) {
          exifText = ` + "`" + `Taken on ${datetime}` + "`" + `;
        }
        
        lightboxExif.textContent = exifText;
        lightboxExif.style.display = 'block';
      } else {
        lightboxExif.style.display = 'none';
      }
      
      lightbox.style.display = 'block';
      document.body.style.overflow = 'hidden';
    });
  });
  
  const closeLightbox = () => {
    lightbox.style.display = 'none';
    document.body.style.overflow = '';
    
    // Stop video if playing
    if (lightboxVideo.src) {
      lightboxVideo.pause();
      lightboxVideo.src = '';
    }
  };
  
  closeBtn.addEventListener('click', closeLightbox);
  lightbox.addEventListener('click', (e) => {
    if (e.target === lightbox || e.target === lightbox.querySelector('.lightbox-content')) {
      closeLightbox();
    }
  });
  
  // Keyboard navigation
  document.addEventListener('keydown', (e) => {
    if (lightbox.style.display === 'block' && e.key === 'Escape') {
      closeLightbox();
    }
  });
}

// Lazy loading for images
function initLazyLoad() {
  const images = document.querySelectorAll('img[loading="lazy"]');
  
  if ('loading' in HTMLImageElement.prototype) {
    // Browser supports native lazy loading
    return;
  }
  
  // Fallback for older browsers
  const imageObserver = new IntersectionObserver((entries, observer) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = entry.target;
        img.src = img.dataset.src;
        img.classList.add('loaded');
        observer.unobserve(img);
      }
    });
  });
  
  images.forEach(img => imageObserver.observe(img));
}

// Infinite scroll
let loading = false;
let currentPage = 1;
const photosPerPage = 50;

function initInfiniteScroll() {
  const grid = document.querySelector('.masonry-grid');
  if (!grid) return;
  
  const photos = Array.from(grid.querySelectorAll('.photo-item'));
  const totalPhotos = photos.length;
  
  // Initially hide photos beyond first page
  photos.forEach((photo, index) => {
    if (index >= photosPerPage) {
      photo.style.display = 'none';
    }
  });
  
  window.addEventListener('scroll', () => {
    if (loading) return;
    
    const scrollHeight = document.documentElement.scrollHeight;
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const clientHeight = window.innerHeight;
    
    if (scrollTop + clientHeight >= scrollHeight - 200) {
      loadMorePhotos(photos);
    }
  });
}

function loadMorePhotos(photos) {
  loading = true;
  
  const start = currentPage * photosPerPage;
  const end = Math.min(start + photosPerPage, photos.length);
  
  if (start >= photos.length) {
    loading = false;
    return;
  }
  
  // Show loading indicator
  const loadingDiv = document.createElement('div');
  loadingDiv.className = 'loading';
  loadingDiv.textContent = 'Loading more photos...';
  document.querySelector('.masonry-container').appendChild(loadingDiv);
  
  // Simulate loading delay
  setTimeout(() => {
    for (let i = start; i < end; i++) {
      photos[i].style.display = 'block';
    }
    
    currentPage++;
    loading = false;
    loadingDiv.remove();
    
    // Reinitialize lightbox for new photos
    initLightbox();
  }, 300);
}

// Initialize Masonry
function initMasonry() {
  // Photo gallery masonry
  const photoGrid = document.querySelector('.masonry-grid');
  if (photoGrid) {
    imagesLoaded(photoGrid, function() {
      const msnry = new Masonry(photoGrid, {
        itemSelector: '.photo-item',
        columnWidth: '.photo-item',
        gutter: 20,
        fitWidth: true,
        stagger: 30,
        resize: true
      });
      
      // Center the grid
      photoGrid.style.margin = '0 auto';
    });
  }

  // Album grid masonry
  const albumGrid = document.querySelector('.album-grid');
  if (albumGrid) {
    imagesLoaded(albumGrid, function() {
      const msnry = new Masonry(albumGrid, {
        itemSelector: '.album-item',
        columnWidth: '.album-item',
        gutter: 20,
        fitWidth: true,
        stagger: 30,
        resize: true
      });
      
      // Center the grid
      albumGrid.style.margin = '0 auto';
    });
  }
}

// Initialize map if GPS data exists
function initMap() {
  const mapElement = document.getElementById('map');
  if (!mapElement) return;
  
  // Collect all photos with GPS data
  const photosWithGPS = [];
  document.querySelectorAll('.photo-item').forEach(photo => {
    const lat = parseFloat(photo.dataset.lat);
    const lng = parseFloat(photo.dataset.lng);
    if (!isNaN(lat) && !isNaN(lng)) {
      photosWithGPS.push({
        lat: lat,
        lng: lng,
        title: photo.querySelector('img').alt,
        element: photo
      });
    }
  });
  
  if (photosWithGPS.length === 0) {
    mapElement.innerHTML = '<p style="text-align: center; padding: 2rem; color: #666;">No location data available for these photos</p>';
    return;
  }
  
  // For now, show a simple message
  // In production, you'd integrate with a mapping library like Leaflet
  mapElement.innerHTML = '<p style="text-align: center; padding: 2rem; color: #666;">' + 
    photosWithGPS.length + ' photos have location data. Map integration coming soon!</p>';
}

// Initialize everything
document.addEventListener('DOMContentLoaded', () => {
  initVideoHover();
  initLightbox();
  initLazyLoad();
  initMasonry();
  initMap();
});
`

// HTML templates
const baseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}} - {{.SiteTitle}}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;700&family=Roboto+Slab:wght@300;400;700&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="{{if .BaseURL}}{{.BaseURL}}/{{end}}css/style.css">
  <script src="https://unpkg.com/masonry-layout@4/dist/masonry.pkgd.min.js"></script>
  <script src="https://unpkg.com/imagesloaded@5/imagesloaded.pkgd.min.js"></script>
</head>
<body>
  <header class="header">
    <a href="{{.BaseURL}}/" class="site-title">{{.SiteTitle}}</a>
  </header>
  
  {{if .Breadcrumbs}}
  <nav class="breadcrumbs">
    <a href="{{.BaseURL}}/">Albums</a> /
    <span>{{.Title}}</span>
  </nav>
  {{end}}
  
  <main>
    {{.Content}}
  </main>
  
  <script src="{{.BaseURL}}/js/gallery.js"></script>
</body>
</html>
`

const indexContent = `<div class="album-grid-container">
  <div class="album-grid">
    {{range .Albums}}
    <a href="{{$.BaseURL}}/albums/{{.ID}}.html" class="album-item">
      {{if .Photos}}
      {{with index .Photos 0}}
      <img src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "medium"}}" 
           alt="{{.Title}}" 
           class="album-thumbnail"
           loading="lazy">
      {{end}}
      {{end}}
      <h2 class="album-title">{{.Title}}</h2>
      <p class="album-count">{{len .Photos}} photos</p>
    </a>
    {{end}}
  </div>
</div>
`

const albumContent = `<div class="masonry-container">
  <div class="masonry-grid">
    {{range .Album.Photos}}
    <div class="photo-item {{if .IsVideo}}video-item{{end}}" data-aspect="{{.AspectRatio}}" 
         data-lat="{{if .EXIF}}{{if .EXIF.GPS}}{{.EXIF.GPS.Latitude}}{{end}}{{end}}"
         data-lng="{{if .EXIF}}{{if .EXIF.GPS}}{{.EXIF.GPS.Longitude}}{{end}}{{end}}"
         data-camera="{{if .EXIF}}{{.EXIF.Camera}}{{end}}"
         data-datetime="{{if .EXIF}}{{.EXIF.DateTime.Format "Jan 2, 2006 at 3:04PM"}}{{end}}"
         data-video="{{if .IsVideo}}true{{end}}"
         data-video-src="{{if .IsVideo}}{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{.VideoPath}}{{end}}">
      {{if .IsVideo}}
      <div class="video-container">
        <img src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "poster"}}" 
             alt="{{.Title}}" 
             loading="lazy"
             class="video-poster">
        <video src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{.VideoPath}}" 
               muted 
               loop
               preload="none"
               class="video-preview">
        </video>
        <div class="play-button">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <circle cx="24" cy="24" r="22" fill="rgba(0,0,0,0.7)" stroke="white" stroke-width="2"/>
            <path d="M19 16L32 24L19 32V16Z" fill="white"/>
          </svg>
        </div>
      </div>
      {{else}}
      <img src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "medium"}}" 
           data-full="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "full"}}" 
           alt="{{.Title}}" 
           loading="lazy">
      {{end}}
      <div class="photo-title">{{.Title}}</div>
    </div>
    {{end}}
  </div>
</div>

{{if .HasGPS}}
<div class="map-section">
  <h2>Photo Locations</h2>
  <div id="map" class="photo-map"></div>
</div>
{{end}}
`