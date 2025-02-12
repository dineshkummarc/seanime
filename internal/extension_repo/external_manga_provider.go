package extension_repo

import (
	"seanime/internal/extension"
	"seanime/internal/util"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Manga
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *Repository) loadExternalMangaExtension(ext *extension.Extension) (err error) {
	defer util.HandlePanicInModuleWithError("extension_repo/loadExternalMangaExtension", &err)

	switch ext.Language {
	case extension.LanguageJavascript:
		err = r.loadExternalMangaExtensionJS(ext, extension.LanguageJavascript)
	case extension.LanguageTypescript:
		err = r.loadExternalMangaExtensionJS(ext, extension.LanguageTypescript)
	}

	if err != nil {
		return
	}

	return
}

func (r *Repository) loadExternalMangaExtensionJS(ext *extension.Extension, language extension.Language) error {
	provider, _, err := NewGojaMangaProvider(ext, language, r.logger, r.gojaRuntimeManager)
	if err != nil {
		return err
	}

	// Add the extension to the map
	retExt := extension.NewMangaProviderExtension(ext, provider)
	r.extensionBank.Set(ext.ID, retExt)
	return nil
}
