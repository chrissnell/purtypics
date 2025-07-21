// Lightbox functionality
export class Lightbox {
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
    this.photos = Array.from(document.querySelectorAll('.photo-card'));
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
      const link = photo.querySelector('.photo-link');
      if (link) {
        link.addEventListener('click', (e) => {
          e.preventDefault();
          this.open(index);
        });
      }
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
    const link = photo.querySelector('.photo-link');
    const isVideo = photo.classList.contains('video-item');
    const fullSrc = link ? link.getAttribute('href') : '';
    
    // Extract title and description from photo element
    const photoImg = photo.querySelector('img');
    const photoTitle = photo.querySelector('.photo-title');
    const title = photoTitle ? photoTitle.textContent : (photoImg ? photoImg.getAttribute('alt') : '');
    const description = ''; // We can add description support later if needed
    
    // Show/hide media elements
    if (isVideo) {
      this.img.style.display = 'none';
      this.video.style.display = 'block';
      this.video.src = photo.dataset.videoSrc || fullSrc;
    } else {
      this.video.style.display = 'none';
      this.video.pause();
      this.video.src = '';
      this.img.style.display = 'block';
      this.img.src = fullSrc;
    }
    
    // Update info
    this.info.innerHTML = title ? `<h3>${title}</h3>` : '';
    if (description) {
      this.info.innerHTML += `<p>${description}</p>`;
    }
    
    // Update EXIF data
    if (photo.dataset.exif) {
      const exifData = JSON.parse(photo.dataset.exif);
      let exifHtml = '<div class="exif-data">';
      
      if (exifData.camera) exifHtml += `<span>${exifData.camera}</span>`;
      if (exifData.lens) exifHtml += `<span>${exifData.lens}</span>`;
      if (exifData.settings) {
        exifHtml += `<span>${exifData.settings.join(' • ')}</span>`;
      }
      if (exifData.datetime) exifHtml += `<span>${exifData.datetime}</span>`;
      
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