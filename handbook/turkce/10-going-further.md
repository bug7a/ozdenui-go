# 10. Daha ileri

## Bir hatayı bulmak

- **Go'nun bir hatası** (yanlış bir alan, yanlış bir tür) derlenmez: `ui.Props`'ta ya da bir `...Props` türünde
  olmayan bir alan bir derleme hatasıdır ve editörünüz onu program çalışmadan işaretler.
- **Derleyicinin denetleyemediği bir ad** (`Set()` içinde, `Extra`'da, bir bileşenin bir `ui.P`'sinde) programınızın
  satırıyla bir kez söylenir:

  ```
  ozdenui: main.go:17: Label: "fontSzie" is not a value of Label (did you mean "fontSize"?)
  ozdenui: main.go:24: TextField: "warningTxt" is not a value of TextField (did you mean "warningText"?)
  ```

- **Olmayan bir metot** (`obj.Call("setValu", 3)`) ve **hiç verilmemiş bir değer** (`"auto"` bir label'ın
  `label.Get("width")`'i) da aynı yoldan söylenir; ikincisi `Measure()`'ı anar.
- **Bir olaydaki ya da bir zamanlayıcıdaki bir panic** pencereyi içindeki işinizle birlikte kapatmaz: terminal
  `ozdenui: a task failed: ...` yazar, ona götüren program satırlarınızı da; pencere çalışmaya devam eder.
- **Sizin kendi bir goroutine'inizden değiştirilen bir nesne**, mesajında çözümüyle panic verir:
  `ui.Do(func() { ... })`.
- **Sayfanın bir mesajı** (istenen şeyi yapamayan bir bileşen) `ozdenui (page ...): ...` olarak gelir.
- `ui.Window{Debug: true}` pencereye web görünümünün geliştirici araçlarını verir: sağ tık > Inspect.

Kısacası: programı çalıştırın ve **terminali okuyun**.

## Go'dan basic.js

basic.js'in ve kitin her şeyine ulaşılabilir; henüz Go adı olmayanlara da:

```go
gauge.Call("setValue", 80)            // sayfadaki nesnenin bir metodu; beklenmez
text := editor.Ask("getSelectedText") // bir metot ve sonucu
width := label.Measure("width")       // sayfanın şu anki değeri ("auto" bir boyut)
field.Part("input").Call("select")    // bir bileşenin içindeki bir kutu

// Herhangi bir JavaScript; label.JS() nesnenin kendisidir.
value := ui.Evaluate(label.JS() + ".elem.offsetWidth")

// Kitin herhangi bir bileşeni adıyla, alanı olmayan bir parametre de kitteki adıyla.
gauge = ui.Component("Gauge", ui.P{"value": 72, "width": 200})
ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Extra: ui.P{"ariaLabel": "Dosyayı kaydet"}})
```

- `Call()` işin öbür komutlarıyla birlikte gider. `Ask()`, `Measure()` ve `Evaluate()` sayfayı bekler; bu yüzden onları
  döngülerin dışında tutun.
- Sayfadan gelen şey düzdür: bir sayı `float64`, kitin bir bayrağı `1` / `0`, bir liste `[]any`, bir nesne
  `map[string]any`'dir. `ui.Number(v)`, `ui.Text(v)` ve `ui.IsOn(v)` onları okur.
- Bir bileşenin `go doc`'u metotlarını sayar; her birinin ne aldığı kitteki dosyasındadır:
  [web/comp](https://github.com/bug7a/ozdenui-go/tree/main/web/comp). Parametreler camelCase adlarını korur.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	title *ui.Object
	line  *ui.Object
	name  *ui.Object
	info  *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Ölçüldü", Width: 480, Height: 340})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14, Padding: 24, Align: "left top"}, func() {
		title = ui.Label(ui.Props{Text: "Sayfanın ölçtüğü", FontSize: 26, Bold: true})
		line = ui.Box(ui.Props{Width: 20, Height: 4, Color: "#2A78D6", Round: 2})
		name = ui.TextField(ui.TextFieldProps{TitleText: "Yeni bir başlık", Value: "Merhaba", Width: 300, OnEnter: retitle})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Kullan", Tone: "cyan", OnClick: retitle})
		info = ui.Label(ui.Props{Text: "-", FontSize: 14, TextColor: "#5A6270"})
	})
	line.Call("setMotion", "width 0.3s") // genişliği sayfa hareket ettirir; Go yalnız sonunu verir
	underline()
}

func retitle() {
	title.SetText(name.Text())
	underline()
	name.Part("input").Call("select") // alanın metni seçilir, bir sonrakine hazır
}

// underline çizgiyi başlık kadar geniş yapar: "auto" bir boyutu yalnız sayfa bilir.
func underline() {
	width := ui.Number(title.Measure("width"))
	line.Set("width", width)
	pixels := ui.Number(ui.Evaluate("window.devicePixelRatio"))
	info.SetText(fmt.Sprintf("Başlık %.0f px genişliğinde; bu ekranın piksel oranı %.0f.", width, pixels))
}
```

## Goroutine'ler, bir kez daha

| | nesneleri değiştirebilir | sayfaya sorabilir (`Ask`, `Measure`, `Evaluate`) ve bir diyalog gösterebilir |
|---|---|---|
| `start()`, olaylar, zamanlayıcılar (pencerenin goroutine'i) | evet | evet |
| sizin bir goroutine'iniz | `ui.Do(func() { ... })` içinde | `ui.Do(...)` içinde, ama goroutine onu beklemez |

Her goroutine `ui.Do()`'yu, `ui.Close()`'u, `ui.SetTheme()`'yi çağırabilir ve `logView.Writer(level)`'a yazabilir.

## Hız

Bir pencereyi yavaşlatan neredeyse her zaman iki şeyden biridir:

- **Bir döngüde sayfaya bir soru**: her `Measure()`, `Ask()` ya da `Evaluate()` bir cevap bekler. Programınızın
  verdiğini ya da bir bileşenin söylediğini okumak (`obj.Get()`, `field.Text()`, `box.Checked()`) hiçbir şey sormaz.
- **Pencerenin goroutine'inde yavaş bir iş**: pencere o bitene kadar cevap vermez. Onu bir goroutine'e koyun
  (4. bölüm).

Bir işin bütün komutları, iş bitince sayfaya birlikte gider; bu yüzden bir döngüde 500 nesne yapmak 500 değil, tek bir
mesajdır. Pencerenin goroutine'inde kalması gereken uzun bir iş, ilerlemesini `ui.Flush()` ile gösterebilir: o ana kadar
yapılan komutlar sayfaya hemen gider.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	limit = 10_000_000
	parts = 20
)

var (
	bar    *ui.Object
	result *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Asal sayılar", Width: 420, Height: 260})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16, Padding: 24}, func() {
		bar = ui.ProgressBar(ui.ProgressBarProps{TitleText: "Bakılan sayılar", MaxValue: parts, Width: 320})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Asalları say", Tone: "cyan", OnClick: countPrimes})
		result = ui.Label(ui.Props{Text: fmt.Sprintf("%d altındaki asallar: ?", limit), FontSize: 15})
	})
}

// countPrimes pencerenin goroutine'inde çalışır: Flush her parçayı biter bitmez gösterir.
func countPrimes() {
	count := 0
	step := limit / parts
	for part := 0; part < parts; part++ {
		for number := part * step; number < (part+1)*step; number++ {
			if isPrime(number) {
				count++
			}
		}
		bar.SetValue(part + 1)
		ui.Flush()
	}
	result.SetText(fmt.Sprintf("%d altındaki asallar: %d", limit, count))
}

func isPrime(number int) bool {
	if number < 2 {
		return false
	}
	for divisor := 2; divisor*divisor <= number; divisor++ {
		if number%divisor == 0 {
			return false
		}
	}
	return true
}
```

Çubuk ilerler, ama bir tıklama sayma bitene kadar bekler. Saniyeler süren bir iş, `ui.Do()` ile bir goroutine'e aittir.

## Pencerenin resmi, kişisiz bir denetim

```
OZDENUI_LOOK=app.png go run .                               # programı çalıştırır, app.png'yi kaydeder, biter; pencere yok
OZDENUI_THEME=dark OZDENUI_LOOK=dark.png go run .           # koyu pencere, kod değişmeden
OZDENUI_LOOK_WAIT=3 OZDENUI_LOOK=app.png go run .           # ilk çizimden 3 sn sonra (varsayılan 1.5)
OZDENUI_WINDOW=headless go run .                            # gizli bir pencere, program ui.Close()'u çağırana kadar
OZDENUI_WINDOW=browser go run .                             # Edge'in ya da Chrome'un uygulama penceresi, macOS'ta da
OZDENUI_WINDOW=tab go run .                                 # varsayılan tarayıcının bir sekmesi
```

- Bir resim için bilgisayarda Chrome, Edge ya da Chromium gerekir. Hatasız çalışan bir programın da bir parçası örtülü,
  kesik ya da sıkışmış olabilir: resim bunu gösterir.
- Bir programın denetimi onu koddan kapatır: `start()`'ın ilk satırı olarak `ui.After(3000, ui.Close)`. Bir diyaloğa
  ulaşmamalıdır: bir diyalog bir kişiyi bekler.

## Bir yapay zekâ asistanıyla

OzdenUI yeni olduğu için bir yapay zekâ asistanı (Claude Code, Codex, Cursor, Copilot...) onu henüz tanımaz ve bildiği
araçlardan tahmin yürütür. Tek bir komut, modülün kurallarını projenizin klasörüne yazar; asistan onları kendiliğinden
okur:

```
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents              # AGENTS.md dosyasına
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents CLAUDE.md    # Claude Code'un okuduğu dosyaya
go run github.com/bug7a/ozdenui-go/cmd/ozdenui@latest agents --print      # yalnızca yazdırır
```

- Kurallar kısadır: bir program nasıl kurulur, yerleşim, olaylar, bileşenler, kaçınılacak hatalar ve asistanın yaptığı
  pencereye bakması için `OZDENUI_LOOK`.
- Var olan bir dosyada satırlarınız kalır: kurallar dosyanın sonuna, iki işaretin arasına yazılır. Modülü yükselttikten
  sonra komutu yeniden çalıştırın: yalnızca o bölüm değişir.

## Verilecek tek bir dosya

```
go build -o notes .                       # tek dosya: pencerenin sayfası ve kit içindedir
GOOS=windows GOARCH=amd64 go build .      # Windows için bir program, başka bir bilgisayarda derlenmiş
```

Sayfa, basic.js ve bileşenler gömülüdür (`go:embed`); bu yüzden derlenmiş program başka bir dosya istemez. Küçük bir
program yaklaşık 15 MB'tır. macOS'ta derleme cgo ister (sistemin web görünümü); Windows ve Linux C derleyicisi istemez.

## v0.1'in sınırları

- Denenen sistem macOS'tur. Windows ve Linux, Edge'in ya da Chrome'un uygulama penceresini açar (hiçbiri yoksa
  varsayılan tarayıcının bir sekmesini); onları deponun testleri denetler, elle henüz denenmediler.
- Bir program için tek pencere.
- Hemen cevap vermesi gereken bir callback bir Go fonksiyonu olamaz, çünkü cevap sonra gelir: bir doğrulayıcı, bir
  satırın nesnelerini kuran bir fonksiyon (`createItemView`). Bir tablonun değerlerini Go'da biçimlendirin; bir
  `GraphBox`'ın ya da bir `Stepper`'ın sayıları bir biçim metni alır (8. bölüm).
- Pencereye bırakılan bir dosya bir `FilePicker`'a yolunu değil adını verir (6. bölüm).
- Bir `Tooltip`'in balonu sizin nesnelerinizi değil, bir metin gösterir.
- On beş bileşen sayfadayken temasını değiştirebilir; öbürleri onu yapılırken alır (7. bölüm).
- Hot reload yok: yeniden `go run .`.

## Bundan sonra

- Modülün **örnekleri**, en küçüğünden başlayarak: `counter`, `greeting`, `basics`, `notes`, `todo`, `form`,
  `settings`, `calculator`, `dialogs`, `clock`, `game`, `charts`, `csv-viewer`, `json-viewer`, `image-viewer`,
  `editor`, `markdown-editor`, `job-log`, `chat`, `gallery` (ana bileşenler), `minimal-theme`.
- Modülün **README**'si: her bileşen tek bir tabloda.
- **`go doc github.com/bug7a/ozdenui-go SelectBoxProps`**: bir bileşenin her parametresi, varsayılanıyla;
  `go doc github.com/bug7a/ozdenui-go SelectBox`: görünümleri ve metotları.
- **[basic.js](https://bug7a.github.io/basic.js/)** ve **[bileşenler](https://bug7a.github.io/js-components/)**:
  altındaki kütüphane ve her bileşenin canlı bir örneği.

[Başa dön](README.md)
