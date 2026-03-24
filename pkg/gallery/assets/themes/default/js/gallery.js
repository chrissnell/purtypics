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
        const visibleImages = gridElem.querySelectorAll('.grid-item img, .album-cover img');
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
                img.closest('.grid-item, .album-card')?.classList.add('loaded');
            } else {
                // Wait for image to load
                img.addEventListener('load', function() {
                    handleImageLoad();
                    // Add loaded class with stagger effect
                    setTimeout(function() {
                        img.closest('.grid-item, .album-card')?.classList.add('loaded');
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
    
    // Video hover preview
    initVideoHover();

    // Simple lightbox functionality for photo pages
    const photoLinks = document.querySelectorAll('.photo-link[data-lightbox]');
    if (photoLinks.length > 0) {
        initializeLightbox(photoLinks);
    }
});

function initVideoHover() {
    const videoItems = document.querySelectorAll('.video-item');
    videoItems.forEach(item => {
        const preview = item.querySelector('.video-preview');
        if (!preview) return;

        item.addEventListener('mouseenter', () => {
            preview.play().catch(() => {});
        });

        item.addEventListener('mouseleave', () => {
            preview.pause();
            preview.currentTime = 0;
        });
    });
}

// Basic lightbox implementation
function initializeLightbox(links) {
    // Create lightbox elements
    const lightbox = document.createElement('div');
    lightbox.className = 'lightbox';
    lightbox.innerHTML = `
        <div class="lightbox-content">
            <img class="lightbox-image" src="" alt="">
            <video class="lightbox-video" controls preload="metadata"></video>
            <div class="lightbox-info"></div>
            <div class="lightbox-exif"></div>
            <button class="lightbox-close">&times;</button>
            <button class="lightbox-prev">&#10094;</button>
            <button class="lightbox-next">&#10095;</button>
        </div>
    `;
    document.body.appendChild(lightbox);

    const lightboxImage = lightbox.querySelector('.lightbox-image');
    const lightboxVideo = lightbox.querySelector('.lightbox-video');
    const lightboxExif = lightbox.querySelector('.lightbox-exif');
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
        const card = link.closest('.photo-card');
        const isVideo = card && card.dataset.video === 'true';

        // Pause any hover preview
        const preview = card && card.querySelector('.video-preview');
        if (preview) {
            preview.pause();
        }

        if (isVideo) {
            const videoSrc = card.dataset.videoSrc;
            lightboxVideo.src = videoSrc;
            lightboxVideo.style.display = 'block';
            lightboxImage.style.display = 'none';
        } else {
            lightboxVideo.pause();
            lightboxVideo.src = '';
            lightboxVideo.style.display = 'none';
            lightboxImage.src = link.href;
            lightboxImage.alt = link.querySelector('img').alt;
            lightboxImage.style.display = 'block';
        }

        // Show EXIF data
        if (card) {
            const parts = [];
            if (card.dataset.camera) parts.push(card.dataset.camera);
            if (card.dataset.lens) parts.push(card.dataset.lens);

            const settings = [];
            if (card.dataset.focal) settings.push(card.dataset.focal + 'mm');
            if (card.dataset.aperture) settings.push('f/' + card.dataset.aperture);
            if (card.dataset.shutter) settings.push(card.dataset.shutter + 's');
            if (card.dataset.iso) settings.push('ISO ' + card.dataset.iso);
            if (settings.length) parts.push(settings.join(' · '));

            if (card.dataset.datetime) parts.push(card.dataset.datetime);

            if (parts.length) {
                lightboxExif.innerHTML = parts.map(p => '<span>' + p + '</span>').join('');
                lightboxExif.style.display = 'block';
            } else {
                lightboxExif.style.display = 'none';
            }
        } else {
            lightboxExif.style.display = 'none';
        }
    }

    function closeLightbox() {
        lightbox.classList.remove('active');
        lightboxVideo.pause();
        lightboxVideo.src = '';
        lightboxVideo.style.display = 'none';
        lightboxExif.style.display = 'none';
        lightboxExif.innerHTML = '';
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
    max-width: 80%;
    max-height: 90%;
    cursor: default;
}

.lightbox-image {
    max-width: 100%;
    max-height: 90vh;
    object-fit: contain;
}

.lightbox-close {
    position: absolute;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    padding: 10px 15px;
    transition: background 0.3s;
    top: 20px;
    right: 20px;
}

.lightbox-close:hover {
    background: rgba(255, 255, 255, 0.2);
}

.lightbox-prev,
.lightbox-next {
    position: fixed;
    background: rgba(255, 255, 255, 0.15);
    color: white;
    border: none;
    font-size: 2.5rem;
    font-weight: 700;
    cursor: pointer;
    padding: 1.5rem 0.75rem;
    transition: background 0.2s, opacity 0.2s;
    opacity: 0.85;
    border-radius: 4px;
    line-height: 1;
    z-index: 1001;
    top: 50%;
    transform: translateY(-50%);
}

.lightbox-prev:hover,
.lightbox-next:hover {
    background: rgba(255, 255, 255, 0.3);
    opacity: 1;
}

.lightbox-prev {
    left: 0.75rem;
}

.lightbox-next {
    right: 0.75rem;
}

.lightbox-info {
    text-align: center;
    color: rgba(255, 255, 255, 0.9);
    padding: 0.5rem 0 0;
}

.lightbox-exif {
    display: none;
    text-align: center;
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.8rem;
    padding: 0.25rem 0;
}

.lightbox-exif span {
    display: inline-block;
}

.lightbox-exif span + span::before {
    content: ' · ';
    margin: 0 0.25rem;
}

@media (max-width: 768px) {
    .lightbox-prev,
    .lightbox-next {
        padding: 5px 10px;
        font-size: 1.5rem;
    }

    .lightbox-exif {
        font-size: 0.7rem;
    }

    .lightbox-exif span {
        display: block;
    }

    .lightbox-exif span + span::before {
        content: '';
        margin: 0;
    }
}
`;

// Inject lightbox styles
const styleSheet = document.createElement('style');
styleSheet.textContent = lightboxStyles;
document.head.appendChild(styleSheet);