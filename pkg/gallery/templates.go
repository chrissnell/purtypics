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
  font-size: 2.5rem;
  font-weight: 500;
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
  padding: 0 2rem 2rem 2rem;
  text-align: left; /* Left-aligns the masonry grid */
}

.masonry-grid {
  margin: 0 auto;
  width: 100%;
}

/* Album grid container */
.album-grid {
  margin: 0;
  width: 100%;
}

/* Grid sizer for Masonry - defines column width */
.grid-sizer,
.photo-item {
  width: 20%; /* 5 columns on desktop */
}

/* Album sizer not needed - albums use CSS Grid */

.photo-item {
  margin-bottom: 3px;
  float: left;
}

/* Reduce margin on mobile for better space usage */
@media (max-width: 600px) {
  .photo-item {
    margin-bottom: 2px;
  }
}

/* Responsive column widths */
@media (max-width: 1200px) {
  .grid-sizer,
  .photo-item {
    width: 25%; /* 4 columns on medium screens */
  }
}

@media (max-width: 900px) {
  .grid-sizer,
  .photo-item {
    width: 33.333%; /* 3 columns on tablets */
  }
}

@media (max-width: 600px) {
  .grid-sizer,
  .photo-item {
    width: 50%; /* 2 columns on large phones */
  }
}

/* Specific handling for iPhone Pro Max and similar devices */
@media (min-width: 414px) and (max-width: 600px) {
  .grid-sizer,
  .photo-item {
    width: 50%; /* Ensure 2 columns on iPhone Pro Max */
  }
}

@media (max-width: 380px) {
  .grid-sizer,
  .photo-item {
    width: 100%; /* 1 column on small phones */
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

/* Aspect ratio container to prevent reflow */
.photo-item .aspect-ratio-box {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.photo-item .aspect-ratio-box::before {
  content: "";
  display: block;
  padding-top: var(--aspect-ratio, 100%);
}

.photo-item .aspect-ratio-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
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

.photo-item .aspect-ratio-content img {
  position: absolute;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-title {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 0;
  background: rgba(0,0,0,0.5);
  color: #ffffff;
  font-family: var(--font-mono);
  font-size: 0.8rem;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.photo-item:hover .photo-title {
  opacity: 1;
}

/* Album Grid - Using Masonry.js */
.album-grid-container {
  padding: 2rem;
  text-align: left;
}

.album-item {
  width: calc(25% - 15px); /* 4 columns with gap */
  margin-bottom: 20px;
  border: 1px dotted var(--border-color);
  padding: 1rem;
  transition: all 0.2s ease;
  display: block;
}

@media (max-width: 1200px) {
  .album-item {
    width: calc(33.333% - 14px); /* 3 columns with gap */
  }
}

@media (max-width: 800px) {
  .album-item {
    width: calc(50% - 10px); /* 2 columns with gap */
  }
  
  .album-thumbnail {
    width: 100%;
    height: auto;
    object-fit: cover;
  }
}

@media (max-width: 400px) {
  .album-item {
    width: 100%; /* 1 column */
  }
  
  .album-thumbnail {
    width: 100%;
    height: auto;
    object-fit: cover;
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
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.album-description {
  font-size: 1rem;
  color: #555;
  margin-bottom: 0.5rem;
  line-height: 1.4;
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

.lightbox-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  font-size: 3rem;
  color: white;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.3s;
  padding: 1rem;
  user-select: none;
}

.lightbox-nav:hover {
  opacity: 1;
}

.lightbox-prev {
  left: 2rem;
}

.lightbox-next {
  right: 2rem;
}

/* Video styling */
.video-container {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.aspect-ratio-content .video-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.video-poster,
.video-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
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
  background: #f8f9fa;
}

.map-section-inner {
  margin: 0;
  width: 100%;
  max-width: none;
  position: relative;
}

.map-section h2 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 500;
  color: #17a2b8;
  text-transform: uppercase;
  font-family: var(--font-mono);
}

.photo-map {
  width: 100%;
  height: 500px;
  border: 1px dotted var(--border-color);
  background: #f0f0f0;
  border-radius: 0;
  overflow: hidden;
}

/* Leaflet popup customization */
.leaflet-popup-content {
  font-family: var(--font-main);
  font-size: 14px;
  padding: 10px;
}

.leaflet-popup-content strong {
  color: #17a2b8;
}

.leaflet-popup-content img {
  display: block;
  object-fit: cover;
  border: 1px dotted var(--border-color);
}

.leaflet-popup-content-wrapper {
  border-radius: 0;
  border: 1px dotted var(--border-color);
}

.leaflet-popup-tip {
  background: white;
  border: 1px dotted var(--border-color);
  border-top: none;
  border-right: none;
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
  /* Remove conflicting grid styles for masonry */
  
  .header {
    padding: 1rem;
  }
  
  .breadcrumbs {
    padding: 0.5rem 1rem;
  }
  
  .masonry-container {
    padding: 0 1rem 1rem 1rem;
  }
  
  .album-grid-container {
    padding: 1rem;
  }
  
  /* Let Masonry handle the layout */
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
      const playButton = item.querySelector('.play-button');
      if (playButton) {
        playButton.style.display = 'none';
      }
    });
    
    item.addEventListener('mouseleave', () => {
      video.style.opacity = '0';
      video.pause();
      video.currentTime = 0;
      const playButton = item.querySelector('.play-button');
      if (playButton) {
        playButton.style.display = '';
      }
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
      <div class="lightbox-nav lightbox-prev">&lt;</div>
      <div class="lightbox-nav lightbox-next">&gt;</div>
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
  const prevBtn = lightbox.querySelector('.lightbox-prev');
  const nextBtn = lightbox.querySelector('.lightbox-next');
  
  let currentPhotoIndex = 0;
  
  const showPhoto = (index) => {
    const photo = photos[index];
    if (!photo) return;
    
    currentPhotoIndex = index;
    
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
  };
  
  photos.forEach((photo, index) => {
    photo.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      showPhoto(index);
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
  
  const navigatePhoto = (direction) => {
    let newIndex = currentPhotoIndex + direction;
    
    // Wrap around at the ends
    if (newIndex < 0) {
      newIndex = photos.length - 1;
    } else if (newIndex >= photos.length) {
      newIndex = 0;
    }
    
    showPhoto(newIndex);
  };
  
  closeBtn.addEventListener('click', closeLightbox);
  prevBtn.addEventListener('click', () => navigatePhoto(-1));
  nextBtn.addEventListener('click', () => navigatePhoto(1));
  
  lightbox.addEventListener('click', (e) => {
    if (e.target === lightbox || e.target === lightbox.querySelector('.lightbox-content')) {
      closeLightbox();
    }
  });
  
  // Keyboard navigation
  document.addEventListener('keydown', (e) => {
    if (lightbox.style.display === 'block') {
      switch(e.key) {
        case 'Escape':
          closeLightbox();
          break;
        case 'ArrowLeft':
          navigatePhoto(-1);
          break;
        case 'ArrowRight':
          navigatePhoto(1);
          break;
      }
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
    // Add grid-sizer element if it doesn't exist
    if (!photoGrid.querySelector('.grid-sizer')) {
      const gridSizer = document.createElement('div');
      gridSizer.className = 'grid-sizer';
      photoGrid.prepend(gridSizer);
    }
    
    imagesLoaded(photoGrid, function() {
      // Use smaller gutter on mobile devices
      const gutterSize = window.innerWidth <= 600 ? 2 : 3;
      
      const msnry = new Masonry(photoGrid, {
        itemSelector: '.photo-item',
        columnWidth: '.grid-sizer',
        percentPosition: true,
        gutter: gutterSize,
        fitWidth: false,
        stagger: 30
      });
      
      // Store masonry instance globally for resize handling
      window.photoMsnry = msnry;
      
      // Adjust map width after layout is complete
      msnry.on('layoutComplete', function() {
        adjustMapWidth();
      });
      
      // Layout after all images loaded
      msnry.layout();
    });
  }

  // Album grid masonry
  const albumGrid = document.querySelector('.album-grid');
  if (albumGrid) {
    imagesLoaded(albumGrid, function() {
      const msnry = new Masonry(albumGrid, {
        itemSelector: '.album-item',
        columnWidth: '.album-item',
        percentPosition: true,
        gutter: 0
      });
    });
  }
}

// Adjust map section width to match photo grid full width
function adjustMapWidth() {
  const masonryContainer = document.querySelector('.masonry-container');
  const mapSectionInner = document.querySelector('.map-section-inner');
  const photoItems = document.querySelectorAll('.photo-item');
  
  if (masonryContainer && mapSectionInner && photoItems.length > 0) {
    // Calculate the full width from leftmost to rightmost photo
    let minLeft = Infinity;
    let maxRight = -Infinity;
    
    photoItems.forEach(item => {
      const rect = item.getBoundingClientRect();
      minLeft = Math.min(minLeft, rect.left);
      maxRight = Math.max(maxRight, rect.right);
    });
    
    const fullWidth = maxRight - minLeft;
    
    // Set the map width
    mapSectionInner.style.width = fullWidth + 'px';
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
      const img = photo.querySelector('img');
      const videoPoster = photo.querySelector('.video-poster');
      photosWithGPS.push({
        lat: lat,
        lng: lng,
        title: (img || videoPoster).alt,
        element: photo,
        isVideo: photo.dataset.video === 'true'
      });
    }
  });
  
  if (photosWithGPS.length === 0) {
    mapElement.innerHTML = '<p style="text-align: center; padding: 2rem; color: #666;">No location data available for these photos</p>';
    return;
  }
  
  // Initialize Leaflet map
  const map = L.map('map').setView([photosWithGPS[0].lat, photosWithGPS[0].lng], 10);
  
  // Add OpenStreetMap tile layer
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  }).addTo(map);
  
  // Create markers for each photo
  const markers = [];
  photosWithGPS.forEach(photo => {
    // Get the thumbnail URL from the photo element
    const img = photo.element.querySelector('img');
    const videoPoster = photo.element.querySelector('.video-poster');
    const thumbSrc = (img || videoPoster).src;
    
    // Create popup content with thumbnail
    const popupContent = ` + "`" + `
      <div style="text-align: center; min-width: 200px;">
        <img src="${thumbSrc}" alt="${photo.title}" style="max-width: 200px; max-height: 150px; border-radius: 4px; margin-bottom: 8px;">
        <div><strong>${photo.title}</strong>${photo.isVideo ? ' (Video)' : ''}</div>
      </div>
    ` + "`" + `;
    
    const marker = L.marker([photo.lat, photo.lng])
      .bindPopup(popupContent, {
        maxWidth: 250,
        minWidth: 200
      })
      .addTo(map);
    
    // Add click handler to scroll to photo when marker is clicked
    marker.on('popupopen', () => {
      photo.element.scrollIntoView({ behavior: 'smooth', block: 'center' });
      // Highlight the photo briefly
      photo.element.style.transition = 'outline 0.3s';
      photo.element.style.outline = '3px solid #17a2b8';
      setTimeout(() => {
        photo.element.style.outline = '';
      }, 2000);
    });
    
    markers.push(marker);
  });
  
  // Fit map to show all markers
  if (photosWithGPS.length > 1) {
    const group = new L.featureGroup(markers);
    map.fitBounds(group.getBounds().pad(0.1));
  }
}

// Initialize everything
document.addEventListener('DOMContentLoaded', () => {
  initVideoHover();
  initLightbox();
  initLazyLoad();
  initMasonry();
  initMap();
});

// Re-initialize Masonry on window resize with debounce
let resizeTimer;
window.addEventListener('resize', function() {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(function() {
    if (window.photoMsnry) {
      // Update gutter size based on new window width
      const newGutterSize = window.innerWidth <= 600 ? 2 : 3;
      window.photoMsnry.gutter = newGutterSize;
      window.photoMsnry.layout();
    }
  }, 250);
});
`

// HTML templates
const baseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{if ne .Title .SiteTitle}}{{.SiteTitle}} :: {{.Title}}{{else}}{{.Title}}{{end}}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inconsolata:wght@400;700&family=Roboto+Slab:wght@300;400;700&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="{{if .BaseURL}}{{.BaseURL}}/{{end}}css/style.css">
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
  <script src="https://unpkg.com/masonry-layout@4/dist/masonry.pkgd.min.js"></script>
  <script src="https://unpkg.com/imagesloaded@5/imagesloaded.pkgd.min.js"></script>
  <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
</head>
<body>
  <header class="header">
    <a href="{{.BaseURL}}/" class="site-title">{{.SiteTitle}}</a>
  </header>
  
  {{if .Breadcrumbs}}
  <nav class="breadcrumbs">
    {{range $i, $crumb := .Breadcrumbs}}
      {{if $i}} :: {{end}}
      {{if $crumb.URL}}
        <a href="{{$crumb.URL}}">{{$crumb.Title}}</a>
      {{else}}
        <span>{{$crumb.Title}}</span>
      {{end}}
    {{end}}
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
      {{$album := .}}
      {{$coverPhoto := ""}}
      {{if $album.CoverPhoto}}
        {{range $album.Photos}}
          {{if eq .Filename $album.CoverPhoto}}
            {{$coverPhoto = .}}
          {{end}}
        {{end}}
      {{end}}
      {{if not $coverPhoto}}
        {{$coverPhoto = index $album.Photos 0}}
      {{end}}
      {{with $coverPhoto}}
      {{if index .Thumbnails "medium"}}
      <img src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "medium"}}" 
           alt="{{.Title}}" 
           class="album-thumbnail"
           loading="lazy">
      {{else if index .Thumbnails "small"}}
      <img src="{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{index .Thumbnails "small"}}" 
           alt="{{.Title}}" 
           class="album-thumbnail"
           loading="lazy">
      {{else}}
      <div class="album-thumbnail" style="background: #f0f0f0; display: flex; align-items: center; justify-content: center; height: 200px; color: #999;">No thumbnail</div>
      {{end}}
      {{end}}
      {{end}}
      <h2 class="album-title">{{.Title}}</h2>
      {{if .Description}}<p class="album-description">{{.Description}}</p>{{end}}
      <p class="album-count">{{len .Photos}} {{if eq (len .Photos) 1}}photo{{else}}photos{{end}}</p>
    </a>
    {{end}}
  </div>
</div>
`

const albumContent = `<div class="masonry-container">
  <div class="masonry-grid">
    <div class="grid-sizer"></div>
    {{range .Album.Photos}}
    <div class="photo-item {{if .IsVideo}}video-item{{end}} {{if gt .Height .Width}}portrait{{end}}" 
         data-aspect="{{.AspectRatio}}" 
         data-lat="{{if .EXIF}}{{if .EXIF.GPS}}{{.EXIF.GPS.Latitude}}{{end}}{{end}}"
         data-lng="{{if .EXIF}}{{if .EXIF.GPS}}{{.EXIF.GPS.Longitude}}{{end}}{{end}}"
         data-camera="{{if .EXIF}}{{.EXIF.Camera}}{{end}}"
         data-datetime="{{if .EXIF}}{{.EXIF.DateTime.Format "Jan 2, 2006 at 3:04PM"}}{{end}}"
         data-video="{{if .IsVideo}}true{{end}}"
         data-video-src="{{if .IsVideo}}{{if $.BaseURL}}{{$.BaseURL}}{{end}}{{.VideoPath}}{{end}}"
         data-width="{{.Width}}"
         data-height="{{.Height}}"
         style="--aspect-ratio: {{if and .Width .Height}}{{if gt .Width 0}}calc({{.Height}} / {{.Width}} * 100%){{else}}100%{{end}}{{else}}100%{{end}}">
      <div class="aspect-ratio-box">
        <div class="aspect-ratio-content">
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
        </div>
      </div>
      <div class="photo-title">{{.Title}}</div>
    </div>
    {{end}}
  </div>
</div>

{{if .HasGPS}}
<div class="map-section">
  <div class="map-section-inner">
    <h2>Photo Locations</h2>
    <div id="map" class="photo-map"></div>
  </div>
</div>
{{end}}
`
