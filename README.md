![# Header](https://raw.githubusercontent.com/keshon/assets/main/pixelita-image-compressor/header.png)

# Pixelita Image Compressor

Pixelita is a Go-based image processing utility designed to facilitate encoding, decoding, and conversion of images into formats: JPEG, PNG, and WebP.

## Features

- **Format Conversion:** Easily convert images between different formats including JPEG, PNG, and WebP.
- **Quality Adjustment:** Tailor the quality of images during conversion to meet specific requirements.
- **Concurrent Processing:** Utilize parallel processing to optimize multiple images simultaneously.

![# Main Window Example](https://raw.githubusercontent.com/keshon/assets/main/pixelita-image-compressor/gui.jpg)

## Getting Started

### Prerequisites

- **Go (Golang):** Install Go on your system.
- **Go Modules:** Fetch required modules by running `go mod download`.

### Usage

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/keshon/pixelita-image-compressor.git
   ```

2. **Install Dependencies:**

   ```bash
   cd pixelita
   go mod download
   ```

4. **Run Pixelita:**

   Execute the following command to run Pixelita:

   ```bash
   go run cmd/main.go
   ```

    You can use supplied build scripts (`bat` for Windows and `sh` for Linux:
    - build-and-run - debug version (run after build).
    - build-release - release version with terminal window visible.
    - build-release-no-console - release version without terminal window.

### File Structure

- **config:** Config to hold various parameters.  
- **imageencode:** Contains functionalities for encoding images in JPEG, PNG, and WebP formats.
- **imagetype:** Manages image format types and their corresponding extensions.
- **pixelita:** Main utility for image processing and format conversion.
- **internal:** Includes additional utility functions and configurations used within the application.
- **filesystem:** File IO related helper object.

## Acknowledgment

The banner images used in this project were sourced from [Freepik](https://www.freepik.com), attributed to contributors [@Freepik](https://www.freepik.com/author/freepik) and [@starline](https://www.freepik.com/author/starline).


### Contribution

Contributions to Pixelita are welcome! You can create issues, submit pull requests, or suggest enhancements to improve the utility.

### License

Pixelita is licensed under the [MIT License](https://opensource.org/licenses/MIT).