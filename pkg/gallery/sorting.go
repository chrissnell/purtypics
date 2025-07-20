package gallery

import (
	"sort"
	"time"
)

// GetNewestPhotoTime returns the most recent photo time in the album
func (a *Album) GetNewestPhotoTime() time.Time {
	var newest time.Time
	for _, photo := range a.Photos {
		if photo.EXIF != nil && !photo.EXIF.DateTime.IsZero() {
			if newest.IsZero() || photo.EXIF.DateTime.After(newest) {
				newest = photo.EXIF.DateTime
			}
		}
	}
	// If no EXIF dates found, use current time
	if newest.IsZero() {
		newest = time.Now()
	}
	return newest
}

// SortPhotosByDate sorts photos by EXIF date (oldest to newest)
func (a *Album) SortPhotosByDate() {
	sort.Slice(a.Photos, func(i, j int) bool {
		timeI := time.Now()
		timeJ := time.Now()
		
		if a.Photos[i].EXIF != nil && !a.Photos[i].EXIF.DateTime.IsZero() {
			timeI = a.Photos[i].EXIF.DateTime
		}
		if a.Photos[j].EXIF != nil && !a.Photos[j].EXIF.DateTime.IsZero() {
			timeJ = a.Photos[j].EXIF.DateTime
		}
		
		return timeI.Before(timeJ)
	})
}

// SortAlbumsByDate sorts albums by their newest photo (newest album first)
func SortAlbumsByDate(albums []Album) {
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].GetNewestPhotoTime().After(albums[j].GetNewestPhotoTime())
	})
}