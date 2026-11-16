# 9. Kendi bileşeniniz

Birden çok kez gereken bir parça (bir kart, bir formun bir satırı, bir puanlama) sizin bir türünüz olur: nesnelerini
yapan, değiştireceklerini saklayan ve metotları olan, kurucu fonksiyonlu bir struct. JavaScript gerekmez, bir yere
kaydetmek de gerekmez: düz Go'dur.

## Bir kart

```go
package main

import ui "github.com/bug7a/ozdenui-go"

// Card bu programın bir bileşenidir: bir başlık, bir metin ve bir düğme.
type Card struct {
	Box   *ui.Object // kartı tutan grup
	title *ui.Object
}

// NewCard çağrıldığı grupta bir kart yapar. Düğmesine tıklanınca onOpen çalışır.
func NewCard(title, text string, onOpen func(card *Card)) *Card {
	card := &Card{}
	card.Box = ui.VGroup(ui.Props{Width: 220, Height: "auto", Align: "left top",
		Padding: 16, Gap: 8, Color: "white", Round: 12}, func() {
		card.title = ui.Label(ui.Props{Text: title, FontSize: 18, Bold: true})
		ui.Label(ui.Props{Text: text, FontSize: 15, TextColor: "#5A6270"})
		ui.SmartButton(ui.SmartButtonProps{LabelText: "Aç", OnClick: func() {
			if onOpen != nil {
				onOpen(card)
			}
		}})
	})
	return card
}

// SetTitle başlığını değiştirir.
func (card *Card) SetTitle(title string) {
	card.title.SetText(title)
}

func main() {
	ui.Run(start, ui.Window{Title: "Kartlar", Width: 760, Height: 240})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.HGroup(ui.Props{Gap: 16, Padding: 24, Wrap: true}, func() {
		for _, name := range []string{"Satış", "Stok", "Siparişler"} {
			NewCard(name, "Bugün güncellendi.", func(card *Card) { card.SetTitle("Açıldı") })
		}
	})
}
```

- **Kurucu fonksiyon nesneleri çağrıldığı yerde yapar**, kitin bir bileşeninin yaptığı gibi: bir grubun fonksiyonunun
  içinde kart o gruba girer.
- **Struct nesnelerini saklar** (`Box`, `title`). Dışa açık bir alan (`Box`) kartı kullanan program içindir; küçük harfli
  bir alan (`title`) içeride kalır.
- **Metotlar onu değiştirir** (`SetTitle()`), bileşenlerin metotlarının yaptığı gibi.
- **Nesneleri olağan nesnelerdir:** `card.Box.Remove()` kartı ve içindeki her şeyi siler;
  `card.Box.On("mouseenter", fn)` ve `card.title.Measure("width")` her nesnede olduğu gibi çalışır.
- Bir struct nesneyi bunun yerine gömebilir (`type Card struct { *ui.Object; title *ui.Object }`, sonra
  `card.Object = ui.VGroup(...)`): kartın bir nesnenin bütün metotları olur, `card.Hide()`, `card.Remove()`.

Hiçbir şey saklamayan bir parçanın türe ihtiyacı yoktur: nesnelerini yapıp `*ui.Object`'i döndüren bir fonksiyon
yeter.

```go
// caption bir sayfanın bir bölümünün üstündeki küçük başlıktır.
func caption(text string) *ui.Object {
	return ui.Label(ui.Props{Text: strings.ToUpper(text), FontSize: 11, Bold: true, TextColor: "#5A6270"})
}
```

## Kendi durumu ve kendi callback'i

Sizin bir bileşeniniz durumunu Go'da tutar ve kullanıcı onu değiştirince programa haber verir:

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const (
	starOn  = "#E0A038"
	starOff = "#C9CDD4"
)

// Rating beş yıldızdır: Value (0 ile 5 arası); kullanıcı bir yıldıza tıklayınca OnChange çalışır.
type Rating struct {
	Value    int
	OnChange func(rating *Rating)
	Box      *ui.Object
	stars    []*ui.Object
}

func NewRating(value int, onChange func(rating *Rating)) *Rating {
	rating := &Rating{Value: value, OnChange: onChange}
	rating.Box = ui.HGroup(ui.Props{Hug: true, Gap: 4}, func() {
		for number := 1; number <= 5; number++ {
			star := ui.Label(ui.Props{Text: "★", FontSize: 28, OnClick: func() { rating.click(number) }})
			rating.stars = append(rating.stars, star)
		}
	})
	rating.show()
	return rating
}

// SetValue kod içindir: bileşenlerin yaptığı gibi OnChange'i çağırmaz.
func (rating *Rating) SetValue(value int) {
	rating.Value = value
	rating.show()
}

func (rating *Rating) click(number int) {
	rating.SetValue(number)
	if rating.OnChange != nil {
		rating.OnChange(rating)
	}
}

func (rating *Rating) show() {
	for index, star := range rating.stars {
		color := starOff
		if index < rating.Value {
			color = starOn
		}
		star.Set("textColor", color)
	}
}

func main() {
	ui.Run(start, ui.Window{Title: "Puan", Width: 320, Height: 200})
}

func start(page *ui.Object) {
	var said *ui.Object
	ui.VGroup(ui.Props{Gap: 12}, func() {
		said = ui.Label(ui.Props{Text: "Puan verin", FontSize: 16})
		NewRating(3, func(rating *Rating) { said.SetText(fmt.Sprintf("%d yıldız", rating.Value)) })
	})
}
```

`func() { rating.click(number) }` her yıldıza kendi numarasını verir: Go 1.22'den beri bir `for` döngüsünün her turunun
kendi değişkeni vardır.

## Nesne tutan bir bileşen

Nesneleri, `ui.VGroup`'un yaptığı gibi bir fonksiyon olarak alın ve onları tutan kutuyu döndürün; böylece program
sonra `Add()` ile daha fazlasını ekleyebilir:

```go
// Section bir formun başlıklı bir bölümüdür. content nesnelerini yapar; section.Add(...) sonradan daha fazlasını koyar.
func Section(title string, content func()) *ui.Object {
	var body *ui.Object
	ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 8, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: strings.ToUpper(title), FontSize: 11, Bold: true, TextColor: "#5A6270"})
		body = ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 10, Align: "left top"}, content)
	})
	return body
}

sound := Section("Ses", func() {
	ui.CheckBox(ui.CheckBoxProps{LabelText: "Ses çal"})
	ui.Slider(ui.SliderProps{Value: 7, MinValue: 0, MaxValue: 10})
})
sound.Add(func() { ui.Toggle(ui.ToggleProps{Value: true}) })
```

## Paylaşmak: bir paket

Birden çok dosyanın kullandığı bir bileşen kendi paketine girer: modülünüzün bir klasörü.

```
myapp/
  go.mod          module myapp
  main.go         import "myapp/widgets"  ...  widgets.NewCard("Satış", "Bugün güncellendi.", nil)
  widgets/
    card.go       package widgets  (ui "github.com/bug7a/ozdenui-go" paketini alır)
```

- Öbür paketlerin kullandığı adlar büyük harfle başlar: `NewCard`, `Card.Box`, `SetTitle`.
- Birçok program için pakete kendi modülünü verin (`go mod init github.com/you/widgets`) ve her programda
  `go get github.com/you/widgets` yazın.
- Dışa açık her adın üstündeki bir `go doc` yorumu, kullanıcısına onun ne aldığını söyler; kitin fonksiyonlarının
  yaptığı gibi.

## Hangisi Go'nun, hangisi sayfanın

**"Bir tıklama, sonra bir değişiklik" Go'nundur.** Fonksiyonunuz çalışır, sayfa değişikliği fonksiyon bitince gösterir.

**"Her kare, farenin her hareketi" sayfanındır.** Go'dan gönderilen çok adımlı bir hareket takılarak ilerler. Sayfaya
nasıl hareket edeceğini söyleyin, sonra varılacak yeri verin:

```go
ball.Call("setMotion", "left 0.3s, opacity 0.2s") // bu değerler artık yumuşakça değişir
ball.Set("left", 300)                             // sayfa onu 0.3 sn'de oraya götürür
```

Bir `Canvas` bunu kendiliğinden yapar: serbest çizgileri ve sürüklenen şekilleri sayfa çizer, Go yalnız sonlarını duyar.

### Hemen cevap vermesi gereken bir callback

Bileşenlerin bazı callback'leri hemen bir cevap döndürmek zorundadır (bir listenin bir satırını çizen bir fonksiyon,
bir tıklamayı durduran bir denetim). Go'nun cevabı biraz sonra gelir; bu yüzden böyle bir callback bir Go fonksiyonu
olamaz. Bunun yerine değerleri sayfaya gitmeden önce Go'da biçimlendirin ve bileşenlerin verdiği parametreleri kullanın
(10. bölüm bunları sayar).

## Sonraki

[10. bölüm: Daha ileri](10-going-further.md).
