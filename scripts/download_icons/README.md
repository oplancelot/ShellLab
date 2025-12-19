# Download Item Icons to Local

This tool downloads all item icons locally for offline usage.

## Usage

```bash
cd tools/download_icons
go run main.go
```

## Features

- Read all items from data/items.json
- Extract unique icon names
- Download icons from Wowhead CDN
- Save to frontend/public/items/icons/
- Concurrent downloads (10 concurrent), intelligent rate limiting
- Skip existing files
- Display download progress

## Output

Icons save location: `frontend/public/items/icons/`

- File format: lowercase iconname.jpg
- Example: inv_sword_39.jpg

## Estimated Time

- ~5000 icons
- ~10-20KB each
- Total size: ~50-100MB
- Download time: ~10-15 minutes

## Next Steps

Rebuild the application after downloading:

```bash
wails build
```

The application will prioritize local icons and display normally even when network connection fails.
