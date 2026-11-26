# 1. İlk pencere

## Kurulum

Programınızın klasöründe (bir Go modülü: `go mod init myapp` bir tane yapar):

```
go get github.com/bug7a/ozdenui-go
```

Go 1.22 ya da daha yenisi gerekir.

- **macOS:** pencere, programınızın içindeki sistemin web görünümüdür (WKWebView). cgo gerektirir, bu yüzden Xcode'un
  komut satırı araçları kurulu olmalıdır (`xcode-select --install`). İlk derleme yaklaşık on saniye sürer ve web
  görünümünün C kodundan uyarılar yazabilir: bunlar hata değildir, sonraki derlemeler hızlıdır.
- **Windows:** pencere, programınızın içindeki sistemin web görünümüdür (WebView2). C derleyicisi gerekmez. WebView2
  her Windows 11'de ve çoğu Windows 10'da vardır; yoksa pencere Edge'in ya da Chrome'un bir uygulama penceresidir.
- **Linux:** pencere, programın içinde sistemin web görünümüdür (WebKitGTK); C derleyicisi gerekmez. Bilgisayarda
  WebKitGTK yoksa: `sudo apt install libwebkit2gtk-4.1-0` (Debian, Ubuntu); o olmadan pencere Chrome'un, Chromium'un
  ya da Edge'in bir uygulama penceresidir.

Derlenen program (`go build`) tek bir dosyadır: pencerenin sayfası ve bileşenler onun içindedir.

## En küçük pencere

```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(func(page *ui.Object) {})
}
```

800 x 600 px boyutunda, başlığı "OzdenUI" olan boş bir pencere. `go run .` onu açar ve siz kapatınca döner.

## Sayaçlı bir pencere

```go
package main

import (
	"strconv"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	count = 0
	label *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Sayaç", Width: 400, Height: 320})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16}, func() {
		label = ui.Label(ui.Props{Text: "0", FontSize: 64, Bold: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 12}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "−", OnClick: func() { change(-1) }})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "+", Tone: "cyan", OnClick: func() { change(1) }})
		})
	})
}

func change(step int) {
	count += step
	label.SetText(strconv.Itoa(count))
}
```

Yukarıdan aşağı okuyun:

- **`ui.Run(start, ui.Window{...})`** penceredir. Onu açar, sayfa hazır olunca `start`'ı çağırır ve pencere kapanınca
  döner. `main()`'in sonunda bir kez çağırın.
- **`start(page)`** nesneleri yapar. `page` pencerenin zeminidir: bir rengi, bir genişliği ve bir yüksekliği vardır;
  kendi kutusu olmayan her nesne onun üstündedir. Nesneler `start()`'ta ya da sonra (bir olayda, bir zamanlayıcıda)
  yapılır, hiçbir zaman `ui.Run()`'dan önce değil.
- **`ui.VGroup(ui.Props{Gap: 16}, func() { ... })`** bir gruptur: fonksiyonunun içinde yapılan nesneler ona girer, alt
  alta ve 16 px aralıkla. İçindeki `HGroup` iki düğmeyi yan yana koyar; `Hug: true` onu yalnız düğmeleri kadar büyük
  yapar.
- **`ui.Label(ui.Props{Text: "0", FontSize: 64, Bold: true})`** bir metindir. Değerleri `ui.Props`'un alanlarıdır.
- **`OnClick: func() { change(-1) }`** olaydır: bir tıklama fonksiyonu çalıştırır.
- **`label.SetText(...)`** etiketi değiştirir. `label.Text()` onu pencereye sormadan geri okur.

Her nesne bir `*ui.Object`'tir. Sonradan değiştireceklerinizi değişkenlerde saklayın (burada `label`); öbürlerine
gerek yoktur.

Modülün `examples/counter`'ı bunun gibi bir programdır.

## Nasıl çalışır

Pencere, programınızın içindeki bir sayfayı gösterir. Sayfada basic.js ve kitin bileşenleri vardır; modül onları
yanında getirir, internet gerekmez. Sayfa yalnız bu bilgisayara sunulur.

```
Go kodunuz ──komutlar──▶ sayfa (basic.js nesneleri) ──▶ gördüğünüz
      ▲                                  │
      └──────── olaylar (bir tıklama) ◀──┘
```

- Go'daki bir nesne küçük bir vekildir: bir kimlik ve kendisine verilen değerleri tutar, gerçek nesnenin yaşadığı
  sayfaya komutlar gönderir. Bir işin (`start`, bir olay, bir zamanlayıcı) bütün komutları, iş bitince sayfaya
  birlikte gider.
- Sayfadaki bir tıklama, Go fonksiyonunuzu pencerenin goroutine'inde çağırır (4. bölüm).
- Kodunuz düz Go'dur, istediğiniz her paketle: dosyalar, veritabanları, `net/http`, bir yapay zekâ istemcisi.

## Window

```go
ui.Window{
	Title:    "Notlar", // pencerenin başlığı (boşsa "OzdenUI")
	Width:    800,      // sayfanın px olarak boyutu (0 ise 800 x 600)
	Height:   600,
	Theme:    "light", // "light", "dark" ya da "auto" (sistem ne istiyorsa): 7. bölüm
	Style:    "",      // "" ya da "minimal": 7. bölüm
	Language: "tr",    // metinlerinin dili: "tr" büyük harfleri Türkçe yazar ("İ"); boşsa "en"
	Debug:    false,   // true: web görünümünün geliştirici araçları (sağ tık > Inspect)
}
```

Pencerenin öbür fonksiyonları `start()`'ta, bir olayda ya da bir zamanlayıcıda çalışır (`ui.Close()` her
goroutine'den de):

```go
ui.SetTitle("Notlar — kaydedildi")
ui.Close()          // pencereyi kapatır; Run() döner
ui.Theme()          // "light" ya da "dark"
ui.SetTheme("dark") // pencere başka bir temada yeniden: 7. bölüm
```

## Birden çok fonksiyonlu bir program

Çoğu program, sonradan değiştireceği nesneleri saklar. Dosyanın başındaki değişkenler bunun en basit yeridir:

```go
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	name   *ui.Object
	answer *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Selamlama", Width: 420, Height: 320})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14, Padding: 40}, func() {
		name = ui.TextField(ui.TextFieldProps{TitleText: "Adınız ne?", Placeholder: "Ada", Width: 300, OnEnter: greet})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Selamla", Tone: "cyan", OnClick: greet})
		answer = ui.Label(ui.Props{FontSize: 18})
	})
}

func greet() {
	typed := strings.TrimSpace(name.Text())
	if typed == "" {
		answer.SetText("Önce adınızı yazın.")
		return
	}
	answer.SetText("Merhaba, " + typed + "!")
}
```

`name.Text()` her zaman alandaki metindir: sayfa onu her tuşta gönderir. Bir `Label` yazılan metni olduğu gibi
gösterir (bir `<b>`, `<b>` olarak kalır).

Bir struct da aynı iyi çalışır (`app.name = ui.TextField(...)`, sonra `ui.Run(app.start)`): 9. bölüm kendi
bileşeninizi böyle yapar. `examples/greeting` yukarıdaki programdır.

## Bir pencereye onu açmadan bakmak

```
OZDENUI_LOOK=app.png go run .
```

programı olduğu gibi çalıştırır, penceresinin bir resmini `app.png` olarak kaydeder ve biter; hiçbir pencere açılmaz.
Bilgisayarda Chrome, Edge ya da Chromium gerekir. `OZDENUI_THEME=dark` koyu pencereyi gösterir, `OZDENUI_LOOK_WAIT=3`
resimden önce üç saniye bekler. Hatasız çalışan bir programın da bir parçası örtülü ya da kesik olabilir: resme bir
bakış bunu söyler.

## Sayfa yeniden yüklenince

`Debug: true` ile sayfayı yeniden yüklemek (sağ tık > Reload) nesneleri yeniden yapar: `start()` yeniden çalışır,
zamanlayıcılar durur ve terminal `ozdenui: the page was loaded again: start() runs again` yazar. Verisini Go'da tutan
bir program (bir slice, bir dosya) onu yeniden gösterir; başka bir şey gerekmez.

## Sonraki

[2. bölüm: Nesneler](02-objects.md): beş düz nesne ve değerleri.
