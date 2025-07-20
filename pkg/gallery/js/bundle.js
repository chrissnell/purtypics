// Bundled Purtypics Gallery JavaScript (ES5 compatible)

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

// Lightbox class
class Lightbox {
  constructor() {
    this.photos = [];
    this.currentIndex = 0;
    this.lightbox = null;
    this.img = null;
    this.video = null;
    this.info = null;
    this.exif = null;
  }

  init() {
    this.photos = Array.from(document.querySelectorAll('.photo-item'));
    this.createLightbox();
    this.bindEvents();
  }

  createLightbox() {
    this.lightbox = document.createElement('div');
    this.lightbox.className = 'lightbox';
    this.lightbox.innerHTML = `
      <div class="lightbox-content">
        <span class="lightbox-close">&times;</span>
        <div class="lightbox-image-container">
          <img src="" alt="" style="display:none;">
          <video controls style="display:none; max-width: 100%; max-height: calc(100vh - 8rem);" preload="metadata"></video>
          <div class="lightbox-info"></div>
          <div class="lightbox-exif"></div>
        </div>
      </div>
      <button class="lightbox-prev">‹</button>
      <button class="lightbox-next">›</button>
    `;
    
    document.body.appendChild(this.lightbox);
    
    this.img = this.lightbox.querySelector('img');
    this.video = this.lightbox.querySelector('video');
    this.info = this.lightbox.querySelector('.lightbox-info');
    this.exif = this.lightbox.querySelector('.lightbox-exif');
  }

  bindEvents() {
    // Close button
    this.lightbox.querySelector('.lightbox-close').addEventListener('click', () => this.close());
    
    // Click outside image to close
    this.lightbox.addEventListener('click', (e) => {
      if (e.target === this.lightbox || e.target.classList.contains('lightbox-content')) {
        this.close();
      }
    });
    
    // Navigation
    this.lightbox.querySelector('.lightbox-prev').addEventListener('click', () => this.navigate(-1));
    this.lightbox.querySelector('.lightbox-next').addEventListener('click', () => this.navigate(1));
    
    // Keyboard navigation
    document.addEventListener('keydown', (e) => {
      if (!this.isOpen()) return;
      
      switch(e.key) {
        case 'Escape':
          this.close();
          break;
        case 'ArrowLeft':
          this.navigate(-1);
          break;
        case 'ArrowRight':
          this.navigate(1);
          break;
      }
    });
    
    // Photo click events
    this.photos.forEach((photo, index) => {
      photo.addEventListener('click', (e) => {
        e.preventDefault();
        this.open(index);
      });
    });
  }

  open(index) {
    this.currentIndex = index;
    this.showPhoto(this.photos[index]);
    this.lightbox.style.display = 'flex';
    document.body.style.overflow = 'hidden';
    this.updateNavigation();
  }

  close() {
    this.lightbox.style.display = 'none';
    document.body.style.overflow = '';
    this.video.pause();
    this.video.src = '';
  }

  navigate(direction) {
    this.currentIndex = (this.currentIndex + direction + this.photos.length) % this.photos.length;
    this.showPhoto(this.photos[this.currentIndex]);
    this.updateNavigation();
  }

  showPhoto(photo) {
    const fullSrc = photo.dataset.fullSrc || photo.querySelector('img').src;
    const title = photo.dataset.title || photo.querySelector('img').alt || '';
    const isVideo = photo.classList.contains('video-item');
    
    // Show/hide media elements
    if (isVideo) {
      this.img.style.display = 'none';
      this.video.style.display = 'block';
      this.video.src = photo.dataset.videoSrc;
    } else {
      this.video.style.display = 'none';
      this.video.pause();
      this.video.src = '';
      this.img.style.display = 'block';
      this.img.src = fullSrc;
    }
    
    // Update info
    this.info.innerHTML = title ? `<h3>${title}</h3>` : '';
    
    // Update EXIF data
    if (photo.dataset.camera || photo.dataset.datetime) {
      let exifHtml = '<div class="exif-data">';
      
      if (photo.dataset.camera) exifHtml += `<span>${photo.dataset.camera}</span>`;
      if (photo.dataset.datetime) exifHtml += `<span>${photo.dataset.datetime}</span>`;
      
      exifHtml += '</div>';
      this.exif.innerHTML = exifHtml;
    } else {
      this.exif.innerHTML = '';
    }
  }

  updateNavigation() {
    const prev = this.lightbox.querySelector('.lightbox-prev');
    const next = this.lightbox.querySelector('.lightbox-next');
    
    prev.style.display = this.photos.length > 1 ? 'block' : 'none';
    next.style.display = this.photos.length > 1 ? 'block' : 'none';
  }

  isOpen() {
    return this.lightbox.style.display === 'flex';
  }
}

// Lazy loading functionality
function initLazyLoading() {
  const lazyImages = document.querySelectorAll('img[loading="lazy"]');
  
  if ('loading' in HTMLImageElement.prototype) {
    // Browser supports native lazy loading
    return;
  }
  
  // Fallback for browsers that don't support native lazy loading
  const imageObserver = new IntersectionObserver((entries, observer) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = entry.target;
        img.src = img.dataset.src;
        img.classList.remove('lazy');
        imageObserver.unobserve(img);
      }
    });
  });
  
  lazyImages.forEach(img => imageObserver.observe(img));
}

// Masonry layout functionality
function initMasonry() {
  const grid = document.querySelector('.masonry-grid');
  if (!grid) return null;
  
  // Initialize Masonry
  const msnry = new Masonry(grid, {
    itemSelector: '.photo-item',
    columnWidth: '.photo-item',
    gutter: 3,
    fitWidth: true,
    transitionDuration: 0
  });
  
  // Layout after all images have loaded
  imagesLoaded(grid, () => {
    msnry.layout();
  });
  
  // Handle image load events for lazy loaded images
  grid.addEventListener('load', (e) => {
    if (e.target.tagName === 'IMG') {
      msnry.layout();
    }
  }, true);
  
  return msnry;
}

// Album grid layout
function initAlbumGrid() {
  const albumGrid = document.querySelector('.album-grid');
  if (!albumGrid) return null;
  
  const msnry = new Masonry(albumGrid, {
    itemSelector: '.album-item',
    columnWidth: '.album-item',
    gutter: 3,
    fitWidth: true,
    transitionDuration: 0
  });
  
  imagesLoaded(albumGrid, () => {
    msnry.layout();
  });
  
  return msnry;
}

// Photo map functionality
class PhotoMap {
  constructor(options = {}) {
    this.mapId = options.mapId || 'map';
    this.map = null;
    this.markers = [];
  }

  init() {
    const mapElement = document.getElementById(this.mapId);
    if (!mapElement) return;
    
    // Initialize Leaflet map
    this.map = L.map(this.mapId).setView([0, 0], 2);
    
    // Add tile layer
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
      maxZoom: 19
    }).addTo(this.map);
    
    // Add photo markers
    this.addPhotoMarkers();
    
    // Fit bounds to show all markers
    if (this.markers.length > 0) {
      const group = new L.featureGroup(this.markers);
      this.map.fitBounds(group.getBounds().pad(0.1));
    }
  }

  addPhotoMarkers() {
    const photos = document.querySelectorAll('.photo-item[data-lat][data-lng]');
    
    photos.forEach(photo => {
      const lat = parseFloat(photo.dataset.lat);
      const lng = parseFloat(photo.dataset.lng);
      const title = photo.dataset.title || 'Photo';
      const thumb = photo.querySelector('img').src;
      
      if (!isNaN(lat) && !isNaN(lng)) {
        // Create custom icon
        const icon = L.divIcon({
          className: 'photo-marker',
          html: `<img src="${thumb}" alt="${title}">`,
          iconSize: [60, 60],
          iconAnchor: [30, 30]
        });
        
        // Create marker
        const marker = L.marker([lat, lng], { icon });
        
        // Add popup
        marker.bindPopup(`
          <div class="map-popup">
            <img src="${thumb}" alt="${title}">
            <h4>${title}</h4>
          </div>
        `);
        
        // Add click handler to open in lightbox
        marker.on('click', () => {
          const index = Array.from(photos).indexOf(photo);
          if (window.lightbox && index >= 0) {
            window.lightbox.open(index);
          }
        });
        
        marker.addTo(this.map);
        this.markers.push(marker);
      }
    });
  }
}

// Main gallery initialization
function initGallery() {
  // Initialize video hover effects
  initVideoHover();
  
  // Initialize lightbox
  window.lightbox = new Lightbox();
  window.lightbox.init();
  
  // Initialize lazy loading
  initLazyLoading();
  
  // Initialize masonry layout
  window.masonryInstance = initMasonry() || initAlbumGrid();
  
  // Initialize map if present
  if (document.getElementById('map')) {
    window.photoMap = new PhotoMap();
    window.photoMap.init();
  }
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initGallery);
} else {
  initGallery();
}