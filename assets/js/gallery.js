// Modern Gallery JavaScript with proper Masonry implementation

document.addEventListener('DOMContentLoaded', function() {
    // Initialize Masonry for album or photo grids
    const grids = document.querySelectorAll('.masonry-grid');
    
    grids.forEach(function(gridElem) {
        // Initialize Masonry
        const msnry = new Masonry(gridElem, {
            itemSelector: '.grid-item',
            columnWidth: '.grid-sizer',
            percentPosition: true,
            gutter: 0, // Gutter is handled by CSS margins
            transitionDuration: '0.3s',
            initLayout: true // Layout immediately with visible images
        });
        
        // Store masonry instance on element for later access
        gridElem.masonry = msnry;
        
        // Progressive layout approach: layout visible images first, then as others load
        const visibleImages = gridElem.querySelectorAll('.grid-item img');
        let loadedCount = 0;
        const totalImages = visibleImages.length;
        
        // Function to handle individual image load
        const handleImageLoad = function() {
            loadedCount++;
            // Re-layout periodically as images load (batch updates)
            if (loadedCount % 5 === 0 || loadedCount === totalImages) {
                msnry.layout();
            }
        };
        
        // Set up load handlers for each image
        visibleImages.forEach(function(img, index) {
            if (img.complete && img.naturalWidth !== 0) {
                // Image already loaded (from cache)
                handleImageLoad();
                img.parentElement.parentElement.classList.add('loaded');
            } else {
                // Wait for image to load
                img.addEventListener('load', function() {
                    handleImageLoad();
                    // Add loaded class with stagger effect
                    setTimeout(function() {
                        img.parentElement.parentElement.classList.add('loaded');
                    }, index * 20); // Reduced stagger time
                });
                img.addEventListener('error', handleImageLoad); // Handle errors too
            }
        });
        
        // Initial layout with whatever is visible
        msnry.layout();
        
        // Re-layout on window resize
        let resizeTimer;
        window.addEventListener('resize', function() {
            clearTimeout(resizeTimer);
            resizeTimer = setTimeout(function() {
                msnry.layout();
            }, 250);
        });
    });
    
    // Lazy loading for images
    if ('loading' in HTMLImageElement.prototype) {
        // Browser supports native lazy loading
        const images = document.querySelectorAll('img[loading="lazy"]');
        images.forEach(img => {
            // Trigger layout when each image loads
            img.addEventListener('load', function() {
                const grid = img.closest('.masonry-grid');
                if (grid && grid.masonry) {
                    grid.masonry.layout();
                }
            });
        });
    } else {
        // Fallback for browsers that don't support native lazy loading
        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/lazysizes/5.3.2/lazysizes.min.js';
        document.body.appendChild(script);
    }
    
    // Simple lightbox functionality for photo pages
    const photoLinks = document.querySelectorAll('.photo-link[data-lightbox]');
    if (photoLinks.length > 0) {
        initializeLightbox(photoLinks);
    }
});

// Basic lightbox implementation
function initializeLightbox(links) {
    // Create lightbox elements
    const lightbox = document.createElement('div');
    lightbox.className = 'lightbox';
    lightbox.innerHTML = `
        <div class="lightbox-content">
            <img class="lightbox-image" src="" alt="">
            <button class="lightbox-close">&times;</button>
            <button class="lightbox-prev">‹</button>
            <button class="lightbox-next">›</button>
        </div>
    `;
    document.body.appendChild(lightbox);
    
    const lightboxImage = lightbox.querySelector('.lightbox-image');
    const closeBtn = lightbox.querySelector('.lightbox-close');
    const prevBtn = lightbox.querySelector('.lightbox-prev');
    const nextBtn = lightbox.querySelector('.lightbox-next');
    
    let currentIndex = 0;
    const photos = Array.from(links);
    
    // Open lightbox
    links.forEach((link, index) => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            currentIndex = index;
            showPhoto(currentIndex);
            lightbox.classList.add('active');
        });
    });
    
    // Close lightbox
    closeBtn.addEventListener('click', closeLightbox);
    lightbox.addEventListener('click', function(e) {
        if (e.target === lightbox) {
            closeLightbox();
        }
    });
    
    // Navigation
    prevBtn.addEventListener('click', function() {
        currentIndex = (currentIndex - 1 + photos.length) % photos.length;
        showPhoto(currentIndex);
    });
    
    nextBtn.addEventListener('click', function() {
        currentIndex = (currentIndex + 1) % photos.length;
        showPhoto(currentIndex);
    });
    
    // Keyboard navigation
    document.addEventListener('keydown', function(e) {
        if (!lightbox.classList.contains('active')) return;
        
        if (e.key === 'Escape') closeLightbox();
        if (e.key === 'ArrowLeft') prevBtn.click();
        if (e.key === 'ArrowRight') nextBtn.click();
    });
    
    function showPhoto(index) {
        const link = photos[index];
        lightboxImage.src = link.href;
        lightboxImage.alt = link.querySelector('img').alt;
    }
    
    function closeLightbox() {
        lightbox.classList.remove('active');
    }
}

// Add lightbox styles
const lightboxStyles = `
.lightbox {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.9);
    z-index: 1000;
    cursor: pointer;
}

.lightbox.active {
    display: flex;
    align-items: center;
    justify-content: center;
}

.lightbox-content {
    position: relative;
    max-width: 90%;
    max-height: 90%;
    cursor: default;
}

.lightbox-image {
    max-width: 100%;
    max-height: 90vh;
    object-fit: contain;
}

.lightbox-close,
.lightbox-prev,
.lightbox-next {
    position: absolute;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    padding: 10px 15px;
    transition: background 0.3s;
}

.lightbox-close:hover,
.lightbox-prev:hover,
.lightbox-next:hover {
    background: rgba(255, 255, 255, 0.2);
}

.lightbox-close {
    top: 20px;
    right: 20px;
}

.lightbox-prev {
    left: 20px;
    top: 50%;
    transform: translateY(-50%);
}

.lightbox-next {
    right: 20px;
    top: 50%;
    transform: translateY(-50%);
}

@media (max-width: 768px) {
    .lightbox-prev,
    .lightbox-next {
        padding: 5px 10px;
        font-size: 1.5rem;
    }
}
`;

// Inject lightbox styles
const styleSheet = document.createElement('style');
styleSheet.textContent = lightboxStyles;
document.head.appendChild(styleSheet);