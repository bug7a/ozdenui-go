# 7. Görünüm

En genişinden en küçüğüne üç düzey vardır: pencerenin **teması** (açık ya da koyu), pencerenin **stili** (olağan görünüm
ya da sizin renginizdeki minimal görünüm) ve **tek bir bileşenin görünümü** (paketi ve renkleri). Kodunuzun verdiği bir
görünüm, pencerenin görünümüne her zaman baskın gelir.

## Koyu bir pencere

```go
ui.Run(start, ui.Window{Theme: "dark"}) // "light" (varsayılan), "dark" ya da "auto": sistem ne istiyorsa
```

- Her bileşen koyu görünümünü alır; diyaloglar ve toast'lar da.
- Sayfa, bir `Label`, bir `Box` ve bir `Input` da uyar: koyu bir zemin, açık bir metin.
- Bir `SmartButton`'ın `Tone`'u koyu tonuna döner (`"cyan"` → `"cyanDark"`).
- `ui.IsDark()` koyu bir pencerede true'dur (`"auto"`dan sonra da), `ui.Theme()` ise `"light"` ya da `"dark"`'tır:
  kendi nesnelerinizin renkleri için kullanın.

```go
textSoft := map[bool]string{false: "#5A6270", true: "#9AA3B5"}[ui.IsDark()]
ui.Label(ui.Props{Text: "Bugün güncellendi", TextColor: textSoft})
```

`OZDENUI_THEME=dark go run .` herhangi bir programı, kodunu değiştirmeden koyu açar: ikisini de görmenin hızlı bir yolu.

## Program çalışırken temayı değiştirmek

On beş bileşen sayfadayken temasını değiştirebilir: `LogView`, `ImageView`, `MarkdownView`, `JsonView`, `ChatView`,
`RangeSlider`, `StepBar`, `Canvas`, `MenuBar`, `ToolBar`, `StatusBar`, `Splitter`, `CodeEditor`, `Accordion` ve
`GraphBox`:

```go
logView.Set("theme", "dark")
```

Öbürleri görünümlerini yapılırken alır; bu yüzden bütün bir pencere temasını sayfasını yeniden yaparak değiştirir.
Bunu `ui.SetTheme("dark")` yapar (`"light"`, `"dark"` ya da `"auto"`): sayfa yeniden yapılır ve **`start()` yeniden
çalışır**. Her nesne yenidir, bir alana yazılmış metin gider ve önceden kurulan zamanlayıcılar durur. Programınızın kendi
değişkenlerinde tuttukları kalır; bu yüzden veriyi `start()`'ın dışında tutun, nesneleri `start()` içinde ondan
doldurun ve zamanlayıcıları orada kurun:

```go
package main

import (
	"strconv"

	ui "github.com/bug7a/ozdenui-go"
)

// Bu programın her temadaki renkleri: anahtar ui.IsDark()'tır.
var (
	pageColor = map[bool]string{false: "#F4F5F7", true: "#141413"}
	cardColor = map[bool]string{false: "#FFFFFF", true: "#232322"}
	textSoft  = map[bool]string{false: "#5A6270", true: "#9AA3B5"}
)

var clicks = 0 // start()'ın dışında: sayfa yeniden yapılınca kalır

func main() {
	ui.Run(start, ui.Window{Title: "Tema", Width: 420, Height: 320})
}

func start(page *ui.Object) {
	dark := ui.IsDark()
	other := map[bool]string{false: "dark", true: "light"}[dark]
	page.Set("color", pageColor[dark])
	var count *ui.Object
	ui.VGroup(ui.Props{Gap: 16}, func() {
		ui.VGroup(ui.Props{Width: 300, Height: "auto", Padding: 20, Gap: 10, Color: cardColor[dark], Round: 12}, func() {
			ui.Label(ui.Props{Text: "Tıklamalar", FontSize: 13, TextColor: textSoft[dark]})
			count = ui.Label(ui.Props{Text: strconv.Itoa(clicks), FontSize: 40, Bold: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Tıkla", Tone: "cyan", OnClick: func() {
				clicks++
				count.SetText(strconv.Itoa(clicks))
			}})
		})
		ui.SmartButton(ui.SmartButtonProps{LabelText: map[string]string{"dark": "Koyu yap", "light": "Açık yap"}[other], OnClick: func() { ui.SetTheme(other) }})
	})
}
```

`examples/gallery`'nin ToolBar'ında bunun için bir Dark / Light düğmesi var. `examples/editor` ve
`examples/markdown-editor` yalnız çubuklarını ve görünümlerini `Call("setTheme", ...)` ile değiştirir; böylece
düzenleyicideki metin kalır.

## Sizin renginizde minimal görünüm

```go
ui.Run(start, ui.Window{Style: "minimal"})                    // küçük ve çizgilerden yapılmış; Theme: "dark" de çalışır
ui.Run(start, ui.Window{Style: "minimal", Accent: "#2A78D6"}) // aynısı, sizin renginizde
```

**`Style: "minimal"`** 40 kadar bileşene minimal görünümlerini verir: 28 px denetimler, 13 px metin, kendilerine ait
dolgu ve gölge yok; böylece pencerenin zeminini alırlar. Bir araç, bir düzenleyici, bir ayarlar paneli bu görünümde
sakin durur.

**`Accent`** minimal görünüme tek bir renk verir (yalnız `Style: "minimal"` ile çalışır). Renk, görünümün mürekkebini
çizdiği yerlere gider: açık bir anahtar, işaretli bir kutu, bir kaydırıcının ve bir ilerleme çubuğunun dolgusu, seçili
sekme, satır ve adım, yazarken bir alanın kenarı ve ana düğme. Çizgiler, metinler ve odak halkası oldukları gibi kalır.

```go
package main

import ui "github.com/bug7a/ozdenui-go"

const accent = "#2A78D6" // beyaz üstünde okunan bir renk: seçili sekmenin metni de bu renktir

func main() {
	ui.Run(start, ui.Window{Title: "Minimal", Width: 420, Height: 380, Style: "minimal", Accent: accent})
}

func start(page *ui.Object) {
	page.Set("color", map[bool]string{false: "#FFFFFF", true: "#1A1A19"}[ui.IsDark()])
	ui.VGroup(ui.Props{Gap: 14, Padding: 30, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Hesap", FontSize: 18, Bold: true})
		ui.Tabs(ui.TabsProps{Tabs: []string{"Genel", "Gizlilik", "Gelişmiş"}, Panels: false, Width: 300})
		ui.CheckBox(ui.CheckBoxProps{LabelText: "Beni hatırla", Checked: true})
		ui.HGroup(ui.Props{Hug: true, Gap: 10}, func() {
			ui.Toggle(ui.ToggleProps{Value: true})
			ui.Label(ui.Props{Text: "Sesler", FontSize: 13})
		})
		ui.ProgressBar(ui.ProgressBarProps{TitleText: "Yükleniyor", Value: 62, Width: 260})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "İptal"})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Tone: "accent"})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Sil", StyleName: "minimalDanger"})
		})
	})
}
```

- **`Tone: "accent"`** her pencerenin ana işlemidir: burada vurgu rengi, vurgusuz minimal görünümde görünümün mürekkebi,
  olağan görünümde `"cyan"`.
- **`StyleName: "minimalDanger"`**: minimal görünümde bir silme düğmesi (koyu pencerede koyu ikizi).
- **Koyu bir pencerede** renk, koyu zeminde okunana kadar açılır. `AccentDark: "#7FB6DE"` onu sizin seçmenizi sağlar.
- Minimal görünümü olmayan bir bileşen olağan görünümünü korur. `examples/minimal-theme` hepsini gösterir.

## Tek bir bileşenin görünümü

### Bir paket: StyleName

Çoğu bileşenin paket denen birkaç görünümü vardır. Fonksiyonun `go doc`'u paketlerini sayar
(`go doc github.com/bug7a/ozdenui-go CheckBox`):

```go
ui.CheckBox(ui.CheckBoxProps{LabelText: "Beni hatırla", StyleName: "modern"})
ui.ItemList(ui.ItemListProps{Items: photos, Layout: "grid", StyleName: "gallery"})
ui.SmartButton(ui.SmartButtonProps{LabelText: "İleri", StyleName: "pill", Tone: "blue"})
```

Adını verdiğiniz bir paket her pencerede kalır; koyu bir pencerede, koyu ikizi olan açık bir paket ikizini alır
(`"gallery"` → `"galleryDark"`, `"minimalDanger"` → `"minimalDangerDark"`). Bir bileşen paketiyle yapılır: sonraki bir
`Set("styleName", ...)` onu değiştirmez (terminal bunu söyler); yukarıdaki on beşin `"theme"`'i dışında.

### Sizin renkleriniz: Style

`Style` bir bileşenin parçalarını değiştirir. Anahtarları kitin adlarıdır, camelCase:

```go
ui.SmartButton(ui.SmartButtonProps{LabelText: "Git", Style: ui.P{
	"box":   ui.P{"color": "#2A78D6"},
	"hover": ui.P{"color": "#2367BA"},
	"label": ui.P{"textColor": "white"},
}})
ui.Toggle(ui.ToggleProps{Value: true, StyleName: "classic", BackgroundStyle: ui.P{"selectedColor": "#6AC966"}})
```

- `Style` paketin üstüne gelir: verdiğiniz kazanır, kalanı kalır.
- Parçaların adları (`box`, `hover`, `label`...) bileşenin basic.js'teki adlarıdır;
  [web/comp](https://github.com/bug7a/ozdenui-go/tree/main/web/comp) içindeki dosyası ve
  [katalog](https://bug7a.github.io/js-components/) onları gösterir.
- Adı `Style` ile biten bir alan (`BackgroundStyle`, `ButtonStyle`) de bir stildir.
- `Unframed: true` bileşenin çerçevesini kaldırır (kutusunda kenarlık ve yuvarlak köşe yok): iki çubuğun arasını
  dolduran bir `LogView` ya da `CodeEditor`.

### Düz nesneler

Bir `Box`, bir `Label`, bir grup renklerini değerlerinden alır (`Color`, `TextColor`, `BorderColor`): 2. bölüm. Bir
programın renklerini dosyasının başında sabitler olarak (ya da `ui.IsDark()`'a göre bir map olarak) tutun, böylece bir
değişiklik tek satır olur:

```go
const (
	pageColor = "#F4F5F7"
	cardColor = "#FFFFFF"
	textSoft  = "#5A6270"
)
```

## Sonraki

[8. bölüm: Veri ve resimler](08-data-and-pictures.md): grafikler, büyük tablolar, kayıtlar, Markdown, JSON, kod ve
resimler.
