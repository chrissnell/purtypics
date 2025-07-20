// Infinite scroll functionality
export class InfiniteScroll {
  constructor(options = {}) {
    this.container = options.container || '.masonry-grid';
    this.itemSelector = options.itemSelector || '.photo-item';
    this.loadMoreThreshold = options.threshold || 200;
    this.onLoadMore = options.onLoadMore || (() => {});
    this.loading = false;
    this.hasMore = true;
    this.page = 1;
  }

  init() {
    this.bindScrollEvent();
  }

  bindScrollEvent() {
    let ticking = false;
    
    const handleScroll = () => {
      if (!ticking) {
        window.requestAnimationFrame(() => {
          this.checkScroll();
          ticking = false;
        });
        ticking = true;
      }
    };
    
    window.addEventListener('scroll', handleScroll);
    window.addEventListener('resize', handleScroll);
  }

  checkScroll() {
    if (this.loading || !this.hasMore) return;
    
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const windowHeight = window.innerHeight;
    const documentHeight = document.documentElement.scrollHeight;
    
    if (scrollTop + windowHeight >= documentHeight - this.loadMoreThreshold) {
      this.loadMore();
    }
  }

  async loadMore() {
    this.loading = true;
    
    try {
      const newItems = await this.onLoadMore(this.page + 1);
      
      if (newItems && newItems.length > 0) {
        this.page++;
        this.appendItems(newItems);
      } else {
        this.hasMore = false;
      }
    } catch (error) {
      console.error('Error loading more items:', error);
    } finally {
      this.loading = false;
    }
  }

  appendItems(items) {
    const container = document.querySelector(this.container);
    if (!container) return;
    
    items.forEach(item => {
      container.appendChild(item);
    });
    
    // Trigger masonry layout update if available
    if (window.masonryInstance) {
      window.masonryInstance.appended(items);
      window.masonryInstance.layout();
    }
    
    // Re-initialize any necessary components for new items
    this.initNewItems(items);
  }

  initNewItems(items) {
    // Re-initialize lightbox for new items
    if (window.lightbox) {
      window.lightbox.photos = Array.from(document.querySelectorAll(this.itemSelector));
    }
    
    // Initialize lazy loading for new images
    const lazyImages = items.reduce((acc, item) => {
      const imgs = item.querySelectorAll('img[loading="lazy"]');
      return acc.concat(Array.from(imgs));
    }, []);
    
    if (window.imageObserver) {
      lazyImages.forEach(img => window.imageObserver.observe(img));
    }
  }

  reset() {
    this.page = 1;
    this.hasMore = true;
    this.loading = false;
  }

  disable() {
    this.hasMore = false;
  }
}