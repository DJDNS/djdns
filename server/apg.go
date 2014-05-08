package server

type AliasPageGetter struct {
	Aliases map[string]string
	Child   PageGetter
}

func NewAliasPageGetter(child PageGetter) AliasPageGetter {
	return AliasPageGetter{
		Aliases: make(map[string]string),
		Child:   child,
	}
}

func (apg AliasPageGetter) GetPage(url string, ab Aborter) (Page, error) {
	transformed, ok := apg.Aliases[url]
	if ok {
		url = transformed
	}
	return apg.Child.GetPage(url, ab)
}
