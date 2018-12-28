package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
)

func (ds *DatastoreDB) Put(ctx context.Context, asset DSAsset, parent *datastore.Key) (*datastore.Key, error) {
	key := composeKey(asset, parent)

	key, err := ds.client.Put(ctx, key, asset)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (ds *DatastoreDB) Get(ctx context.Context, asset DSAsset) error {
	key := composeKey(asset, nil)

	err := ds.client.Get(ctx, key, asset)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return DSmissErr{AssetKey: key.String(), AssetNamespace: asset.GetDSNamespace()}
		}
		return err
	}
	return nil
}

func composeKey(asset DSAsset, parent *datastore.Key) *datastore.Key {
	kind := asset.GetDSKind()
	var key *datastore.Key
	if idKey, ok := asset.GetDSIDKey(); ok {
		key = datastore.IDKey(kind, idKey, parent)
	} else if nameKey, ok := asset.GetDSNameKey(); ok {
		key = datastore.NameKey(kind, nameKey, parent)
	} else {
		key = datastore.IncompleteKey(kind, parent)
	}
	key.Namespace = asset.GetDSNamespace()
	return key
}


type DSmissErr struct {
	AssetKey       string
	AssetNamespace string
}

func (e DSmissErr) Error() string {
	return fmt.Sprintf("asset not found in datastore - Key: %s - kind: %s", e.AssetKey, e.AssetNamespace)
}

