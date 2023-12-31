package mongodb

import (
	"SilentPaymentAppBackend/src/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// todo add unique lock
//  db.members.createIndex( { groupNumber: 1, lastname: 1, firstname: 1 }, { unique: true } )

func CreateIndices() {
	common.InfoLogger.Println("creating database indices")
	CreateIndexTransactions()
	CreateIndexCFilters()
	CreateIndexTweaks()
	CreateIndexUTXOs()
	CreateIndexSpentTXOs()
	CreateIndexHeaders()
	common.InfoLogger.Println("created database indices")
}

func CreateIndexTransactions() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transactions").Collection("taproot_transactions")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"txid": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func CreateIndexCFilters() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("filters").Collection("general")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			// in rare case counting is off we can then reindex from local DB data
			"blockheader": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func CreateIndexUTXOs() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transaction_outputs").Collection("unspent")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "txid", Value: 1},
			{Key: "vout", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func CreateIndexSpentTXOs() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transaction_outputs").Collection("spent")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "txid", Value: 1},
			{Key: "vout", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func CreateIndexTweaks() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("tweak_data").Collection("tweaks")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"txid": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func CreateIndexHeaders() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("headers").Collection("headers")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"block_hash": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	nameIndex, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// will panic because it only runs on startup and should be executed
		panic(err)
	}
	common.DebugLogger.Println("Created Index with name:", nameIndex)
}

func SaveTransactionDetails(transaction *common.Transaction) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transactions").Collection("taproot_transactions")

	result, err := coll.InsertOne(context.TODO(), transaction)
	if err != nil {
		common.ErrorLogger.Println(err)
		//panic(err)
		return
	}

	log.Printf("Transaction inserted with ID: %s\n", result.InsertedID)
}

func SaveFilter(filter *common.Filter) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			common.ErrorLogger.Println(err)
			return
		}
	}()

	coll := client.Database("filters").Collection("general")

	result, err := coll.InsertOne(context.TODO(), filter)
	if err != nil {
		//todo don't log duplicate keys as error but rather as debug
		common.ErrorLogger.Println(err)
		return
	}

	log.Println("Filter inserted", "ID", result.InsertedID)
}

func SaveFilterTaproot(filter *common.Filter) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			common.ErrorLogger.Println(err)
			return
		}
	}()

	coll := client.Database("filters").Collection("taproot")

	result, err := coll.InsertOne(context.TODO(), filter)
	if err != nil {
		//todo don't log duplicate keys as error but rather as debug
		common.ErrorLogger.Println(err)
		return
	}

	log.Println("Taproot Filter inserted", "ID", result.InsertedID)
}

func SaveLightUTXO(utxo *common.LightUTXO) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transaction_outputs").Collection("unspent")

	result, err := coll.InsertOne(context.TODO(), utxo)
	if err != nil {
		common.ErrorLogger.Println(err)
		//panic(err)
		return
	}

	log.Printf("UTXO inserted with ID: %s\n", result.InsertedID)
}

func SaveTweakData(tweak *common.TweakData) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("tweak_data").Collection("tweaks")

	result, err := coll.InsertOne(context.TODO(), tweak)
	if err != nil {
		common.ErrorLogger.Println(err)
		//panic(err)
		return
	}

	log.Printf("Tweak inserted with ID: %s\n", result.InsertedID)
}

func SaveSpentUTXO(utxo *common.SpentUTXO) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transaction_outputs").Collection("spent")

	result, err := coll.InsertOne(context.TODO(), utxo)
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	log.Printf("Spent Transaction output inserted with ID: %s\n", result.InsertedID)
}

func SaveBulkHeaders(headers []*common.Header) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("headers").Collection("headers")

	// Convert []*common.Header to []interface{}
	var interfaceHeaders []interface{}
	for _, header := range headers {
		interfaceHeaders = append(interfaceHeaders, header)
	}

	result, err := coll.InsertMany(context.TODO(), interfaceHeaders)
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	log.Printf("Bulk inserted %d new headers\n", len(result.InsertedIDs))
}

func RetrieveLastHeader() *common.Header {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("headers").Collection("headers")
	var result common.Header
	filter := bson.D{}                                                   // no filter, get all documents
	optionsQuery := options.FindOne().SetSort(bson.D{{"timestamp", -1}}) // sort by timestamp in descending order

	err = coll.FindOne(context.TODO(), filter, optionsQuery).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No documents found!")
			return nil
		}
		panic(err)
	}

	return &result
}

func RetrieveAllHeaders() []*common.Header {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("headers").Collection("headers")

	filter := bson.D{}                                               // no filter, get all documents
	optionsQuery := options.Find().SetSort(bson.D{{"timestamp", 1}}) // sort by timestamp in descending order

	cursor, err := coll.Find(context.TODO(), filter, optionsQuery)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	var results []*common.Header
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func RetrieveTransactionsByHeight(blockHeight uint32) []*common.Transaction {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transactions").Collection("taproot_transactions")
	filter := bson.D{{"status.blockheight", blockHeight}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	var results []*common.Transaction
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func RetrieveLightUTXOsByHeight(blockHeight uint32) []*common.LightUTXO {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("transaction_outputs").Collection("unspent")
	filter := bson.D{{"blockheight", blockHeight}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	var results []*common.LightUTXO
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func RetrieveSpentUTXOsByHeight(blockHeight uint32) []*common.SpentUTXO {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("transaction_outputs").Collection("spent")
	filter := bson.D{{"blockheight", blockHeight}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	var results []*common.SpentUTXO
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func RetrieveCFilterByHeight(blockHeight uint32) *common.Filter {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("filters").Collection("general")
	filter := bson.D{{"blockheight", blockHeight}}

	result := coll.FindOne(context.TODO(), filter)
	var cFilter common.Filter

	err = result.Decode(&cFilter)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	return &cFilter
}

func RetrieveCFilterByHeightTaproot(blockHeight uint32) *common.Filter {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("filters").Collection("taproot")
	filter := bson.D{{"blockheight", blockHeight}}

	result := coll.FindOne(context.TODO(), filter)
	var cFilter common.Filter

	err = result.Decode(&cFilter)
	if err != nil {
		common.ErrorLogger.Println(err)
		return nil
	}

	return &cFilter
}

func RetrieveTweakDataByHeight(blockHeight uint32) []*common.TweakData {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("tweak_data").Collection("tweaks")
	filter := bson.D{{"blockheight", blockHeight}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		common.ErrorLogger.Println(err)
	}

	var results []*common.TweakData
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}

func DeleteLightUTXOByTxIndex(txId string, vout uint32) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(common.MongoDBURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("transaction_outputs").Collection("unspent")
	filter := bson.D{{"txid", txId}, {"vout", vout}}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		common.ErrorLogger.Println(err)
		//panic(err)
		return
	}
	common.DebugLogger.Printf("Deleted %d LightUTXOs\n", result.DeletedCount)

	err = chainedTweakDeletion(client, txId)
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}
}

// chained deletion of tweak data if no more utxos with a certain txid are left
func chainedTweakDeletion(client *mongo.Client, txId string) error {
	coll := client.Database("tweak_data").Collection("tweaks")
	filter := bson.D{{"txid", txId}}
	result := coll.FindOne(context.TODO(), filter)

	var utxo common.LightUTXO

	err := result.Decode(&utxo)
	if err != nil && err.Error() != "mongo: no documents in result" {
		common.ErrorLogger.Println(err)
		return err
	}

	// if no match was found
	if utxo.Txid == "" {
		var resultDelete *mongo.DeleteResult
		resultDelete, err = coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			common.ErrorLogger.Println(err)
			return err
		}

		common.DebugLogger.Printf("Deleted %d tweak data for %s\n", resultDelete.DeletedCount, txId)
		return err
	}
	return nil
}
