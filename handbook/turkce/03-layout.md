# 3. Yerleşim

Kendi yeri olan nesne (`Left: 20, Top: 20`) bir oyun ya da bir çizim için iyidir. Bir programın penceresi ise
**gruplarla** daha kolay kurulur: gruplar nesnelerini yerleştirir, pencerenin boyu değişince de yeniden yerleştirir.

## VGroup ve HGroup

```go
ui.VGroup(ui.Props{Gap: 10, Padding: 20, Align: "left top"}, func() {
	ui.Label(ui.Props{Text: "Bir"})
	ui.Label(ui.Props{Text: "İki"}) // "Bir"in altında, 10 px aşağıda
	ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Evet"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Hayır"}) // "Evet"in yanında
	})
})
```

- `VGroup` nesnelerini alt alta koyar, `HGroup` yan yana.
- Bir grup, kutusu kadar büyüktür: `Width: "100%"`, `Height: "100%"`. `Hug: true` onu yalnız nesneleri kadar büyük
  yapar.
- `Gap`: nesneler arasındaki boşluk. `Padding`: kenarının içindeki boşluk: `20` ya da `[]int{x, y}` (`[]int{30, 12}`).
- `Align`: nesnelerin durduğu yer: `"center"` (varsayılan), `"left top"`, `"right center"`, `"center bottom"`...
- `Wrap: true`: yer kalmayınca nesneler bir sonraki satıra geçer (kartlar, etiketler).
- Bir grup bir `Box`'tır: `Color`, `Round`, `Border` onda da çalışır.

Bir grupta bir nesneye yalnız boyutu gerekir: `ui.Label(ui.Props{Text: "Bir", Width: 200})`. Boyutu verilmezse
içeriği kadar büyüktür.

**Bir grubun içindeki grup** hemen her zaman bunlardan birini ister: `Hug: true` (nesneleri kadar büyük), bir yükseklik
(`Height: "auto"` ya da bir sayı) ya da `Fill: true` (aşağıda). Bir `VGroup`'taki %100'lük iki grup birbirini
sıkıştırır.

**Sığmayınca:** kutusunun alabileceğinden çok nesnesi olan bir grup onları sıkıştırır ve biri ötekinin üstüne gelebilir
(bir kartın altında kalan bir başlık). Pencereye daha çok yer verin ya da içeriği bir `ScrollArea`'ya koyun (aşağıda).

## Kalan yeri almak: Fill

Alışılmış bir pencere: üstte bir çubuk, ortada bütün yeri alan bir bölüm, altta bir satır.

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

func main() {
	ui.Run(start, ui.Window{Title: "Notlar", Width: 640, Height: 420})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 0}, func() {
		ui.HGroup(ui.Props{Height: 52, Padding: []int{16, 0}, Gap: 10, Align: "left center", Color: "white"}, func() {
			ui.Label(ui.Props{Text: "Notlar", FontSize: 18, Bold: true})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Yeni", Tone: "cyan"})
		})
		ui.HGroup(ui.Props{Gap: 0, Fill: true}, func() {
			ui.VGroup(ui.Props{Width: 180, Height: "100%", Padding: 16, Gap: 8, Align: "left top", Color: "#E9EBEF"}, func() {
				for _, folder := range []string{"Tümü", "İş", "Ev"} {
					ui.Label(ui.Props{Text: folder, FontSize: 15})
				}
			})
			ui.ScrollArea(ui.ScrollAreaProps{Padding: 16, Gap: 12}, func() {
				for number := 1; number <= 40; number++ {
					ui.Label(ui.Props{Text: fmt.Sprintf("Not %d", number), FontSize: 15})
				}
			})
		})
		ui.Label(ui.Props{Text: "40 not", FontSize: 13, TextColor: "#5A6270", Height: 28})
	})
}
```

- `Fill: true` bir nesneye öbürlerinin bıraktığı yeri verir: bir `VGroup`'ta kalan **yüksekliği** (çubuk ile alt satır
  arasındaki orta sıra), bir `HGroup`'ta kalan **genişliği**. Her nesne ve her bileşen bunu alabilir:
  `ui.HGroup(ui.Props{Fill: true})`, `ui.LogViewProps{Fill: true}`, `ui.SplitterProps{Fill: true}`.
- `Fill` nesne yapılırken verilir; sonradan `Set()` ile değiştirilemez.
- `ui.ScrollArea` öbürlerinin bıraktığı yeri kendiliğinden alır (`Width` ve `Height`'ı kalan yerin `"100%"`'üdür) ve
  içindekileri kaydırır. Bir grupta değil de doğrudan sayfada yapılırsa bütün penceredir.
- İçindeki nesneler bir `VGroup`'taki gibi, sol üstten başlayarak alt alta dizilir: ona `Align`, `Gap` ve `Padding`
  verin (padding onlarla birlikte kayar). `ScrollX: true` yana da kaydırır. Kitin kaydırma çubuğu ondadır; pencereyle
  birlikte açık ya da koyu olur.

## GridGroup: sütunlar ve satırlar

```go
package main

import ui "github.com/bug7a/ozdenui-go"

func main() {
	ui.Run(start, ui.Window{Title: "Tuşlar", Width: 320, Height: 300})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Padding: 16, Align: "left top"}, func() {
		ui.GridGroup(ui.GridGroupProps{Columns: 4, Gap: 8}, func() {
			for _, key := range "789/456*123-0.=+" {
				ui.SmartButton(ui.SmartButtonProps{LabelText: string(key), Width: "100%"})
			}
			clearKey := ui.SmartButton(ui.SmartButtonProps{LabelText: "Temizle", Width: "100%", Tone: "orange"})
			ui.Current().Place(clearKey, ui.P{"columnSpan": 4})
		})
	})
}
```

- Nesneler hücrelere sırayla girer: soldan sağa, satır satır.
- `Columns: 4`: eşit sütunlar. `Columns: []any{120, "1fr", "auto"}`: 120 px, kalan yer, içeriği kadar geniş.
- Bir nesne hücresini doldurur; ızgaraya verilen `Align: "left center"` ise her nesneyi hücresinin içine yerleştirir.
- `grid.Place(obj, ui.P{"columnSpan": 2})`, `grid.Place(obj, ui.P{"column": 1, "row": 0})`: kendine ait bir hücre
  (0'dan sayılır).
- **Bir grubun fonksiyonunun içinde, grubun kendi değişkeni henüz atanmamıştır** (`grid := ui.GridGroup(..., func() {
  grid... })` orada `nil`'dir): grup `ui.Current()`'tir.
- Öbür gruplardan farklı olarak bir ızgara satırları kadar yüksektir (`Height: "auto"`) ve `Padding`'i yoktur: boşluğu
  çevresindeki bir grup verir. `examples/calculator` eksiksiz bir hesap makinesidir.

## Gruptan taşanlar

Bir grup kenarından taşanı kesmez: bir düğmenin odak halkası, bir gölge ya da fare üstüne gelince büyüyen bir düğme
küçük bir grupta da bütün görünür:

```go
ui.HGroup(ui.Props{Hug: true, Gap: 8}, func() {
	ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Tone: "cyan"})
	ui.SmartButton(ui.SmartButtonProps{LabelText: "İptal"})
})
```

Kesmenin amaç olduğu yerde (bir çerçevedeki resim) `ClipContent: true` keser. Kayan bir grup (`ScrollY: true`) ve bir
`ScrollArea` kendiliğinden keser.

İki varsayılan basic.js'tekiyle aynı değildir; böylece program daha az söyler:

- basic.js'te bir `VGroup` / `HGroup` içeriğini keser; Go'da kesmez.
- `MenuBar`, `ToolBar`, `StatusBar`, `SearchInput` ve `FilePicker` bir grupta küçülmez (`Shrink: 0`): bir `ScrollArea`
  ya da `Fill: true` bir nesne kalan yeri alırken boylarını korurlar.

## Kendi yeri olan nesne

Bir grubun dışında, bir `Box`'ta ya da sayfada, bir nesne onu koyduğunuz yerde durur: `Left`, `Top`, `Width`,
`Height`. Bir oyun, bir çizim ya da başka bir şeyin üstüne gelen bir katman için yol budur:

```go
ui.Box(ui.Props{Width: 400, Height: 300, Color: "#1E2230"}, func() {
	player = ui.Box(ui.Props{Left: 180, Top: 130, Width: 40, Height: 40, Color: "#5BBAC1", Round: 8})
})
player.Set("left", 200) // onu kutusunun solundan 200 px'e götürür
```

Yeri kendi değişkeninizde tutun (`playerLeft += 20`, sonra `player.Set("left", playerLeft)`): onu geri okumak sayfaya
bir sorudur. `examples/game` bir kutuyu ok tuşlarıyla böyle hareket ettirir. Bir boyut metin olabilir:
`Width: "100%"` kutusunu doldurur.

## Pencerenin boyu değişince

Gruplar nesnelerini kendiliğinden yeniden yerleştirir. Kendi hesaplarınız için `page.Width()` ve `page.Height()` her
zaman günceldir; `page.On("resize", fn)` de size haber verir:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var size *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Boyut", Width: 420, Height: 240})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	size = ui.Label(ui.Props{Left: 20, Top: 60, FontSize: 28, Bold: true})
	ui.Label(ui.Props{Left: 20, Top: 20, Text: "Pencerenin boyunu değiştirin.", FontSize: 15, TextColor: "#5A6270"})
	showSize(page)
	page.On("resize", showSize)
}

func showSize(page *ui.Object) {
	size.SetText(fmt.Sprintf("%.0f x %.0f", page.Width(), page.Height()))
}
```

`page.On("resize", fn)` her olayın aldığı türden fonksiyonları alır (4. bölüm): burada sayfayı alan bir tane.

## Sonraki

[4. bölüm: Olaylar ve zaman](04-events-and-time.md): tıklamalar, tuşlar, zamanlayıcılar ve uzun işler.
