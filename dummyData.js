use dummyStore;

var bulk = db.store.initializeUnorderedBulkOp();

bulk.insert( { _id: 1, second: 20, value: 2546 });
bulk.insert( { _id: 2, second: 25, value: 6537 });
bulk.insert( { _id: 3, second: 10, value: 4568 });
bulk.insert( { _id: 4, second: 7, value: 5489 });
bulk.insert( { _id: 5, second: 35, value: 9965 });
bulk.insert( { _id: 6, second: 66, value: 5364 });
bulk.insert( { _id: 7, second: 65, value: 2156 });
bulk.insert( { _id: 8, second: 20, value: 2486 });
bulk.insert( { _id: 9, second: 43, value: 3248 });
bulk.insert( { _id: 10, second: 76, value: 9999 });
bulk.insert( { _id: 11, second: 52, value: 5542 });
bulk.insert( { _id: 12, second: 94, value: 5679 });

bulk.execute();