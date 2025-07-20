# Purtypics Metadata System

Purtypics includes a powerful metadata system that allows you to customize your gallery without modifying the original photos.

## Metadata Editor

The easiest way to manage metadata is using the built-in web editor:

```bash
purtypics edit -source /path/to/photos
```

This will:
1. Start a web server on http://localhost:8080
2. Open your browser automatically
3. Provide a visual interface to edit gallery, album, and photo metadata
4. Auto-save changes to `gallery.yaml`

### Editor Features

- **Gallery Settings**: Edit overall gallery title, description, author, and copyright
- **Album Management**: 
  - Custom titles and descriptions
  - Select cover photos
  - Hide/show albums
  - Configure sort order
- **Photo Management**:
  - Custom titles and descriptions
  - Hide individual photos
  - Visual preview while editing

## Manual Metadata Editing

You can also edit the `gallery.yaml` file directly. Place it in your photos directory:

```yaml
# Gallery-level settings
title: "My Photo Gallery"
description: "A collection of memorable moments"
author: "Your Name"
copyright: "© 2024 Your Name"

# Album metadata
albums:
  "album-folder-name":
    title: "Custom Album Title"
    description: "Album description"
    cover_photo: "specific-photo.jpg"
    hidden: false
    sort_order: "date"  # Options: date, name, custom
    tags:
      - vacation
      - family

# Photo metadata  
photos:
  "album-folder/photo.jpg":
    title: "Custom Photo Title"
    description: "Photo description"
    hidden: false
    tags:
      - sunset
      - landscape
```

## Metadata Fields

### Gallery Metadata
- `title`: Gallery title (overrides CLI -title flag)
- `description`: Gallery description
- `author`: Gallery author/photographer
- `copyright`: Copyright notice

### Album Metadata
- `title`: Album display title
- `description`: Album description
- `cover_photo`: Filename of the photo to use as album cover
- `hidden`: Whether to hide this album (true/false)
- `sort_order`: How to sort photos ("date", "name", "custom")
- `custom_order`: Array of filenames when using custom sort
- `tags`: Array of tags for categorization

### Photo Metadata
- `title`: Photo display title
- `description`: Photo description  
- `hidden`: Whether to hide this photo (true/false)
- `tags`: Array of tags for categorization
- `sort_index`: Number for custom ordering

## Usage Examples

### Hide Sensitive Photos
```yaml
photos:
  "family-album/private-photo.jpg":
    hidden: true
```

### Custom Album Cover
```yaml
albums:
  "vacation-2024":
    title: "Summer Vacation 2024"
    cover_photo: "sunset-beach.jpg"
```

### Custom Photo Order
```yaml
albums:
  "wedding":
    sort_order: "custom"
    custom_order:
      - "ceremony-start.jpg"
      - "vows.jpg"
      - "first-kiss.jpg"
      - "reception.jpg"
```

## Workflow

1. Organize photos into album folders
2. Run `purtypics edit` to launch the metadata editor
3. Customize titles, descriptions, and settings
4. Save changes (auto-saves in the editor)
5. Run `purtypics generate` to build your gallery with metadata

The metadata system ensures your customizations are preserved separately from your photos, making it easy to regenerate the gallery or version control your settings.