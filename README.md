# Pixel Powaahh - Image Optimizer and Converter

## Introduction

Pixel Powaahh is a simple command-line tool for optimizing and converting image files. The idea behind this project is driven by a mix of impatience and curiosity. I recently found myself using a web service called TinyPNG to optimize a bunch of images, but I struggled with the limitations in uploading quotas. I realized I needed my own solution, a tool that could automate this process and make my life a little easier. And that's how Pixel Powaahh was born.

## Features

- Compress JPEG images with adjustable quality settings.
- Compress PNG images using a quantization algorithm.
- Converts images to the modern and efficient WebP format.

## Getting Started

### Installation

Clone this repository and navigate to the project directory.

### Usage

To optimize both PNG and JPEG images:

```bash
pp.exe
```

To optimize only JPEG images (with default quality 80):

```bash
pp.exe --jpg
```

To optimize only JPEG images with a specific quality (e.g., 90):

```bash
pp.exe --jpg --quality 90
```

To optimize PNG images using a quantization algorithm (a.k.a lossy compression):

```bash
pp.exe --png
```

To convert images to the WebP format with a specific quality (default is 80):

```bash
pp.exe --webp --quality 70
```

### Acknowledgments

Pixel Powaahh makes use of the following open-source libraries:

- [Go-ImageQuant](https://github.com/ultimate-guitar/go-imagequant) for PNG optimization using the quantization algorithm.
- [Chai2010's webp](https://github.com/chai2010/webp) for WebP image encoding.