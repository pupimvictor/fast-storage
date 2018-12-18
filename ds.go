package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
)

func (ds *DatastoreDB) Put(ctx context.Context, asset DSAsset, parent *datastore.Key) (*datastore.Key, error) {
	key := composeKey(asset, parent)

	key, err := ds.client.Put(ctx, key, asset)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (ds *DatastoreDB) Get(ctx context.Context, asset DSAsset) (error) {
	key := composeKey(asset, nil)

	err := ds.client.Get(ctx, key, asset)
	if err != nil {
		return err
	}
	return nil
}

func composeKey(asset DSAsset, parent *datastore.Key) *datastore.Key {
	kind := asset.GetDSKind()
	var key *datastore.Key
	if idKey, ok := asset.GetIDKey(); ok {
		key = datastore.IDKey(kind, idKey, parent)
	} else if nameKey, ok := asset.GetNameKey(); ok {
		key = datastore.NameKey(kind, nameKey, parent)
	} else {
		key = datastore.IncompleteKey(kind, parent)
	}
	key.Namespace = asset.GetDSNamespace()
	return key
}