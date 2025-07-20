// Photo map functionality
export class PhotoMap {
  constructor(options = {}) {
    this.mapId = options.mapId || 'map';
    this.photos = options.photos || [];
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
      const fullSrc = photo.dataset.fullSrc;
      
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
            <a href="${fullSrc}" target="_blank">View full size</a>
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

  updateView(lat, lng, zoom = 15) {
    if (this.map) {
      this.map.setView([lat, lng], zoom);
    }
  }

  destroy() {
    if (this.map) {
      this.map.remove();
      this.map = null;
      this.markers = [];
    }
  }
}