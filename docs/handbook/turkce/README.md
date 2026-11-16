# OzdenUI for Go — El Kitabı

Düz Go ile bir masaüstü penceresi kurun: nesneler, gruplar ve hazır bileşenler. Yalnız Go: JavaScript yok, HTML yok,
CSS yok, bir ön yüz için derleme adımı yok. Bu kitap, Go bilen ama hiç web sayfası yazmamış kişiler için konuyu adım
adım anlatır.


```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(start, ui.Window{Title: "Merhaba", Width: 360, Height: 200})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 12}, func() {
		ui.Label(ui.Props{Text: "Go'dan merhaba", FontSize: 24, Bold: true})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Kapat", OnClick: ui.Close})
	})
}
```

## Bölümler

| | Bölüm | Öğrendikleriniz |
|---|---|---|
| 1 | [İlk pencere](01-first-window.md) | kurulum, `ui.Run`, `start(page)`, `ui.Window`, pencere ile kodunuz nasıl konuşur |
| 2 | [Nesneler](02-objects.md) | `Box`, `Label`, `Button`, `Input`, `Icon`; `ui.Props`; verdiğiniz ve okuduğunuz değerler; silmek |
| 3 | [Yerleşim](03-layout.md) | `VGroup`, `HGroup`, `GridGroup`, `Fill: true`, `ScrollArea`, kendi yeri olan nesne |
| 4 | [Olaylar ve zaman](04-events-and-time.md) | tıklamalar ve tuşlar, `ui.Event`, `ui.Every` ve `ui.After`, bir goroutine ve `ui.Do` |
| 5 | [Bileşenler](05-components.md) | hazır parçalar: düğmeler, alanlar, seçimler, listeler, tablolar, sekmeler, çubuklar |
| 6 | [Diyaloglar, mesajlar ve dosyalar](06-dialogs-and-files.md) | `ui.Confirm()`, kısa mesajlar (toast), balonlar, sistemin dosya pencereleri |
| 7 | [Görünüm](07-looks.md) | koyu pencereler, sizin renginizde minimal görünüm, tek bir bileşenin renkleri |
| 8 | [Veri ve resimler](08-data-and-pictures.md) | grafikler, büyük tablolar, kayıtlar, Markdown, JSON, kod, resimler |
| 9 | [Kendi bileşeniniz](09-your-own-component.md) | nesnelerden yapılmış kendi türünüz, bir paket olarak paylaşmak |
| 10 | [Daha ileri](10-going-further.md) | Go'dan basic.js, hız, goroutine'ler, sınırlar, hata bulmak |

1'den 4'e kadar sırayla okuyun: bütün fikir onlarda. Sonra gereken bölümleri seçin.

## Nasıl okunur

- `package main` ile başlayan her program eksiksizdir: onu `go get github.com/bug7a/ozdenui-go` yapılmış bir modülün
  `main.go`'suna kopyalayın ve `go run .` ile çalıştırın. Modülün bir denetimi bunların hepsini derler ve çizer, bu
  yüzden doğru kalırlar.
- Modülün örnekleri (kaynağındaki `examples/` klasörü) daha uzun programlardır: bir yapılacaklar listesi, bir hesap
  makinesi, bir CSV görüntüleyici, bir düzenleyici, küçük bir oyun. Bölümler onları yeri geldikçe anar. Deponun bir
  kopyasından: `go run ./examples/todo`.
- `kod` biçimindeki bir ad, Go'daki bir addır. Penceredeki sayfa küçük bir JavaScript kütüphanesi olan
  [basic.js](https://bug7a.github.io/basic.js/) ile yapılır. Hiç JavaScript yazmanız gerekmez; 10. bölüm yine de ona
  nasıl ulaşacağınızı gösterir.
- Her fonksiyon ve her alan `go doc github.com/bug7a/ozdenui-go` içinde ve
  [pkg.go.dev](https://pkg.go.dev/github.com/bug7a/ozdenui-go)'dedir. Site: https://bug7a.github.io/ozdenui-go/

## Bir dakikada

- `ui.Run(start, ui.Window{...})` bir pencere açar. `start(page)` nesnelerini yapar; `page` pencerenin zeminidir.
- Bir nesne yazıldığı yerde yapılır: bir kutuya ya da gruba verilen fonksiyonun içindeyse o kutuya ya da gruba girer.
- Bir grup nesnelerini yerleştirir: `ui.VGroup` alt alta, `ui.HGroup` yan yana.
- Bir değer, nesne yapılırken bir alandır (`ui.Props{FontSize: 20}`), sonra bir ad (`label.Set("fontSize", 24)`).
- Bir olay bir fonksiyondur: `OnClick: save`. Hiçbir şey almayabilir, nesneyi ya da nesneyle olayı alabilir.
- Bileşenler durumlarını Go'da tutar: `box.Checked()`, `selectBox.Value()`, `field.Text()`. Bunu okumak sayfaya bir
  şey sormaz.
- Kodunuz tek bir goroutine'de çalışır: bir olay, bir zamanlayıcı ve `start()` hiçbir zaman aynı anda çalışmaz.
