# 5. Bileşenler

Bir bileşen, pencerenin hazır bir parçasıdır ve 2. bölümün nesnelerinden yapılmıştır: bekleme saati olan bir düğme,
başlığı ve uyarısı olan bir alan, arama kutusu olan bir liste, 100000 satırlık bir tablo. Bunların 83'ünün kendi
fonksiyonu ve her birinin kendi `...Props` türü var; `ui.Component("Ad", ...)` herhangi birini adıyla yapar.

## Her bileşen nasıl çalışır

```go
remember := ui.CheckBox(ui.CheckBoxProps{LabelText: "Beni hatırla", Checked: true, OnChange: show})
```

- **Parametreler, Props türünün alanlarıdır.** Bir alan, kitin parametresinin ilk harfi büyük yazılmış hâlidir:
  `TitleText`, `SelectedKeys`, `PanelHeight`. `go doc github.com/bug7a/ozdenui-go CheckBoxProps` her alanı varsayılan
  değeri ve aldığı şeyle birlikte listeler; düzenleyiciniz de onları önerir. Yanlış bir ad derlenmez.
- **Yazılmayan bir alan kitin varsayılan değerini korur.** Yazılan bir sıfır bir değerdir (`Checked: false`,
  `Value: 0`). Tek istisna metindir: `""` "verilmedi" demektir, bu yüzden bilerek boş bırakılan bir metin `ui.Empty`'dir
  (`LabelText: ui.Empty`: yalnız simgesi olan bir düğme).
- **Bir callback bir alandır**: `OnChange: show`. Bir olay gibi ne gerekiyorsa onu alır: `show()`,
  `show(box *ui.Object)` ya da daha fazlası. Nesneden sonraki değerler, alanın açıklamasının saydıklarıdır; her biri
  fonksiyonunuzun istediği türde gelir: `OnSelect: func(table *ui.Object, row Person)` satırı bir `Person` olarak alır
  (`json` etiketlerine göre), `OnStepClick: func(bar *ui.Object, index int)` bir sayı alır.
- **Durumu Go'dadır**: `remember.Checked()`, `selectBox.Value()`, `field.Text()`, `list.Get("selectedKeys")`. Bunları
  okumak sayfaya bir şey sormaz; her fonksiyonun açıklaması neyi tuttuğunu söyler.
- **Değiştirmek**: `remember.SetValue(...)`, `bar.SetValue(40)`, `field.SetText("")` ya da adıyla bir parametre:
  `field.Set("unitText", "kg")`. Kitin bir metodu olan bir parametre (`setUnitText`) o metotla değiştirilir; böylece
  bileşen kendini yeniden çizer.
- **Kitin bir metodu**: `list.Call("addItem", item)` (beklenmez), `list.Ask("getSelectedKeys")` (sonucunu bekler;
  sayfaya bir sorudur, bu yüzden bir döngüde değil). Fonksiyonun açıklaması metotlarını listeler; her birinin ne
  aldığı [kitteki dosyasındadır](https://github.com/bug7a/ozdenui-go/tree/main/web/comp).
- **Kodunuzun yaptığı bir değişiklik callback'i çağırmaz**; yalnız kullanıcınınki çağırır.

## İki sekmeli bir form

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

const formWidth = 400

var name, language, news, volume, sounds *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Ayarlar", Width: 540, Height: 580})
}

func start(page *ui.Object) {
	page.Set("color", "#F4F5F7")
	ui.VGroup(ui.Props{Gap: 16, Padding: 24, Align: "left top"}, func() {
		ui.Label(ui.Props{Text: "Ayarlar", FontSize: 24, Bold: true})
		ui.VGroup(ui.Props{Width: 460, Height: "auto", Padding: 20, Gap: 16, Align: "left top", Color: "white", Round: 12}, func() {
			tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"Profil", "Ses"}, PanelHeight: 240})
			tabs.Panel("Profil", func() {
				ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 14, Padding: []int{0, 16}, Align: "left top"}, func() {
					name = ui.TextField(ui.TextFieldProps{TitleText: "Ad", Required: true, Width: formWidth})
					language = ui.SelectBox(ui.SelectBoxProps{Options: []string{"Türkçe", "English", "Deutsch"}, Value: "Türkçe", Width: formWidth})
					news = ui.CheckBox(ui.CheckBoxProps{LabelText: "Bana haberleri gönder"})
				})
			})
			tabs.Panel("Ses", func() {
				ui.VGroup(ui.Props{Width: "100%", Height: "auto", Gap: 16, Padding: []int{0, 16}, Align: "left top"}, func() {
					volume = ui.Slider(ui.SliderProps{Value: 7, MinValue: 0, MaxValue: 10, Width: formWidth})
					ui.HGroup(ui.Props{Hug: true, Gap: 18}, func() {
						for _, quality := range []string{"Düşük", "Normal", "Yüksek"} {
							ui.RadioButton(ui.RadioButtonProps{Group: "quality", Value: quality, LabelText: quality, Checked: quality == "Normal"})
						}
					})
					sounds = ui.Toggle(ui.ToggleProps{Value: true})
				})
			})
			ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Tone: "cyan", OnClick: save})
		})
	})
}

func save() {
	if !name.Check() { // alanın üstünde uyarıyı gösterir
		return
	}
	fmt.Println(name.Text(), language.Value(), news.Checked(), volume.Value(), ui.RadioValue("quality"), sounds.Checked())
	ui.ToastSuccess("Kaydedildi.")
}
```

`Value()` `any`'dir: kullanıcı değiştirene kadar kodunuzun verdiği değer (kaydırıcının `int` 7'si), sonra sayfanın
söylediği değer (sayfanın bir sayısı bir `float64`'tür). `ui.Number(v)`, `ui.Text(v)` ve `ui.IsOn(v)` iki türü de
okur: `ui.Number(volume.Value())` her zaman bir `float64`'tür. `examples/settings` üç sekmeli, daha uzun bir örnektir;
`examples/form` tarihleri, bir saati ve bir rengi olan bir formdur.

## Bileşenler, yaptıkları işe göre

### Düğmeler

| | |
|---|---|
| `ui.SmartButton(ui.SmartButtonProps{LabelText: "Kaydet", Tone: "cyan", OnClick: save})` | `Tone`: `"cyan"` ana işlem, `"orange"` bir silme, `"yellow"`, `"green"`, `"blue"`, `"white"`; `"accent"`: minimal bir pencerenin ana işlemi (7. bölüm) |
| `ui.SmartButtonProps{LabelText: ui.Empty, IconFile: "plus.svg", AriaLabel: "Ekle"}` | bir simge düğmesi: bir kare |
| `button.Set("waiting", true)` | üstünde bir saat; `button.Set("waiting", false)`'a kadar tıklama almaz |
| `ui.SmartButtonProps{LabelText: "Sil", Ask: "Silmek için yeniden tıklayın."}` | ilk tıklama sorar, ikincisi `OnClick`'i çalıştırır |
| `button.Set("enabled", false)`, `button.SetText("Kaydedildi")` | |

### Alanlar

| | |
|---|---|
| `ui.TextField(ui.TextFieldProps{TitleText: "Ad", Placeholder: "Ada", OnChange: fn, OnEnter: fn})` | `field.Text()`; `OnChange` her değişiklikte, `OnEnter` Enter'da çalışır |
| `EmailField`, `PasswordField`, `NumberField`, `TextAreaField` | aynısı; bir e-posta, bir parola, bir sayı, çok satır için |
| `Required: true`, sonra `field.Check()` | alanın üstünde uyarıyı gösterir; hiç uyarı yoksa `true` |
| `MaxLength: 40`, `ui.PasswordFieldProps{MinLength: 8}` | en çok ve en az harf sayısı |
| `ui.UrlField(ui.UrlFieldProps{TitleText: "Site"})` | protokolüyle bir adres |
| `ui.PhoneField(...)`, `ui.CurrencyField(ui.CurrencyFieldProps{UnitText: "$"})` | maskesindeki numara; üçlü gruplarla yazılan tutar |
| `ui.SearchInput(ui.SearchInputProps{Placeholder: "Ara", OnChange: fn})` | her değişiklikte `fn(search *ui.Object)`; `search.Text()` |
| `ui.TagInput(ui.TagInputProps{Tags: []string{"ses"}, Suggestions: []string{...}})` | `field.Get("tags")` |

`Check()` bir alanın uyarısını gösterir. Bir formun bütün uyarılarını göstermek için önce hepsini denetleyin, sonra
karar verin: `nameOK, mailOK := name.Check(), email.Check()` (`&&` ile ilki tutmayınca ikincisi denetlenmezdi).

### Seçimler

| | |
|---|---|
| `ui.CheckBox(ui.CheckBoxProps{LabelText: "Beni hatırla", Checked: true})` | `box.Checked()` |
| `ui.Toggle(ui.ToggleProps{Value: true})` | `toggle.Checked()` |
| `ui.RadioButton(ui.RadioButtonProps{Group: "size", Value: "M", LabelText: "Orta"})` | `ui.RadioValue("size")` |
| `ui.SelectBox(ui.SelectBoxProps{Options: []string{"Elma", "Armut"}, Value: "Armut"})` | `selectBox.Value()` |
| `ui.Slider(ui.SliderProps{Value: 30, MinValue: 0, MaxValue: 100, OnChangeEnd: fn})` | `slider.Value()`; `OnChangeEnd`: fare onu bırakınca |
| `ui.SliderField(...)`, `ui.Stepper(ui.StepperProps{Value: 2, MinValue: 0, MaxValue: 10})` | `.Value()` |
| `ui.RangeSlider(ui.RangeSliderProps{LowValue: 200, HighValue: 800, MaxValue: 1000, Step: 50})` | `slider.Get("low")`, `slider.Get("high")` |
| `ui.DatePicker(ui.DatePickerProps{Date: "2026-11-20"})`, `ui.TimePicker(ui.TimePickerProps{Time: "19:30"})`, `ui.ColorPicker(ui.ColorPickerProps{SelectedColor: "#3871E0"})` | `.Value()`: `"2026-11-20"`, `"19:30"`, `"#3871E0"` |

### Listeler, tablolar ve ağaçlar

```go
fruits := ui.ItemList(ui.ItemListProps{
	Items:     []any{"Elma", ui.P{"key": "p", "text": "Armut", "description": "Yeşil"}},
	Selection: "multiple", Height: 240,
	OnChange: func(list *ui.Object) { fmt.Println(list.Get("selectedKeys")) },
})
fruits.Call("addItem", ui.P{"key": "c", "text": "Kiraz"})
fruits.Call("removeItem", "p")
```

- **`ItemList`**: `list.Get("selectedKeys")` (anahtarların bir `[]any`'si), `list.Get("selected")`. Metin olarak
  verilen bir öğe, sırasını anahtar olarak alır (`"0"`); `"enabled": false` olan seçilemez.
  `Layout: "grid", StyleName: "gallery"` resimleri gösterir. `examples/todo` bir listeden yapılmış bir yapılacaklar
  listesidir.
- **Bir öğe** bir metin, bir `ui.P` ya da sizin bir struct'ınızdır: `json` etiketleri anahtarlarıdır. Anahtarlar kitin
  adlarıdır: `text`, `description`, `iconFile`, `enabled`...
- **`SmartTable`**: çok sayıda satır (yalnız görünen satırlar yapılır), bir süzgeç kutusu ve sıralama. Sütunları:
  `TitleDataList: []ui.P{{"name": "AD", "dataTitle": "name", "width": 160, "sortable": 1}}`; satırları:
  `ItemDataList` (struct'lar ya da `ui.P`); `OnSelect: func(table *ui.Object, row Person)`. 8. bölümde daha fazlası
  var.
- **`DataTable`** (`Columns`, `Rows`): daha küçük bir tablo; `table.Get("selectedRow")`. **`TreeView`**
  (`Items: []ui.P{{"key": "a", "text": "A", "children": []string{"b"}}}`): `tree.Value()` seçili anahtardır.
  **`SortableList`**: sırayı kullanıcı değiştirir; `list.Get("order")` anahtarları yeni sırasıyla verir.

### Sizin nesnelerinizi tutan alanlar

Bazı bileşenler sizin nesnelerinizi tutar. Her biri, bir kutunun yaptığı gibi, alanına bir fonksiyon verir:

```go
tabs := ui.Tabs(ui.TabsProps{Tabs: []string{"Siparişler", "Müşteriler"}})
tabs.Panel("Siparişler", func() { ui.Label(ui.Props{Text: "Henüz sipariş yok"}) })

folds := ui.Accordion(ui.AccordionProps{Sections: []string{"Genel", "Gelişmiş"}})
folds.Section("Genel", func() { ui.CheckBox(ui.CheckBoxProps{LabelText: "Sistemle birlikte başlat"}) })

split := ui.Splitter(ui.SplitterProps{FirstSize: 240}) // iki alan, aralarında sürüklenen bir çizgi
split.First(func() { ui.Label(ui.Props{Text: "Sol"}) })
split.Second(func() { ui.Label(ui.Props{Text: "Sağ"}) })

panel := ui.SidePanel(ui.SidePanelProps{TitleText: "Süzgeçler"}) // yandan kayarak gelir
panel.Body(func() { ui.CheckBox(ui.CheckBoxProps{LabelText: "Yalnız yeniler"}) })
ui.SmartButton(ui.SmartButtonProps{LabelText: "Süzgeçler", OnClick: func() { panel.Call("open") }})
```

Metin olarak verilen bir sekmenin ya da bölümün anahtarı o metindir. Bir `Splitter` kutusu kadar büyüktür: bir grupta
ona `Fill: true` ya da bir boyut verin.

### Pencerenin çubukları

| | |
|---|---|
| `ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{{"text": "Dosya", "items": []ui.P{{"text": "Kaydet", "shortcut": "Ctrl+S", "onClick": save}}}}})` | sayfanın içinde, pencerenin görünümünde; Mac'te Ctrl, Cmd'dir |
| `ui.ToolBar(ui.ToolBarProps{Items: []any{ui.P{"key": "save", "icon": "save", "hint": "Kaydet"}, "-"}, OnClick: fn})` | `"-"` bir çizgi; `"toggle": true` basılı kalır; `fn(bar *ui.Object, item ui.P)` |
| `ui.StatusBar(ui.StatusBarProps{MessageText: "Hazır", Items: []string{"UTF-8"}})` | `bar.Call("setMessageText", "Kaydedildi", "success", 3000)`, `bar.Call("setBusy", true)` |
| `ui.ContextMenu(ui.ContextMenuProps{Target: table, Items: []any{"Aç", "-", ui.P{"text": "Sil", "onClick": fn}}})` | hedefe sağ tıklamak onu açar |

Menüsü, araç çubuğu, yeri dolduran bir listesi ve durum çubuğu olan bir pencere, bunların `Gap: 0` olan bir
`VGroup`'udur; ortadaki parçaya `Fill: true` verilir (3. bölüm):

```go
package main

import (
	"fmt"

	ui "github.com/bug7a/ozdenui-go"
)

var status *ui.Object

func main() {
	ui.Run(start, ui.Window{Title: "Dosyalar", Width: 620, Height: 420})
}

func start(page *ui.Object) {
	ui.VGroup(ui.Props{Gap: 0}, func() {
		ui.MenuBar(ui.MenuBarProps{Menus: []ui.P{
			{"text": "Dosya", "items": []any{
				ui.P{"text": "Yeni", "shortcut": "Ctrl+N", "onClick": func() { tell("Yeni") }},
				"-",
				ui.P{"text": "Kapat", "onClick": ui.Close},
			}},
		}})
		ui.ToolBar(ui.ToolBarProps{OnClick: pressTool, Items: []any{
			ui.P{"key": "new", "icon": "new", "hint": "Yeni"},
			ui.P{"key": "save", "icon": "save", "hint": "Kaydet"},
		}})
		ui.ItemList(ui.ItemListProps{Fill: true, Unframed: true, OnChange: choose, Items: []ui.P{
			{"key": "report", "text": "rapor.pdf", "description": "120 KB"},
			{"key": "photo", "text": "foto.jpg", "description": "2.4 MB"},
			{"key": "notes", "text": "notlar.txt", "description": "4 KB"},
		}})
		status = ui.StatusBar(ui.StatusBarProps{MessageText: "Hazır", Items: []string{"3 dosya"}})
	})
}

func tell(text string) {
	status.Call("setMessageText", text, "info", 2000)
}

func pressTool(bar *ui.Object, item ui.P) { tell(fmt.Sprintf("%v basıldı", item["hint"])) }

func choose(list *ui.Object) { tell(fmt.Sprintf("Seçilen: %v", list.Get("selectedKeys"))) }
```

`Unframed: true` listenin kenarlığını ve yuvarlak köşelerini kaldırır: liste iki çubuk arasındaki yeri doldurur.
`examples/editor` ve `examples/gallery` bu türden daha büyük pencerelerdir.

### İlerleme ve adımlar

| | |
|---|---|
| `ui.ProgressBar(ui.ProgressBarProps{TitleText: "Kopyalanıyor", Value: 0})` | `bar.SetValue(40)` |
| `ui.LoadingScreen(ui.LoadingScreenProps{TitleText: "Uygulamam", MessageText: "Okunuyor..."})` | pencereyi örter; `loading.Call("setProgress", 0.4)`, `loading.Call("close")` |
| `ui.StepBar(ui.StepBarProps{Steps: []string{"Sepet", "Adres", "Ödeme"}, CurrentStep: 1})` | `bar.Call("next")`, `bar.Call("previous")`, `bar.Get("step")` |
| `ui.Tooltip(ui.TooltipProps{Target: button, HintText: "Dosyayı kaydeder"})` | fare üstündeyken bir ipucu; öbür biçimleri 6. bölümde |

### Görünümler: 8. bölüm

`GraphBox` (grafikler), `LogView`, `ImageView`, `MarkdownView`, `JsonView`, `JsonEditor`, `CodeEditor`, `Canvas`,
`ChatView` ve `WebView` veri, resim, kod, bir konuşma ve web sayfası gösterir (`JsonEditor` ve `CodeEditor` yazdırır
da): 8. bölüm.

### Bileşenlerin kendi metinleri

Bir bileşenin kendi yazdığı metinler (bir uyarı, bir diyaloğun düğmeleri, boş bir tablonun mesajı) İngilizcedir. Her
biri bir alanla değişir:

```go
ui.TextField(ui.TextFieldProps{TitleText: "Ad", Required: true, RequiredText: "Gerekli"}) // "Required" yerine
ui.Confirm("Dosya silinecek.", ui.DialogOptions{Title: "Silinsin mi?", ConfirmText: "Sil", CancelText: "Vazgeç", Danger: true})
ui.SmartTable(ui.SmartTableProps{NoDataFoundAlert: "Kayıt yok", AllTitlesText: "TÜMÜ"})
ui.DatePicker(ui.DatePickerProps{Date: "2026-11-20", Language: "tr"}) // ay ve gün adları Türkçe
ui.JsonEditor(ui.JsonEditorProps{Texts: ui.P{"addRow": "Satır", "search": "Satırlarda ara", "table": "Tablo"}})
```

Bir bileşenin metin alanları `Text` ile biter ya da ne olduklarını söyler (`Placeholder`, `EmptyText`); Props türünün
`go doc`'u her birini İngilizce varsayılan değeriyle gösterir. Bir `JsonEditor`'ın metinleri, kitin adlarıyla, tek bir
`Texts` map'indedir (`addRow`).

## Herhangi bir bileşen adıyla: ui.Component()

Props türü olan bir fonksiyon kolay yoldur. `ui.Component()` aynı bileşeni adından ve parametrelerinin, kitteki
adlarıyla bir `ui.P`'sinden yapar; programınızın çalışırken seçtiği bir ad için yol budur:

```go
ui.Gauge(ui.GaugeProps{Value: 72, SubText: "Disk"})
ui.Component("Gauge", ui.P{"value": 35, "subText": "Bellek"}) // aynı bileşen, adıyla
```

- Sayfa bir bileşeni, programınız onu ilk andığında, gerekenleriyle birlikte yükler: bir program yalnız kullandığının
  bedelini öder. `ui.Run()`'dan önce hiçbir şey listelenmez.
- Bileşenin tanımadığı bir `ui.P` adı terminalde, programınızın satırıyla söylenir.
- Props türünde alanı olmayan bir parametre: `Extra: ui.P{"name": value}`, kitteki adıyla.
- Her bileşen, her birinin bir örneğiyle: [katalog](https://bug7a.github.io/js-components/) ve
  `go doc github.com/bug7a/ozdenui-go`. Bir bileşenin içinde, nesne koymak ya da çağırmak için bir kutu:
  `obj.Part("name")`.

## Sonraki

[6. bölüm: Diyaloglar, mesajlar ve dosyalar](06-dialogs-and-files.md).
