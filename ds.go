package faststorage

import (
	"cloud.google.com/go/datastore"
	"context"
)

func (ds *DatastoreDB) Put(ctx context.Context, asset Asset, parent *datastore.Key) (*datastore.Key, error) {
	var key *datastore.Key
	if idKey, ok := asset.GetIDKey(); ok{
		key = datastore.IDKey(asset.GetDSKind(), idKey, parent)
	} else if nameKey, ok := asset.GetNameKey(); ok {
		key = datastore.NameKey(asset.GetDSKind(), nameKey, parent)
	} else {
		key = datastore.IncompleteKey(asset.GetDSKind(), parent)
	}
	key.Namespace = asset.GetDSNamespace()

	key, err := ds.client.Put(ctx, key, asset)
	if err != nil {
		return nil, err
	}
	return key, nil
}

