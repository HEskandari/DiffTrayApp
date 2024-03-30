### Generating image bundles

- First install the fynd command tools:
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
fyne package -os darwin -icon icon.png
```

Images courtesy of Flaticon:

- Dwi Ridwanto
- Tanah Basah
- Smashicons