# Convrt - Fast Image Converter

A lightning-fast command-line tool for converting and optimizing images, written in Go.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Make (optional, but recommended)
- Git

### Quick Start (After Cloning)

1. Clone the repository:
```bash
git clone https://github.com/melker12345/convrt.git
cd convrt
```

2. Install the tool:

#### On Windows
```bash
# Option 1: Using PowerShell (Recommended)
Right-click install.ps1 → Run with PowerShell

# Option 2: Using Make
make install

# Option 3: Manual build
go mod tidy
go build -o convrt.exe
```

#### On Linux/macOS
```bash
# Option 1: Using install script (Recommended)
chmod +x install.sh
./install.sh

# Option 2: Using Make
make install

# Option 3: Manual build
go mod tidy
go build
```

3. Restart your terminal

4. Test the installation:
```bash
convrt --help
```

## Usage Examples

### Basic Conversion
```bash
# Convert JPEG to PNG
convrt image.jpg .png

# Convert PNG to JPEG with 85% quality
convrt image.png .jpg -q 85

# Convert to WebP
convrt image.jpg .webp
```

### Optimization
```bash
# Optimize a single image
convrt image.jpg -o

# Optimize all images in a directory
convrt images/ -o -r
```

### Batch Processing
```bash
# Convert all JPEGs to WebP
convrt images/*.jpg .webp

# Optimize all PNGs in a directory
convrt images/*.png -o
```

## Features
- Convert between multiple formats (JPEG, PNG, WebP, GIF, TIFF)
- Optimize images with quality control
- Batch processing with wildcard support
- Recursive directory processing
- Progress bar for large files
- Beautiful command-line interface
- Cross-platform support (Windows, Linux, macOS)

## Command Options
- `-o, --optimize` - Optimize the image
- `-q, --quality N` - Set quality (1-100, default: 90)
- `-w, --width N` - Resize to width N
- `-h, --height N` - Resize to height N
- `-r, --recursive` - Process directories recursively

## Troubleshooting

### Common Issues

1. Command not found
```bash
# On Windows, run:
refreshenv
# Or restart your terminal

# On Linux/macOS, run:
source ~/.bashrc  # for bash
source ~/.zshrc   # for zsh
```

2. Permission denied
```bash
# On Linux/macOS
chmod +x install.sh
chmod +x uninstall.sh
```

3. Build errors
```bash
# Clean and rebuild
go mod tidy
go mod verify
make clean
make build
```

## Uninstallation

### Windows
```bash
# Option 1: Using PowerShell
Right-click uninstall.ps1 → Run with PowerShell

# Option 2: Using Make
make uninstall
```

### Linux/macOS
```bash
# Option 1: Using script
./uninstall.sh

# Option 2: Using Make
make uninstall
```

## Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](LICENSE)
