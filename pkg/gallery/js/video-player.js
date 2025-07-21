// Video hover functionality
export function initVideoHover() {
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
      video.style.opacity = '1';
      poster.style.opacity = '0';
      video.play().catch(e => {
        console.error('Video play failed:', e);
      });
      
      if (playButton) {
        playButton.style.opacity = '0';
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