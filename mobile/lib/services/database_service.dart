import 'dart:convert';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

import '../models/models.dart';

class DatabaseService {
  static Database? _database;

  Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await initDatabase();
    return _database!;
  }

  Future<Database> initDatabase() async {
    String path = join(await getDatabasesPath(), 'warehouse.db');
    return await openDatabase(
      path,
      version: 1,
      onCreate: (db, version) async {
        await db.execute('''
          CREATE TABLE warehouses (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            location TEXT,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
          )
        ''');
        await db.execute('''
          CREATE TABLE shelves (
            id TEXT PRIMARY KEY,
            warehouse_id TEXT NOT NULL,
            name TEXT NOT NULL,
            rows INTEGER NOT NULL DEFAULT 5,
            columns INTEGER NOT NULL DEFAULT 5,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL,
            FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
          )
        ''');
        await db.execute('''
          CREATE TABLE boxes (
            id TEXT PRIMARY KEY,
            shelf_id TEXT NOT NULL,
            row INTEGER NOT NULL,
            column INTEGER NOT NULL,
            code TEXT NOT NULL UNIQUE,
            contents TEXT,
            quantity INTEGER NOT NULL DEFAULT 0,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL,
            FOREIGN KEY (shelf_id) REFERENCES shelves(id)
          )
        ''');
        await db.execute('''
          CREATE TABLE sync_actions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            model_type TEXT NOT NULL,
            model_id TEXT NOT NULL,
            action TEXT NOT NULL,
            data TEXT NOT NULL,
            created_at TEXT NOT NULL
          )
        ''');
      },
    );
  }

  Future<void> insertWarehouse(Warehouse warehouse) async {
    final db = await database;
    await db.insert('warehouses', warehouse.toJson());
    await _saveSyncAction('warehouse', warehouse.id, 'create', warehouse.toJson());
  }

  Future<List<Warehouse>> getWarehouses() async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query('warehouses');
    return List.generate(maps.length, (i) => Warehouse.fromJson(maps[i]));
  }

  Future<void> updateWarehouse(Warehouse warehouse) async {
    final db = await database;
    await db.update(
      'warehouses',
      warehouse.toJson(),
      where: 'id = ?',
      whereArgs: [warehouse.id],
    );
    await _saveSyncAction('warehouse', warehouse.id, 'update', warehouse.toJson());
  }

  Future<void> deleteWarehouse(String id) async {
    final db = await database;
    await db.delete('warehouses', where: 'id = ?', whereArgs: [id]);
    await _saveSyncAction('warehouse', id, 'delete', {});
  }

  Future<void> insertShelf(Shelf shelf) async {
    final db = await database;
    await db.insert('shelves', shelf.toJson());
    await _saveSyncAction('shelf', shelf.id, 'create', shelf.toJson());
  }

  Future<List<Shelf>> getShelves(String warehouseId) async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query(
      'shelves',
      where: 'warehouse_id = ?',
      whereArgs: [warehouseId],
    );
    return List.generate(maps.length, (i) => Shelf.fromJson(maps[i]));
  }

  Future<void> updateShelf(Shelf shelf) async {
    final db = await database;
    await db.update(
      'shelves',
      shelf.toJson(),
      where: 'id = ?',
      whereArgs: [shelf.id],
    );
    await _saveSyncAction('shelf', shelf.id, 'update', shelf.toJson());
  }

  Future<void> deleteShelf(String id) async {
    final db = await database;
    await db.delete('shelves', where: 'id = ?', whereArgs: [id]);
    await _saveSyncAction('shelf', id, 'delete', {});
  }

  Future<void> insertBox(Box box) async {
    final db = await database;
    await db.insert('boxes', box.toJson());
    await _saveSyncAction('box', box.id, 'create', box.toJson());
  }

  Future<List<Box>> getBoxes(String shelfId) async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query(
      'boxes',
      where: 'shelf_id = ?',
      whereArgs: [shelfId],
    );
    return List.generate(maps.length, (i) => Box.fromJson(maps[i]));
  }

  Future<void> updateBox(Box box) async {
    final db = await database;
    await db.update(
      'boxes',
      box.toJson(),
      where: 'id = ?',
      whereArgs: [box.id],
    );
    await _saveSyncAction('box', box.id, 'update', box.toJson());
  }

  Future<void> deleteBox(String id) async {
    final db = await database;
    await db.delete('boxes', where: 'id = ?', whereArgs: [id]);
    await _saveSyncAction('box', id, 'delete', {});
  }

  Future<void> _saveSyncAction(String modelType, String modelId, String action, Map<String, dynamic> data) async {
    final db = await database;
    await db.insert('sync_actions', {
      'model_type': modelType,
      'model_id': modelId,
      'action': action,
      'data': json.encode(data),
      'created_at': DateTime.now().toIso8601String(),
    });
  }

  Future<List<SyncAction>> getPendingSyncActions() async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query('sync_actions');
    return List.generate(maps.length, (i) {
      return SyncAction(
        modelType: maps[i]['model_type'],
        modelId: maps[i]['model_id'],
        action: maps[i]['action'],
        data: json.decode(maps[i]['data']),
      );
    });
  }

  Future<void> clearSyncActions() async {
    final db = await database;
    await db.delete('sync_actions');
  }

  Future<void> syncFromServer(List<Warehouse> warehouses, List<Shelf> shelves, List<Box> boxes) async {
    final db = await database;
    await db.transaction((txn) async {
      await txn.delete('boxes');
      await txn.delete('shelves');
      await txn.delete('warehouses');

      for (var w in warehouses) {
        await txn.insert('warehouses', w.toJson());
      }
      for (var s in shelves) {
        await txn.insert('shelves', s.toJson());
      }
      for (var b in boxes) {
        await txn.insert('boxes', b.toJson());
      }
    });
  }
}
