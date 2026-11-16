# 6. Diyaloglar, mesajlar ve dosyalar

## Diyaloglar cevaplarını döndürür

OzdenUI'da bir diyalog, terminaldeki bir soru gibi, cevabı döndüren bir fonksiyon çağrısıdır:

```go
if ui.Confirm("Dosya silinecek.", ui.DialogOptions{Title: "Silinsin mi?", ConfirmText: "Sil", Danger: true}) {
	deleteTheFile()
}

name, ok := ui.Prompt("Klasör için bir ad:", ui.DialogOptions{Value: "Fotoğraflar", Required: true}) // ok false: vazgeçildi
choice := ui.Ask("Değişiklikler kaydedilsin mi?", "Kaydet", "Kaydetme", "Vazgeç")                    // tıklanan düğmenin metni
choice = ui.AskWithTitle("Kaydedilmemiş değişiklikler", "Kaydedilsin mi?", "Kaydet", "Kaydetme")
ui.Alert("Kopyalama bitti.", ui.DialogOptions{Title: "Bitti"})
```

- **Kodunuz beklerken pencere çalışmaya devam eder**: öbür olaylar ve zamanlayıcılar çalışır, bir saat işlemeyi sürdürür.
  Cevabı yalnız soran fonksiyon bekler.
- `ui.DialogOptions` metinleri tutar: `Title`, `ConfirmText`, `CancelText`; `Prompt` için ayrıca `Value`, `Placeholder`,
  `FieldText` ve `Required` (boş bir cevap alınmaz). Yazmadığınız metin kitin metni olarak kalır.
- `Danger: true` onay düğmesini bir silme düğmesi olarak çizer. `ui.Ask()` ilk düğmeyi ana düğme yapar.
- Escape ya da dışarıya bir tıklama "hayır"dır: `Confirm` için `false`, `Prompt` için `ok == false`, `Ask` için `""`.
- Bir diyalog, burada olduğu gibi, bir olaydan ya da bir zamanlayıcıdan sorulur. Sizin kendi bir goroutine'iniz
  pencereye yalnız `ui.Do()` ile ulaşır (4. bölüm).

`examples/dialogs` her biri için bir düğme taşır ve cevabı Go'nun yazdığı gibi gösterir.

## Kısa mesajlar: toast'lar

```go
ui.ToastSuccess("Kaydedildi.")
ui.ToastError("Dosya kaydedilemedi.", "Hata") // mesajın ardından bir başlık gelebilir
ui.ToastWarning("Disk dolmak üzere.")
ui.ToastInfo("3 yeni mesaj")
```

Bir toast bir köşede belirir ve kendiliğinden gider. Hiç beklemez: kod hemen devam eder.

## Balonlar: ShowHint ve Tooltip

```go
save := ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet"})
ui.Tooltip(ui.TooltipProps{Target: save, HintText: "Dosyayı kaydeder (Ctrl+S)"}) // fare üstündeyken, kendiliğinden

ui.ShowHint(save, "Kaydedildi.", 2000)              // şimdi, düğmenin yanında, 2 sn ("top" onun yanıdır)
hint := ui.ShowHint(save, "Çalışıyor", 0, "bottom") // 0: hint.Remove() çağrılana kadar kalır
```

- `Trigger: "click"` alan bir `Tooltip` bir popover'dır: hedefine bir tıklama onu açar ve dışarıya bir tıklamaya kadar
  açık kalır. `HintPosition` onun yanıdır: `"top"`, `"bottom"`, `"left"`, `"right"`.
- v0.1'de bir balon bir metin gösterir. Sizin nesnelerinizi tutan bir balon henüz yok.

## Sistemin dosya pencereleri

```go
path := ui.OpenFile(ui.FileOptions{Types: []string{"png", "jpg"}}) // bir yol ya da "" (vazgeçildi)
paths := ui.OpenFiles()                                            // bir []string, vazgeçilince boş
target := ui.SaveFile(ui.FileOptions{FileName: "rapor.csv", Types: []string{"csv"}})
folder := ui.ChooseFolder(ui.FileOptions{Title: "Fotoğraflar nerede?"})
```

- Bunlar bilgisayarın kendi pencereleridir (Linux'ta: `zenity`), bu yüzden yollar gerçektir: onları `os` ile okuyup
  yazın.
- `ui.FileOptions{Title, Directory, FileName, Types}`: `Directory` açıldığı klasör, `FileName` `SaveFile`'ın önerdiği
  ad, `Types` gösterdiği uzantılardır.
- `SaveFile` yalnız yolu verir: dosyayı programınız yazar.

### Küçük bir düzenleyici

```go
package main

import (
	"os"
	"path/filepath"

	ui "github.com/bug7a/ozdenui-go"
)

const pageColor = "#F4F5F7"

var fileTypes = []string{"txt", "md"}

var (
	path string // notun dosyası; açılana ya da kaydedilene kadar ""
	note *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Not", Width: 560, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", pageColor)
	ui.VGroup(ui.Props{Gap: 12, Padding: 20, Align: "left top"}, func() {
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Aç", OnClick: openNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Tone: "cyan", OnClick: saveNote})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Temizle", Tone: "orange", OnClick: clearNote})
		})
		note = ui.TextAreaField(ui.TextAreaFieldProps{TitleText: "Not", Width: "100%", Height: 260})
	})
}

func openNote() {
	chosen := ui.OpenFile(ui.FileOptions{Types: fileTypes})
	if chosen == "" {
		return
	}
	content, err := os.ReadFile(chosen)
	if err != nil {
		ui.ToastError(err.Error(), "Açılamadı")
		return
	}
	path = chosen
	note.SetText(string(content))
	note.Set("titleText", filepath.Base(path))
}

func saveNote() {
	if path == "" {
		path = ui.SaveFile(ui.FileOptions{FileName: "not.txt", Types: fileTypes})
	}
	if path == "" {
		return
	}
	if err := os.WriteFile(path, []byte(note.Text()), 0o644); err != nil {
		ui.ToastError(err.Error(), "Kaydedilemedi")
		return
	}
	note.Set("titleText", filepath.Base(path))
	ui.ToastSuccess("Kaydedildi.", filepath.Base(path))
}

func clearNote() {
	if ui.Confirm("Penceredeki metin silinecek.", ui.DialogOptions{Title: "Not temizlensin mi?", ConfirmText: "Temizle", Danger: true}) {
		note.SetText("")
	}
}
```

`note.Text()` sayfaya hiçbir şey sormaz: alan, kullanıcı yazarken metnini programa söyler.
`examples/csv-viewer`, `examples/image-viewer` ve `examples/markdown-editor` dosyalarını aynı yoldan açar.

## Pencerenin üstünde bir bekleme: LoadingScreen

```go
loading := ui.LoadingScreen(ui.LoadingScreenProps{TitleText: "Kopyalanıyor", MessageText: "Başlıyor...", Progress: 0})
go func() {
	for step := 1; step <= 4; step++ {
		copyOneFile(step) // yavaş iş, goroutine'de
		ui.Do(func() {
			loading.Call("setProgress", float64(step)/4) // 0 .. 1
			loading.Call("setMessageText", fmt.Sprintf("Dosya %d / 4", step))
		})
	}
	ui.Do(func() { loading.Call("close") })
}()
```

Nerede yapılırsa yapılsın bütün pencereyi örter. `Progress: -1` (varsayılanı) çubuk göstermez; `MinDuration` (ms) onu
bir an tutar, böylece hızlı bir iş göz kırpmaz; `RemoveOnClose: true` onu `close`'dan sonra siler, bir kez kullanılan
bir ekran için.

## Dosya bırakılacak bir yer: FilePicker

`FilePicker`, dosya bırakılan ya da tıklayınca dosya penceresini açan bir yerdir. `picker.Get("files")` onun listesidir,
map'lerden oluşan bir `[]any`: `name`, `path`, `size`, `sizeText`, `type`, `extension`, `lastModified`.

Üstüne bırakılan bir dosya sayfadan gelir ve sayfa bir dosyanın nerede olduğunu hiç bilmez: `path` yalnız onun adıdır.
Tıklamaya `OnOpenDialog` ile sistemin dosya penceresini verin; dosyaları gerçek yollarıyla gelir:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	list *ui.Object
	rows []*ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Dosyalar", Width: 520, Height: 460})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14, Padding: 24, Align: "left top"}, func() {
		ui.FilePicker(ui.FilePickerProps{Width: "100%", ZoneHeight: 140, Multiple: true, ShowList: false,
			TitleText: "Dosyaları buraya bırakın", OnOpenDialog: chooseFiles, OnChange: showFiles})
		list = ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 6, Align: "left top"}, func() {
			rows = append(rows, ui.Label(ui.Props{Text: "Henüz dosya yok.", FontSize: 14, TextColor: "#5A6270"}))
		})
	})
}

// chooseFiles sayfanınkinin yerine sistemin dosya penceresini açar.
func chooseFiles(picker *ui.Object) {
	files := []ui.P{}
	for _, path := range ui.OpenFiles() {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		files = append(files, ui.P{"name": filepath.Base(path), "path": path, "size": info.Size(),
			"lastModified": info.ModTime().UnixMilli()})
	}
	picker.Call("addFiles", files)
}

func showFiles(picker *ui.Object) {
	for _, row := range rows {
		row.Remove()
	}
	rows = nil
	files, _ := picker.Get("files").([]any)
	list.Add(func() {
		for _, item := range files {
			file, _ := item.(map[string]any)
			text := fmt.Sprintf("%v (%v)", file["path"], file["sizeText"])
			rows = append(rows, ui.Label(ui.Props{Text: text, FontSize: 14}))
		}
	})
}
```

- `Accept: ".csv,.txt"` ya da `"image/*"`; `MaxSize` (bayt) ve `MaxFiles`: geçemeyenler `OnError`'a gider.
  `Multiple: true` listeye ekler; bu olmadan yeni dosya eskisinin yerini alır.
- `picker.Call("clear")` listeyi boşaltır; `TitleText`, `DescriptionText`: metinleri; `IconFile`: sizin bir resminiz.
- Bırakılan bir dosyanın baytları Go'ya ulaşmaz. Bir dosyanın içeriği için, yukarıdaki gibi sistemin penceresini
  kullanın.

## Sonraki

[7. bölüm: Görünüm](07-looks.md): koyu pencereler, minimal görünüm ve sizin renkleriniz.
