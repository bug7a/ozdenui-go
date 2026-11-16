# 8. Veri ve resimler

Bu bölümün bileşenleri, bir Go programında çoğunlukla bulunanları gösterir: sayılar, satırlar, kayıtlar, metinler, kod
ve resimler. Hiçbiri başka bir modül istemez: onları sayfa çizer.

## Grafikler: GraphBox

```go
ui.GraphBox(ui.GraphBoxProps{Type: "bar", Labels: []string{"Oca", "Şub", "Mar"}, Values: []int{1200, 1900, 700},
	FormatValue: "{:,.0f} TL"})
```

- `Type`: `"line"`, `"area"`, `"bar"`, `"pie"`, `"doughnut"`, `"radar"`, `"scatter"`, `"bubble"`.
- `Values` herhangi bir sayı türünden bir slice'tır. Daha çok çizgi ya da çubuk:
  `Datasets: []ui.P{{"label": "2025", "data": []int{8, 11, 9}}, {"label": "2026", "data": []int{10, 14, 12}}}`.
- `FormatValue` bir fonksiyon değil bir biçim metnidir, çünkü sayfa çizerken bir sayının metnine ihtiyaç duyar:
  `"{}"`, `"{:.2f}"`, `"{:,}"`, `"{:,.0f} TL"`, `"{:.1%}"`. Bir biçim alan metot onu `ui.Format` ile alır:
  `graph.Call("setFormatValue", ui.Format("{:.0f} kg"))`.
- Sonra: `graph.Call("setValues", []int{...})`, `graph.Call("setType", "line")`, `graph.Call("addValue", "Nis", 15)`.
  `OnClick: func(graph *ui.Object, item ui.P)`: `item["label"]`, `item["value"]`.

### Canlı bir grafik ve bir kayıt

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	points    = 30
	firstOnes = 12 // başta çizilen değerler, birer saniye arayla
)

var (
	chart   *ui.Object
	logView *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Bir sensör", Width: 560, Height: 520})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20}, func() {
		chart = ui.GraphBox(ui.GraphBoxProps{Type: "area", TitleText: "Sıcaklık", Labels: []string{}, Values: []float64{},
			Curve: "smooth", ShowPoints: false, MaxCount: points, FormatValue: "{:.1f}°", Width: "100%", Height: 240})
		logView = ui.LogView(ui.LogViewProps{Fill: true})
	})
	for back := firstOnes; back > 0; back-- {
		measure(time.Now().Add(-time.Duration(back) * time.Second))
	}
	ui.Every(1000, func() { measure(time.Now()) })
}

func measure(at time.Time) {
	value := math.Round((20+rand.Float64()*6-3)*10) / 10
	chart.Call("addValue", at.Format("15:04:05"), value) // `points` değerden sonra en eskisi gider
	if value > 22 {
		logView.Call("warn", fmt.Sprintf("Sıcak: %.1f°", value))
		return
	}
	logView.Call("info", fmt.Sprintf("%.1f°", value))
}
```

`examples/charts` canlı bir grafik, tıklanabilen çubuklar, bir halka grafik, hedefli iki çizgi ve PNG dosyası olarak
kaydedilen bir grafik gösterir (`graph.Ask("getImage", "image/png", 2)` onun veri adresini verir).

## Tablolar

Bir satır programınızın bir struct'ı olabilir: `json` etiketleri, sütunların andığı anahtarlardır. `OnSelect` satırı
sizin türünüzde geri verir:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

type person struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Orders int    `json:"orders"`
}

var people = []person{{"Ada", "Ankara", 12}, {"Can", "İstanbul", 7}, {"Ece", "Bursa", 21}}

var columns = []ui.P{
	{"name": "Ad", "dataTitle": "name", "width": 180, "sortable": 1},
	{"name": "Şehir", "dataTitle": "city", "width": 160, "sortable": 1},
	{"name": "Siparişler", "dataTitle": "orders", "dataType": "integer", "width": 100, "sortable": 1},
}

var (
	table  *ui.Object
	chosen *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Kişiler", Width: 520, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20, Align: "left top"}, func() {
		table = ui.SmartTable(ui.SmartTableProps{TitleDataList: columns, ItemDataList: people, Width: "100%",
			Height: 260, OnSelect: showPerson})
		ui.HGroup(ui.Props{Hug: true, Gap: 12, Align: "left center"}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Bir kişi ekle", OnClick: addPerson})
			chosen = ui.Label(ui.Props{Text: "Bir satıra tıklayın.", FontSize: 15})
		})
	})
}

func showPerson(table *ui.Object, row person) {
	chosen.SetText(fmt.Sprintf("%s, %s: %d sipariş", row.Name, row.City, row.Orders))
}

func addPerson() {
	people = append(people, person{fmt.Sprintf("Kişi %d", len(people)+1), "İzmir", 0})
	table.Set("itemDataList", people) // listenin bütünü yeniden: tablo onu çizer
}
```

- **`SmartTable`** yalnız görünen satırları yapar; bu yüzden çok büyük bir liste de hızlı kalır (`examples/csv-viewer`'da
  100000 satır için bir düğme var). Bir süzgeç kutusu ve başlıklarında sıralaması vardır. Bir satır bir `ui.P` de
  olabilir.
- Bir sütun: `name` (başlığı), `dataTitle` (satırın anahtarı), `width`, `dataType` (`"string"`, `"integer"`...),
  `sortable`. Değerler metin olarak gösterilir: onları Go'da biçimlendirin (`fmt.Sprintf("%.2f", price)`).
- **`DataTable`** birkaç yüz satıra kadar içindir; sütunları metinlerine göre genişler.
- `examples/csv-viewer` bir CSV dosyasını bir SmartTable'da açar.

## Resimler

Bir resim bir dosya yolu, bytes (bir PNG ya da bir JPEG), programınızın yaptığı bir `image.Image` ya da bir web adresi
olabilir. `ui.Picture()` bunlardan herhangi birini sayfanın gösterebileceği bir adrese çevirir:

```go
ui.Icon(ui.Props{Width: 48, Height: 48, Picture: "logo.png"})               // Picture bir yol, bytes ya da bir image.Image alır
view := ui.ImageView(ui.ImageViewProps{Height: 400, ImageFile: "foto.jpg"}) // adı File ile biten bir alan: bir yol
view.Call("setImage", ui.Picture(img))                                      // bir metotta: ui.Picture()
```

- Bir `ImageView` tekerlekle ya da iki parmakla yakınlaşır, sürükleyince kayar. Yeni bir resim yakınlığı ve yeri korur;
  bu yüzden bir program çizimini tekrar tekrar yenileyebilir.
- Sayfa bu bilgisayardaki bir dosyayı yolundan okuyamaz; bu yüzden onu Go okur ve içeriğini bir `data:` adresi olarak
  verir. Adı `File` ile biten bir alan ve bir `Icon`'un `Picture`'ı bunu kendiliğinden yapar; başka her yerde (bir
  listenin bir öğesi, bir stil, bir metot) `ui.Picture(path)` verin.
- Adres resmin bütününü taşır: çok kez gösterdiğiniz bir resim için `ui.Picture(...)`'ı bir değişkende tutun.

```go
package main

import (
	"image"
	"image/color"

	ui "github.com/bug7a/ozdenui-go"
)

const size = 240

var (
	view  *ui.Object
	shade = 0
)

func main() {
	ui.Run(start, ui.Window{Title: "Go'dan bir resim", Width: 420, Height: 420})
}

// gradient bir resmi, sizin bir programınızın yapacağı gibi, piksel piksel çizer.
func gradient(shade int) image.Image {
	picture := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			picture.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(shade), 255})
		}
	}
	return picture
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 12, Padding: 20}, func() {
		view = ui.ImageView(ui.ImageViewProps{Width: 300, Height: 300, ImageFile: ui.Picture(gradient(shade))})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Yeniden çiz", OnClick: drawAgain})
	})
}

func drawAgain() {
	shade = (shade + 80) % 256
	view.Call("setImage", ui.Picture(gradient(shade)))
}
```

`examples/image-viewer` bir resim dosyası açar ya da bir çizimi tekrar tekrar yeniler.

## Kayıtlar: LogView

```go
logView := ui.LogView(ui.LogViewProps{Height: 300})
logView.Call("info", "Dosya okunuyor")
logView.Call("success", "Bitti")
logView.Call("error", "Dosya orada değil")
log.SetOutput(logView.Writer("info"))             // log paketi onun içine yazar
fmt.Fprintln(logView.Writer("warn"), "Disk dolu") // herhangi bir io.Writer: bir goroutine'den de
```

Düzeyler `"info"`, `"success"`, `"warn"`, `"error"` ve `"debug"`'dır. Görünüm, okuyan kişi yukarı kaydırana kadar son
satırı izler. Bir metin hiçbir zaman HTML olarak okunmaz. `logView.Writer(level)` her goroutine'den kullanılabilir:
kendiliğinden `ui.Do()`'dan geçer ve her `Write` bir satırdır (ya da daha çok). `examples/job-log`.

## Metinler, JSON ve kod

| | |
|---|---|
| `ui.MarkdownView(ui.MarkdownViewProps{MarkdownText: "# Bir rapor\n\nBirkaç **kalın** kelime"})` | `view.Call("appendText", piece)` sona ekler: bir yapay zekânın geldikçe yazdığı cevap; `view.Set("markdownText", text)` |
| `ui.JsonView(ui.JsonViewProps{Value: data, Height: 400, OnSelect: fn})` | herhangi bir Go değeri (bir struct alanlarının sırasını korur); bir tıklama: `OnSelect: func(view *ui.Object, entry ui.P)`, `entry["pathText"]` (`users[0].name`) |
| `ui.JsonEditor(ui.JsonEditorProps{Value: rows, Height: 360, OnChange: fn})` | bir tabloda yazılan JSON listesi ya da map'i; `OnChange: func(editor *ui.Object, list []product)` değerin bütününü sizin türünüzde alır |
| `ui.CodeEditor(ui.CodeEditorProps{Value: "x = 1", Language: "python", Height: 300, OnChange: fn})` | metni sorulur: `editor.Ask("getValue")`; Python, JavaScript, JSON, HTML, CSS, SQL, Bash |
| `ui.ChatView(ui.ChatViewProps{EmptyText: "Bir şey sorun", OnSend: answer})` | `OnSend: func(chat *ui.Object, text string)`; bir cevap: `chat.Call("addMessage", ui.P{"key": key, "role": "assistant"})`, sonra `chat.Call("appendText", key, piece)`; `chat.Set("busy", true)` noktaları ve Durdur'u gösterir |
| `ui.WebView(ui.WebViewProps{HtmlText: page, Height: 400})` | bir web sayfası (`Url`) ya da betikleriyle bir HTML metni; `OnMessage: func(view *ui.Object, data T)` sayfasının `parent.postMessage` ile gönderdiğini alır |

- **Büyük ya da hızlı bir değer saklanmaz, sorulur**: bir `CodeEditor` metnini her tuşta göndermez. Gerektiğinde,
  `OnChange`'inde ya da kullanıcı kaydederken onu sorun (`ChangeDelay: 120` yazmanın durmasını bekler).
- Bir Markdown metni hiçbir zaman HTML olarak okunmaz ve içindeki hiçbir betik çalışmaz. Sayfa bir dosyayı yolundan
  okuyamaz: bu bilgisayardaki bir resim metne `ui.Picture(path)` olarak girer (`examples/markdown-editor` bunu yapar).
  `http` adresli bir resim bağlantı olarak kalır (güvendiğiniz bir metin için `RemotePictures: true`).
- Bir `ChatView` cevabının parçaları `ui.Do()` ile bir goroutine'den gelebilir: pencere çalışmaya devam eder, görünüm de
  okuyan yukarı kaydırana kadar sonu izler (`examples/chat`).
- `examples/markdown-editor`, `examples/json-viewer` ve `examples/editor` (Python programını çalıştıran küçük bir
  düzenleyici) bunları gösterir.

## Çizim: Canvas

```go
ui.Canvas(ui.CanvasProps{Width: "100%", Height: 300, Shapes: []ui.P{
	{"key": "sun", "type": "circle", "x": 100, "y": 80, "radius": 30, "fill": "orange", "draggable": 1},
	{"key": "path", "type": "line", "x1": 20, "y1": 200, "x2": 260, "y2": 200, "color": "tomato", "lineWidth": 4},
	{"key": "title", "type": "text", "x": 20, "y": 12, "text": "Güneşi sürükleyin", "fontSize": 16},
}, OnShapeMove: func(canvas *ui.Object, shape struct{ Key string }) { fmt.Println("taşındı:", shape.Key) }})
```

Sayfa `draggable` bir şekli kendisi hareket ettirir; `OnShapeMove` Go'ya onun nerede durduğunu söyler. `Drawing: true`
fareye bir kalem verir. Metotları (`addShape`, `moveShape`, `removeShape`, `getImage`...) `Call` ile çağrılır.

## Sonraki

[9. bölüm: Kendi bileşeniniz](09-your-own-component.md).
