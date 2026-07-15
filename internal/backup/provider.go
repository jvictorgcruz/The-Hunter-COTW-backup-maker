package backup

type BackupProvider interface {
	ID() string
	IsConfigured() bool
	Send(localZipPath string, filename string) error
	Cleanup(limit int) error
}

var providers = make(map[string]BackupProvider)

func RegisterProvider(p BackupProvider) {
	providers[p.ID()] = p
}

func GetProvider(id string) BackupProvider {
	return providers[id]
}

func ListProviders() []BackupProvider {
	var list []BackupProvider
	for _, p := range providers {
		list = append(list, p)
	}
	return list
}
