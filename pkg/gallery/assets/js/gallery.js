// Modern Gallery JavaScript with video support

document.addEventListener('DOMContentLoaded', function() {
    // Initialize video hover functionality
    initVideoHover();
    
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
            initLayout: false // We'll layout after images load
        });
        
        // Layout Masonry after all images have loaded
        imagesLoaded(gridElem, function() {
            msnry.layout();
            
            // Add loaded class to items for fade-in effect
            const items = gridElem.querySelectorAll('.grid-item');
            items.forEach(function(item, index) {
                setTimeout(function() {
                    item.classList.add('loaded');
                }, index * 50); // Stagger the fade-in
            });
        });
        
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
    
    // Initialize lightbox functionality
    const photoLinks = document.querySelectorAll('.photo-link[data-lightbox]');
    if (photoLinks.length > 0) {
        initializeLightbox();
    }
});

// Video hover functionality
function initVideoHover() {
    const videos = document.querySelectorAll('.video-item');
    
    videos.forEach(item => {
        const videoContainer = item.querySelector('.video-container');
        if (!videoContainer) return;
        
        const poster = videoContainer.querySelector('.video-poster');
        const video = videoContainer.querySelector('.video-preview');
        const playButton = videoContainer.querySelector('.play-button');
        
        if (!video || !poster) return;
        
        // Preload on first interaction
        let hasPreloaded = false;
        
        videoContainer.addEventListener('mouseenter', () => {
            // Preload the video on first hover
            if (!hasPreloaded) {
                video.preload = 'metadata';
                hasPreloaded = true;
            }
            
            // Show video and hide play button
            video.style.display = 'block';
            setTimeout(() => {
                video.style.opacity = '1';
                poster.style.opacity = '0';
            }, 10);
            
            // Ensure video is muted (required for autoplay)
            video.muted = true;
            
            // Try to play the video
            const playPromise = video.play();
            
            if (playPromise !== undefined) {
                playPromise.then(() => {
                    // Autoplay started successfully
                    if (playButton) {
                        playButton.style.opacity = '0';
                    }
                }).catch(error => {
                    console.error('Video autoplay failed:', error);
                    // Revert to poster on autoplay failure
                    video.style.opacity = '0';
                    poster.style.opacity = '1';
                    
                    // Try once more after user interaction
                    if (error.name === 'NotAllowedError') {
                        // Add click handler to play on click
                        videoContainer.addEventListener('click', function playOnClick() {
                            video.play().catch(e => console.error('Video play on click failed:', e));
                            videoContainer.removeEventListener('click', playOnClick);
                        }, { once: true });
                    }
                });
            }
        });
        
        videoContainer.addEventListener('mouseleave', () => {
            // Hide video and show poster
            video.style.opacity = '0';
            poster.style.opacity = '1';
            video.pause();
            video.currentTime = 0;
            
            // Reset play button after transition
            if (playButton) {
                setTimeout(() => {
                    playButton.style.opacity = '1';
                }, 300);
            }
        });
    });
}

// Enhanced lightbox implementation with video support
function initializeLightbox() {
    // Create lightbox elements
    const lightbox = document.createElement('div');
    lightbox.className = 'lightbox';
    lightbox.innerHTML = `
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
    document.body.appendChild(lightbox);
    
    const lightboxImg = lightbox.querySelector('img');
    const lightboxVideo = lightbox.querySelector('video');
    const lightboxInfo = lightbox.querySelector('.lightbox-info');
    const closeBtn = lightbox.querySelector('.lightbox-close');
    const prevBtn = lightbox.querySelector('.lightbox-prev');
    const nextBtn = lightbox.querySelector('.lightbox-next');
    
    let currentIndex = 0;
    const photoCards = Array.from(document.querySelectorAll('.photo-card'));
    
    // Open lightbox on photo click
    photoCards.forEach((card, index) => {
        const link = card.querySelector('.photo-link');
        if (link) {
            link.addEventListener('click', function(e) {
                e.preventDefault();
                currentIndex = index;
                showPhoto(currentIndex);
                lightbox.style.display = 'flex';
                document.body.style.overflow = 'hidden';
                updateNavigation();
            });
        }
    });
    
    // Close lightbox
    closeBtn.addEventListener('click', closeLightbox);
    lightbox.addEventListener('click', function(e) {
        if (e.target === lightbox || e.target.classList.contains('lightbox-content')) {
            closeLightbox();
        }
    });
    
    // Navigation
    prevBtn.addEventListener('click', function() {
        navigate(-1);
    });
    
    nextBtn.addEventListener('click', function() {
        navigate(1);
    });
    
    // Keyboard navigation
    document.addEventListener('keydown', function(e) {
        if (lightbox.style.display !== 'flex') return;
        
        if (e.key === 'Escape') closeLightbox();
        if (e.key === 'ArrowLeft') navigate(-1);
        if (e.key === 'ArrowRight') navigate(1);
    });
    
    function showPhoto(index) {
        const card = photoCards[index];
        const link = card.querySelector('.photo-link');
        const isVideo = card.classList.contains('video-item');
        const fullSrc = link ? link.getAttribute('href') : '';
        
        // Extract title
        const photoImg = card.querySelector('img');
        const photoTitle = card.querySelector('.photo-title');
        const title = photoTitle ? photoTitle.textContent : (photoImg ? photoImg.getAttribute('alt') : '');
        
        // Show/hide media elements
        if (isVideo) {
            lightboxImg.style.display = 'none';
            lightboxVideo.style.display = 'block';
            lightboxVideo.src = card.dataset.videoSrc || fullSrc;
        } else {
            lightboxVideo.style.display = 'none';
            lightboxVideo.pause();
            lightboxVideo.src = '';
            lightboxImg.style.display = 'block';
            lightboxImg.src = fullSrc;
        }
        
        // Update info
        lightboxInfo.innerHTML = title ? `<h3>${title}</h3>` : '';
    }
    
    function navigate(direction) {
        currentIndex = (currentIndex + direction + photoCards.length) % photoCards.length;
        showPhoto(currentIndex);
        updateNavigation();
    }
    
    function updateNavigation() {
        prevBtn.style.display = photoCards.length > 1 ? 'block' : 'none';
        nextBtn.style.display = photoCards.length > 1 ? 'block' : 'none';
    }
    
    function closeLightbox() {
        lightbox.style.display = 'none';
        document.body.style.overflow = '';
        lightboxVideo.pause();
        lightboxVideo.src = '';
    }
}