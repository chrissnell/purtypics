// Video hover functionality
export function initVideoHover() {
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