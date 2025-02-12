package extension_repo

import (
	"fmt"
	"seanime/internal/extension"
	"seanime/internal/util"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Anime Torrent provider
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *Repository) loadExternalAnimeTorrentProviderExtension(ext *extension.Extension) (err error) {
	defer util.HandlePanicInModuleWithError("extension_repo/loadExternalAnimeTorrentProviderExtension", &err)

	switch ext.Language {
	case extension.LanguageJavascript:
		err = r.loadExternalAnimeTorrentProviderExtensionJS(ext, extension.LanguageJavascript)
	case extension.LanguageTypescript:
		err = r.loadExternalAnimeTorrentProviderExtensionJS(ext, extension.LanguageTypescript)
	default:
		err = fmt.Errorf("unsupported language: %v", ext.Language)
	}

	if err != nil {
		return
	}

	return
}

func (r *Repository) loadExternalAnimeTorrentProviderExtensionJS(ext *extension.Extension, language extension.Language) error {
	provider, _, err := NewGojaAnimeTorrentProvider(ext, language, r.logger, r.gojaRuntimeManager)
	if err != nil {
		return err
	}

	// Add the extension to the map
	retExt := extension.NewAnimeTorrentProviderExtension(ext, provider)
	r.extensionBank.Set(ext.ID, retExt)
	return nil
}
