### Generating image bundles

- First install the fyne command tools:
```cmd
go install fyne.io/fyne/v2/cmd/fyne@latest
```
- To bundle the whole folder:
```
fyne bundle -o images.go ./images
```

### Packaging

To package for MacOs:

```cmd
make create-package-macos
```

Images courtesy of Flaticon:

- Dwi Ridwanto
- Tanah Basah
- Smashicons