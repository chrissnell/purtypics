package gallery

import (
	"os"
	"sort"
	"time"

	"github.com/cjs/purtypics/pkg/exif"
)

// GetOldestPhotoTime returns the earliest photo time in the album,
// representing when the album's photos were originally taken.
func (a *Album) GetOldestPhotoTime() time.Time {
	var oldest time.Time
	for _, photo := range a.Photos {
		if photo.EXIF != nil && !photo.EXIF.DateTime.IsZero() {
			if oldest.IsZero() || photo.EXIF.DateTime.Before(oldest) {
				oldest = photo.EXIF.DateTime
			}
		}
	}
	if !oldest.IsZero() {
		return oldest
	}
	// Fall back to file modification time of first photo
	if len(a.Photos) > 0 {
		if info, err := os.Stat(a.Photos[0].Path); err == nil {
			return info.ModTime()
		}
	}
	return time.Time{}
}

// SetCreatedAtFromPhotos sets CreatedAt to the oldest EXIF date in the album
func (a *Album) SetCreatedAtFromPhotos() {
	a.CreatedAt = a.GetOldestPhotoTime()
}

// SortPhotosByDate sorts photos by EXIF date (oldest to newest)
func (a *Album) SortPhotosByDate() {
	sort.Slice(a.Photos, func(i, j int) bool {
		timeI := photoTime(a.Photos[i])
		timeJ := photoTime(a.Photos[j])
		return timeI.Before(timeJ)
	})
}

func photoTime(p Photo) time.Time {
	if p.EXIF != nil && !p.EXIF.DateTime.IsZero() {
		return p.EXIF.DateTime
	}
	// Fall back to file modification time
	if info, err := os.Stat(p.Path); err == nil {
		return info.ModTime()
	}
	return time.Time{}
}

// SortAlbumsByDate sorts albums by their original photo dates (newest album first)
func SortAlbumsByDate(albums []Album) {
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].CreatedAt.After(albums[j].CreatedAt)
	})
}

// SetAlbumDatesFromFirstPhoto extracts EXIF from each album's first photo
// to set CreatedAt without full processing. Used by the editor for sorting.
func SetAlbumDatesFromFirstPhoto(albums []Album) {
	for i := range albums {
		if len(albums[i].Photos) == 0 {
			continue
		}
		photo := &albums[i].Photos[0]
		if data, err := exif.ExtractMetadata(photo.Path); err == nil && !data.DateTime.IsZero() {
			albums[i].CreatedAt = data.DateTime
		} else if info, err := os.Stat(photo.Path); err == nil {
			albums[i].CreatedAt = info.ModTime()
		}
	}
}

// SortAlbumsByCustomOrder reorders albums to match the given order of album IDs.
// Albums in the order list come first (in that order), followed by any remaining
// albums in their current order.
func SortAlbumsByCustomOrder(albums []Album, order []string) {
	pos := make(map[string]int, len(order))
	for i, id := range order {
		pos[id] = i
	}

	sort.SliceStable(albums, func(i, j int) bool {
		pi, okI := pos[albums[i].ID]
		pj, okJ := pos[albums[j].ID]
		if okI && okJ {
			return pi < pj
		}
		if okI {
			return true
		}
		if okJ {
			return false
		}
		return false // preserve existing order for unlisted albums
	})
}
