package storage

const (
	TyUnknowAction = iota
	TyContentStorageAction
	TyHashStorageAction
	TyLinkStorageAction
	TyEncryptStorageAction
	TyEncryptShareStorageAction

	FuncNameQueryStorage = "QueryStorage"
)

const (
	OpCreate = int32(iota)
	OpAdd
)

const StorageX = "storage"

const Addr = "1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs"
