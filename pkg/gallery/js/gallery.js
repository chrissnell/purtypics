// Main gallery module that coordinates all components
import { initVideoHover } from './video-player.js';
import { Lightbox } from './lightbox.js';
import { initLazyLoading } from './lazy-loader.js';
import { InfiniteScroll } from './infinite-scroll.js';
import { initMasonry, initAlbumGrid } from './masonry-layout.js';
import { PhotoMap } from './photo-map.js';

class Gallery {
  constructor() {
    this.lightbox = null;
    this.infiniteScroll = null;
    this.photoMap = null;
    this.masonryInstance = null;
  }

  init() {
    // Initialize video hover effects
    initVideoHover();
    
    // Initialize lightbox
    this.lightbox = new Lightbox();
    this.lightbox.init();
    window.lightbox = this.lightbox; // Make available globally for other modules
    
    // Initialize lazy loading
    initLazyLoading();
    
    // Initialize masonry layout
    this.masonryInstance = initMasonry() || initAlbumGrid();
    window.masonryInstance = this.masonryInstance;
    
    // Initialize infinite scroll if on a photo page
    if (document.querySelector('.masonry-grid')) {
      this.initInfiniteScroll();
    }
    
    // Initialize map if present
    if (document.getElementById('map')) {
      this.photoMap = new PhotoMap();
      this.photoMap.init();
    }
    
    // Handle dynamic content updates
    this.observeContentChanges();
  }

  initInfiniteScroll() {
    // Only initialize if we have pagination data
    const container = document.querySelector('.masonry-grid');
    if (!container || !container.dataset.totalPages) return;
    
    const totalPages = parseInt(container.dataset.totalPages, 10);
    const currentPage = parseInt(container.dataset.currentPage || '1', 10);
    
    if (totalPages <= 1) return;
    
    this.infiniteScroll = new InfiniteScroll({
      container: '.masonry-grid',
      itemSelector: '.photo-item',
      threshold: 400,
      onLoadMore: async (page) => {
        // This would be implemented to load more photos
        // For now, return empty array to indicate no more items
        return [];
      }
    });
    
    this.infiniteScroll.init();
  }

  observeContentChanges() {
    // Re-initialize components when new content is added dynamically
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.addedNodes.length > 0) {
          // Re-initialize video hover for new videos
          initVideoHover();
          
          // Re-layout masonry
          if (this.masonryInstance) {
            this.masonryInstance.layout();
          }
        }
      });
    });
    
    // Observe the main content area
    const mainContent = document.querySelector('main');
    if (mainContent) {
      observer.observe(mainContent, {
        childList: true,
        subtree: true
      });
    }
  }

  // Public method to programmatically open lightbox
  openLightbox(index) {
    if (this.lightbox) {
      this.lightbox.open(index);
    }
  }

  // Public method to refresh layout
  refreshLayout() {
    if (this.masonryInstance) {
      this.masonryInstance.layout();
    }
  }
}

// Initialize gallery when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    window.gallery = new Gallery();
    window.gallery.init();
  });
} else {
  window.gallery = new Gallery();
  window.gallery.init();
}

// Export for use in other scripts if needed
export default Gallery;