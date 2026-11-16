# 2. Nesneler

Her pencere beş düz nesneden yapılır: `Box`, `Label`, `Button`, `Input` ve `Icon`. 5. bölümün bileşenleri de bunlardan
yapılmıştır. Beşi de tek bir değer alır: `ui.Props{...}`.

| Nesne | Nedir | Kendi alanları |
|---|---|---|
| `ui.Box` | bir dikdörtgen; başka nesneleri tutan tek nesne | `Color`, `Border`, `Round`, `ScrollY` |
| `ui.Label` | bir metin | `Text`, `Html`, `FontSize`, `Bold`, `TextColor`, `TextAlign` |
| `ui.Button` | sistemin düz bir düğmesi (daha güzeli `ui.SmartButton`) | `Text` |
| `ui.Input` | tek satırlık bir metin alanı | `Text`, `Placeholder`, `OnEnter` |
| `ui.Icon` | bir resim | `Picture` (bir yol, bir PNG'nin ya da JPEG'in baytları, bir `image.Image`, bir adres) |

## Yer ve boyut

Bir grubun dışındaki nesnenin kendi yeri vardır: `Left`, `Top`, `Width`, `Height`, px olarak:

```go
ui.Box(ui.Props{Left: 20, Top: 20, Width: 100, Height: 100, Color: "orangered"})
```

Bir boyut metin de olabilir: `Width: "100%"` (içinde olduğu kutunun genişliği), `Height: "auto"` (içeriği kadar
yüksek). Yeri verilmemiş bir nesne, kutusunun sol üst köşesinden başlar.

## Değerler

Bir değer, nesne yapılırken `ui.Props`'un bir alanıdır, sonra bir ad:

```go
title := ui.Label(ui.Props{Text: "Rapor", FontSize: 24, Bold: true, TextColor: "#2B2B2B"})
title.SetText("Aylık rapor")                                     // ekranda değiştirir
title.Set("text", "Yıllık rapor", "fontSize", 28, "bold", false) // birden çoğunu bir kerede
fmt.Println(title.Text())                                        // "Yıllık rapor": Go'da okunur, sayfaya sorulmaz
```

- **Adlar:** bir alan `FontSize`'dır; aynı değer `Set()` ve `Get()` içinde, basic.js'in yazdığı gibi `"fontSize"`'dır
  (alanın adı, ilk harfi küçük). Yanlış bir ad kaybolmaz: terminal
  `ozdenui: main.go:12: Label: "fontSzie" is not a value of Label (did you mean "fontSize"?)` yazar.
- **Sayılar ve açık / kapalı değerler:** px olarak bir boyut bir sayıdır; açık / kapalı bir değer `true` / `false`'tur.
- **Renkler:** her CSS rengi: `"red"`, `"#2A78D6"`, `"rgba(0, 0, 0, 0.5)"`.
- **Yazılmayan bir alan varsayılan değerini korur.** Yazılan bir sıfır bir değerdir (`Gap: 0`). Tek istisna metindir:
  `""` "verilmedi" demektir, bu yüzden bilerek boş bırakılan bir metin `ui.Empty`'dir (5. bölüm).
- **Okumak:** `obj.Get("fontSize")` değeri `any` olarak verir; `obj.Str()`, `obj.Num()`, `obj.Int()` ve `obj.Bool()`
  onu bir türde verir: `title.Num("fontSize")` `28`'dir.

Her nesnede olan değerler:

| Alan | Örnek |
|---|---|
| `Color` | zemin: `Color: "#FFFFFF"`, `Color: "transparent"` |
| `Border`, `BorderColor` | `Border: 1, BorderColor: "#DDD"` |
| `Round` | köşeler: `Round: 8`; `Round: 100` bir kareden daire yapar |
| `Opacity` | `Opacity: 0.5` |
| `Visible` | `Visible: false` gizler; bir grupta öbürleri boşluğu kapatır |
| `Clickable` | olayı olan bir nesne tıklamaları kendiliğinden alır; `Clickable: false` onları geçirir |
| `Cursor` | `Cursor: "pointer"` |
| `Css` | CSS'in başka her şeyi: `Css: ui.P{"letterSpacing": "1px"}` |
| `Extra` | basic.js'in başka herhangi bir değeri, adıyla: `Extra: ui.P{"textSize": 12}` |

`ui.P`, `map[string]any`'dir. Onu yeniden göreceksiniz: bir stil, bir listenin bir öğesi, bir tablonun bir satırı.

### Vermediğiniz bir değer

`title.Text()` cevabı Go'dan verir, çünkü onu kodunuz verdi. Sayfanın belirlediği bir değer, ör. genişliği verilmemiş
bir etiketin genişliği, `Measure()` ile sorulur:

```go
width := title.Measure("width") // şu an ekrandaki genişlik, px olarak (bir any içinde bir float64)
```

Hiç verilmemiş bir değerin `Get()`'i `nil`'dir ve terminal bunun için `Measure()`'ı anar. `Measure()` sayfayı bekler:
onu bir döngüde değil, gerektiğinde kullanın. Sayfanın kendisi farklıdır: `page.Width()` ve `page.Height()` her zaman
günceldir.

### Metin ve HTML

Bir `Label`'ın ya da bir `Button`'ın `Text`'i tam yazıldığı gibi gösterilir; bu yüzden bir kullanıcının yazdığı ya da
bir dosyadan okunan metin olduğu gibi verilir (bir `<b>`, `<b>` olarak kalır). HTML, `Html` ile verilir:

```go
note.SetText(typedText)
title.Set("html", "<b>3</b> yeni")
```

basic.js'te `text` HTML'dir; Go'da basic.js'in düz metnidir. Bileşenlerin metinleri de yazıldığı gibi gösterilir
(`ui.SmartButtonProps{LabelText: typed}`); `AllowHtml: true` onları HTML yapar.

## Bir Box'ın içindeki nesneler

Bir `Box`'a verilen fonksiyonda yapılan nesneler onun içine girer. Yerleri kutunun köşesinden sayılır:

```go
package main

import (
	"strings"

	ui "github.com/bug7a/ozdenui-go"
)

// İçinde bir onay işareti olan beyaz bir daire, bir resim adresi olarak: dosya gerekmez.
const check = `data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'>` +
	`<circle cx='12' cy='12' r='11' fill='white'/><path d='M7.5 12.5l3 3 6-7' fill='none' stroke='%232B2B2B' stroke-width='2.4'/></svg>`

var (
	label *ui.Object
	input *ui.Object
)

func main() {
	ui.Run(start, ui.Window{Title: "Bir kart", Width: 360, Height: 220})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.Box(ui.Props{Left: 20, Top: 20, Width: 320, Height: 110, Color: "#5BBAC1", Round: 12}, func() {
		ui.Icon(ui.Props{Left: 20, Top: 25, Width: 60, Height: 60, Picture: check})
		label = ui.Label(ui.Props{Left: 100, Top: 36, Width: 200, Height: 38, Text: "Merhaba", FontSize: 26, Bold: true, TextColor: "white"})
	})
	input = ui.Input(ui.Props{Left: 20, Top: 150, Width: 200, Height: 40, Placeholder: "Bir metin yazın", OnEnter: change})
	ui.Button(ui.Props{Left: 230, Top: 150, Width: 110, Height: 40, Text: "Değiştir", OnClick: change})
}

func change() {
	typed := strings.TrimSpace(input.Text())
	if typed == "" {
		typed = "Merhaba"
	}
	label.SetText(typed)
}
```

- `ui.Box(props, func() { ... })`: fonksiyon hemen çalışır ve yaptıkları kutuya girer. Fonksiyon isteğe bağlıdır:
  `ui.Box(props)` boş bir kutudur.
- **Nesne yalnız bir `Box`'ın (ve kutu olan grupların) içinde durur.** Bir `Label`'a ya da bir `Button`'a verilecek
  fonksiyon yoktur. Bir düğmenin üstüne bir şey koymak için onu düğmenin yanında, aynı kutuda ve aynı yerde yapın.
- `Picture` bu bilgisayardaki bir yolu (`"photo.png"`), bir PNG'nin ya da JPEG'in baytlarını, programınızın yaptığı bir
  `image.Image`'ı ya da bir adresi alır. Bir SVG adresindeki bir rengin `#`'i `%23` olarak yazılır.
- `OnEnter: change`: alanda Enter, düğmenin yaptığını yapar.
- `examples/basics` bu karttır; rengi her basışta değişir.

## Gizlemek, silmek, sonradan eklemek

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var (
	note  *ui.Object
	info  *ui.Object
	shown = true
)

func main() {
	ui.Run(start, ui.Window{Title: "Gizle ve sil", Width: 420, Height: 300})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 14}, func() {
		note = ui.Label(ui.Props{Text: "Gizlenecek bir not", FontSize: 22, Color: "white", Round: 8, Padding: []int{14, 8}})
		ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Gizle / göster", OnClick: toggle})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Ölç", OnClick: measure})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Sil", Tone: "orange", OnClick: remove})
		})
		info = ui.Label(ui.Props{Text: "Genişlik sayfaya sorulur.", FontSize: 14, TextColor: "#5A6270"})
	})
}

func toggle() {
	if note == nil {
		return
	}
	shown = !shown
	if shown {
		note.Show()
	} else {
		note.Hide()
	}
}

func measure() {
	if note == nil {
		return
	}
	info.SetText(fmt.Sprintf("Notun genişliği %v px.", note.Measure("width")))
}

func remove() {
	if note == nil {
		return
	}
	note.Remove()
	note = nil // silinen bir nesne yeniden kullanılmaz
	info.SetText("Not silindi.")
}
```

- `Hide()` ve `Show()` `Visible`'ı değiştirir; gizlenen bir nesne değerlerini korur ve `Show()` onu geri getirir.
- `Remove()` nesneyi ve içindeki her nesneyi kalıcı olarak siler: yeniden göstermek için yenisini yapın. Silinen bir
  nesne yeniden kullanılmaz: yukarıdaki program onu unutur (`note = nil`).
- `box.Add(func() { ... })` bir kutuda ya da bir grupta sonradan, bir olayda ya da bir zamanlayıcıda, daha çok nesne
  yapar:

```go
list.Add(func() {
	ui.Label(ui.Props{Text: "Bir tane daha"})
})
```

## Bir kutuya sığmayınca

`ScrollY: true` bir kutunun içeriğini aşağı yukarı kaydırır (`ScrollX: true` yana). Bir grubun içinde 3. bölümün
`ScrollArea`'sı aynı işi yapar ve kendisine verilen yeri alır.

## Sonraki

[3. bölüm: Yerleşim](03-layout.md): gruplar nesneleri sizin yerinize yerleştirir.
