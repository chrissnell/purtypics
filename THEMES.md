# Purtypics Themes

Purtypics ships with 13 built-in themes. Set a theme in your `gallery.yaml`:

```yaml
theme: darkroom
```

Or select one from the **Gallery Settings** tab in the editor (`purtypics edit`).

---

## Aperture

Clean, modern layout with a light background and card-based album grid. Crisp typography with a bold accent divider.

| Gallery | Album |
|---------|-------|
| [![Aperture gallery](.assets/themes/thumbs/aperture.png)](.assets/themes/aperture.png) | [![Aperture album](.assets/themes/thumbs/aperture-album.png)](.assets/themes/aperture-album.png) |

---

## Atelier

Warm, artisanal feel with a creamy background and soft shadows. Photos presented like works in a studio.

| Gallery | Album |
|---------|-------|
| [![Atelier gallery](.assets/themes/thumbs/atelier.png)](.assets/themes/atelier.png) | [![Atelier album](.assets/themes/thumbs/atelier-album.png)](.assets/themes/atelier-album.png) |

---

## Blueprint

Technical blueprint aesthetic with a deep blue background and white grid lines. Photos displayed on a drafting-table canvas.

| Gallery | Album |
|---------|-------|
| [![Blueprint gallery](.assets/themes/thumbs/blueprint.png)](.assets/themes/blueprint.png) | [![Blueprint album](.assets/themes/thumbs/blueprint-album.png)](.assets/themes/blueprint-album.png) |

---

## Brutalist

High-contrast black and white with bold typography and zero decoration. Raw, uncompromising design.

| Gallery | Album |
|---------|-------|
| [![Brutalist gallery](.assets/themes/thumbs/brutalist.png)](.assets/themes/brutalist.png) | [![Brutalist album](.assets/themes/thumbs/brutalist-album.png)](.assets/themes/brutalist-album.png) |

---

## Coast

Breezy coastal palette with sandy tones and ocean-blue accents. Relaxed typography and airy spacing.

| Gallery | Album |
|---------|-------|
| [![Coast gallery](.assets/themes/thumbs/coast.png)](.assets/themes/coast.png) | [![Coast album](.assets/themes/thumbs/coast-album.png)](.assets/themes/coast-album.png) |

---

## Darkroom

Dark background with warm amber accents — photos glow against the darkness. Perfect for moody photography.

| Gallery | Album |
|---------|-------|
| [![Darkroom gallery](.assets/themes/thumbs/darkroom.png)](.assets/themes/darkroom.png) | [![Darkroom album](.assets/themes/thumbs/darkroom-album.png)](.assets/themes/darkroom-album.png) |

---

## Default

Soft pastel palette with dotted borders and a clean masonry layout. A balanced starting point for any gallery.

| Gallery | Album |
|---------|-------|
| [![Default gallery](.assets/themes/thumbs/default.png)](.assets/themes/default.png) | [![Default album](.assets/themes/thumbs/default-album.png)](.assets/themes/default-album.png) |

---

## Ember

Warm, fiery palette with deep reds and glowing orange highlights. Bold and dramatic.

| Gallery | Album |
|---------|-------|
| [![Ember gallery](.assets/themes/thumbs/ember.png)](.assets/themes/ember.png) | [![Ember album](.assets/themes/thumbs/ember-album.png)](.assets/themes/ember-album.png) |

---

## Kyoto

Japanese-inspired minimalism with muted earth tones, generous whitespace, and refined serif headings.

| Gallery | Album |
|---------|-------|
| [![Kyoto gallery](.assets/themes/thumbs/kyoto.png)](.assets/themes/kyoto.png) | [![Kyoto album](.assets/themes/thumbs/kyoto-album.png)](.assets/themes/kyoto-album.png) |

---

## Mono

Monochrome elegance — grayscale UI with clean lines and understated typography. Lets photos provide all the color.

| Gallery | Album |
|---------|-------|
| [![Mono gallery](.assets/themes/thumbs/mono.png)](.assets/themes/mono.png) | [![Mono album](.assets/themes/thumbs/mono-album.png)](.assets/themes/mono-album.png) |

---

## Nordic

Scandinavian calm with generous whitespace, cool blue-gray tones, and elegant serif headings.

| Gallery | Album |
|---------|-------|
| [![Nordic gallery](.assets/themes/thumbs/nordic.png)](.assets/themes/nordic.png) | [![Nordic album](.assets/themes/thumbs/nordic-album.png)](.assets/themes/nordic-album.png) |

---

## Polaroid

Nostalgic instant-photo cards with a handwriting font and slight rotation. Photos feel like scattered prints on a table.

| Gallery | Album |
|---------|-------|
| [![Polaroid gallery](.assets/themes/thumbs/polaroid.png)](.assets/themes/polaroid.png) | [![Polaroid album](.assets/themes/thumbs/polaroid-album.png)](.assets/themes/polaroid-album.png) |

---

## Salon

Gallery-wall presentation with a rich dark background and elegant framing. Photos hung salon-style for a curated, museum feel.

| Gallery | Album |
|---------|-------|
| [![Salon gallery](.assets/themes/thumbs/salon.png)](.assets/themes/salon.png) | [![Salon album](.assets/themes/thumbs/salon-album.png)](.assets/themes/salon-album.png) |

---

## Custom Themes

Create your own theme by adding a directory under `themes/` in your gallery source:

```
my-gallery/
├── themes/
│   └── mytheme/
│       ├── css/
│       │   └── gallery.css
│       └── templates/       # optional
│           ├── base.html
│           ├── index.html
│           └── album.html
└── gallery.yaml
```

A theme only needs to include files you want to override — anything missing falls back to the built-in default. See the [default theme](pkg/gallery/assets/themes/default/) as a reference.
