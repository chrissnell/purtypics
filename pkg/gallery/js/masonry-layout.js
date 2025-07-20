// Masonry layout functionality
export function initMasonry() {
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
export function initAlbumGrid() {
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