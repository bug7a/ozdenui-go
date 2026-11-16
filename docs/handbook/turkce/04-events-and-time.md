# 4. Olaylar ve zaman

## Bir olay bir fonksiyondur

```go
ui.Button(ui.Props{Text: "Kaydet", OnClick: save}) // yapılırken
button.Set("onClick", save)                        // sonra: öncekinin yerini alır
remove := button.On("click", logClick)             // bir işleyici daha; remove() onu kaldırır
```

`ui.Props`'un olay alanları: `OnClick`, `OnInput`, `OnEnter`, `OnKeyDown`, `OnMouseEnter`, `OnMouseLeave`, `OnFocus`,
`OnBlur`. Sayfanın başka her olayı `On()` ile çalışır; sayfanın yazdığı gibi yazılır: `"dblclick"`, `"mouseenter"`,
`"keyup"`, `"contextmenu"`.

### İşleyici ne gerekiyorsa onu alır

```go
func save()                                  // hiçbir şey
func save(button *ui.Object)                 // nesneyi
func save(button *ui.Object, event ui.Event) // nesneyi ve olayı
```

Bir olayın alanı `any`'dir, bu yüzden üçü de uyar. Fonksiyon olmayan bir değer, programınızın satırıyla panic verir.
Bir metot değeri de çalışır: `OnClick: app.save`.

**Bir döngüde her closure'ın kendi değişkeni vardır** (Go 1.22 ve sonrası), bu yüzden akla gelen kod doğrudur:

```go
for number := 1; number <= 5; number++ {
	ui.SmartButton(ui.SmartButtonProps{LabelText: strconv.Itoa(number), OnClick: func() { pick(number) }})
}
```

`ui.Event`'te olayın bildikleri vardır:

| Alan | Ne zaman |
|---|---|
| `Key`, `Code` | bir tuş: `"Enter"`, `"a"`, `"ArrowLeft"`; `Code` tuşun klavyedeki yeridir (`"KeyA"`) |
| `Shift`, `Ctrl`, `Alt`, `Meta` | bir tuş: basılı tutulurken `true` (`Meta`, Mac'te Cmd'dir) |
| `X`, `Y` | fare, pencerenin köşesinden px olarak |
| `Text` | bir `Input`: metni |
| `Name` | olayın adı: `"click"`, `"keydown"` |

```go
package main

import (
	"fmt"
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	spotColor  = "#5BBAC1"
	hoverColor = "#E0A038"
)

var (
	pressed *ui.Object
	spot    *ui.Object
	where   *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Olaylar", Width: 420, Height: 300})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		pressed = ui.Label(ui.Props{Text: "Bir tuşa basın", FontSize: 28, Bold: true})
		spot = ui.Box(ui.Props{Width: 160, Height: 60, Color: spotColor, Round: 10})
		where = ui.Label(ui.Props{Text: "Kutuya tıklayın", FontSize: 14, TextColor: "#5A6270"})
	})
	page.On("keydown", showKey)
	spot.On("mouseenter", func() { spot.Set("color", hoverColor) })
	spot.On("mouseleave", func() { spot.Set("color", spotColor) })
	spot.On("click", func(box *ui.Object, event ui.Event) {
		where.SetText(fmt.Sprintf("Tıklanan yer: %.0f, %.0f", event.X, event.Y))
	})
}

func showKey(page *ui.Object, event ui.Event) {
	held := []string{}
	if event.Ctrl {
		held = append(held, "Ctrl")
	}
	if event.Alt {
		held = append(held, "Alt")
	}
	if event.Shift {
		held = append(held, "Shift")
	}
	if event.Meta {
		held = append(held, "Meta")
	}
	pressed.SetText(strings.Join(append(held, event.Key), " + "))
}
```

`page.On("keydown", ...)` penceredeki her tuşu duyar; odaktaki bir alan kendi tuşlarını da duyar.

## Tek goroutine: kilit yok

`start()`, her olay ve her zamanlayıcı **sırayla, tek bir goroutine'de** (uygulamanın goroutine'inde) çalışır. İki
tıklama hiçbir zaman aynı anda çalışmaz ve okuduğunuz bir değer şu anın değeridir. Kilide gerek yoktur; programınızın
yalnız bu fonksiyonların kullandığı değişkenleri için de.

Bunun öbür yüzü: bir fonksiyon çalışırken pencere bekler. 3 saniye süren bir fonksiyon pencereyi 3 saniye dondurur.
Bu yüzden **bir olayda hiçbir zaman `time.Sleep()` kullanmayın**: bir zamanlayıcıyla bekleyin, yavaş bir işe de bir
goroutine verin (aşağıda).

## Zaman: Every() ve After()

```go
package main

import (
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	clock     *ui.Object
	stopClock func()
)

func main() {
	ui.Run(start, ui.Window{Title: "Saat", Width: 360, Height: 220})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 12}, func() {
		clock = ui.Label(ui.Props{Text: time.Now().Format("15:04:05"), FontSize: 48, Bold: true})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Saati durdur", OnClick: func() { stopClock() }})
	})
	stopClock = ui.Every(1000, tick)                      // tekrar tekrar, her 1000 ms'de
	ui.After(5000, func() { ui.SetTitle("Saat — 5 sn") }) // bir kez
}

func tick() {
	clock.SetText(time.Now().Format("15:04:05"))
}
```

- `ui.Every(ms, fn)` onu durduran bir fonksiyon döndürür; `ui.After(ms, fn)` onu iptal eden bir fonksiyon.
- Bir zamanlayıcı da bir iştir: uygulamanın goroutine'inde, tıklamaların arasında çalışır, bu yüzden her nesneyi
  değiştirebilir. Öncekisi hâlâ çalışırken gelen bir tık birikmez.
- `OnClick: stopClock` değil, `OnClick: func() { stopClock() }`: düğme yapılırken değişken henüz `nil`'dir; closure onu
  tıklama anında okur.
- `examples/clock` buna bir kronometre ekler.

## Yavaş bir iş: bir goroutine ve ui.Do()

```go
package main

import (
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const steps = 100

var (
	bar    *ui.Object
	button *ui.Object
	result *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Yavaş bir iş", Width: 420, Height: 260})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		bar = ui.ProgressBar(ui.ProgressBarProps{TitleText: "Çalışıyor", Value: 0, MaxValue: steps, Width: 320})
		button = ui.SmartButton(ui.SmartButtonProps{LabelText: "Başlat", Tone: "cyan", OnClick: begin})
		result = ui.Label(ui.Props{FontSize: 16})
	})
}

func begin() {
	button.Set("waiting", true) // düğmede bir saat: şimdi tıklama almaz
	result.SetText("")
	go work()
}

func work() {
	for step := 1; step <= steps; step++ {
		time.Sleep(30 * time.Millisecond) // yavaş kısım: büyük bir dosya, bir indirme, bir model
		ui.Do(func() { bar.SetValue(step) })
	}
	ui.Do(func() {
		button.Set("waiting", false)
		result.SetText("Bitti.")
	})
}
```

- `go work()`: iş kendi goroutine'inde çalışır; pencere çalışmaya devam eder.
- **Nesneler yalnız uygulamanın goroutine'inde değiştirilir.** Sizin bir goroutine'iniz değişikliklerini ona
  `ui.Do(func() { ... })` ile verir: fonksiyon orada, şu an çalışan işten sonra çalışır. `ui.Do`'yu her goroutine
  çağırabilir: bir sunucu, bir soket, dosyaları izleyen bir şey.
- Başka bir goroutine'den değiştirilen bir nesne panic verir ve mesaj çözümü söyler:
  `ozdenui: ... works only in start(), an event or a timer. From another goroutine: ozdenui.Do(func() { ... })`.
  Hiçbir şey sessizce bozulmaz.
- İşe gerekenleri başlamadan önce verin (bir alanın metni, bir yol): bir nesneyi okumak da uygulamanın goroutine'inin
  işidir.
- `ui.Close()` her goroutine'den çalışır.
- `examples/job-log` dosyaları böyle kopyalar ve yaptıklarını bir `LogView`'a yazar.

### İlerlemesini gösteren uzun bir iş

Bir işin komutları sayfaya iş bitince gider. Bir süre uygulamanın goroutine'inde kalması gereken bir iş (bir kerede
çok sayıda nesne) o ana kadar yaptıklarını `ui.Flush()` ile gönderebilir.

## Olayların sırası

Olaylar Go'ya oldukları sırayla ulaşır; çoğu birden geldiğinde de (hızlı yazarken). Her biri kendi işi olarak çalışır.
Hızlı ve sürekli bir olay (`obj.On("mousemove", fn)`) her hareketi Go'ya gönderir. Sürüklenen bir bileşen (bir
`Slider`, bir `Splitter`) hareketlerini sayfada yapar ve size sonucu söyler (`OnChange`, `OnChangeEnd`).

## Sonraki

[5. bölüm: Bileşenler](05-components.md): bir pencerenin hazır parçaları.
