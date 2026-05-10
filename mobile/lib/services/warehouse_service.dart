import 'package:flutter/foundation.dart';
import 'package:uuid/uuid.dart';

import '../models/models.dart';
import 'database_service.dart';

class WarehouseService extends ChangeNotifier {
  final DatabaseService _dbService = DatabaseService();
  List<Warehouse> _warehouses = [];
  List<Shelf> _shelves = [];
  List<Box> _boxes = [];

  List<Warehouse> get warehouses => _warehouses;
  List<Shelf> get shelves => _shelves;
  List<Box> get boxes => _boxes;

  Future<void> loadData() async {
    _warehouses = await _dbService.getWarehouses();
    notifyListeners();
  }

  Future<void> loadShelves(String warehouseId) async {
    _shelves = await _dbService.getShelves(warehouseId);
    notifyListeners();
  }

  Future<void> loadBoxes(String shelfId) async {
    _boxes = await _dbService.getBoxes(shelfId);
    notifyListeners();
  }

  Future<void> addWarehouse(String name, String description, String location) async {
    final warehouse = Warehouse(
      id: const Uuid().v4(),
      name: name,
      description: description,
      location: location,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
    await _dbService.insertWarehouse(warehouse);
    _warehouses.add(warehouse);
    notifyListeners();
  }

  Future<void> updateWarehouse(Warehouse warehouse) async {
    final updatedWarehouse = warehouse.copyWith(updatedAt: DateTime.now());
    await _dbService.updateWarehouse(updatedWarehouse);
    final index = _warehouses.indexWhere((w) => w.id == warehouse.id);
    if (index != -1) {
      _warehouses[index] = updatedWarehouse;
    }
    notifyListeners();
  }

  Future<void> removeWarehouse(String id) async {
    await _dbService.deleteWarehouse(id);
    _warehouses.removeWhere((w) => w.id == id);
    notifyListeners();
  }

  Future<void> addShelf(String warehouseId, String name, int rows, int columns) async {
    final shelf = Shelf(
      id: const Uuid().v4(),
      warehouseId: warehouseId,
      name: name,
      rows: rows,
      columns: columns,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
    await _dbService.insertShelf(shelf);
    _shelves.add(shelf);
    notifyListeners();
  }

  Future<void> updateShelf(Shelf shelf) async {
    final updatedShelf = shelf.copyWith(updatedAt: DateTime.now());
    await _dbService.updateShelf(updatedShelf);
    final index = _shelves.indexWhere((s) => s.id == shelf.id);
    if (index != -1) {
      _shelves[index] = updatedShelf;
    }
    notifyListeners();
  }

  Future<void> removeShelf(String id) async {
    await _dbService.deleteShelf(id);
    _shelves.removeWhere((s) => s.id == id);
    notifyListeners();
  }

  Future<void> addBox(String shelfId, int row, int column, String code, String contents, int quantity) async {
    final box = Box(
      id: const Uuid().v4(),
      shelfId: shelfId,
      row: row,
      column: column,
      code: code,
      contents: contents,
      quantity: quantity,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
    await _dbService.insertBox(box);
    _boxes.add(box);
    notifyListeners();
  }

  Future<void> updateBox(Box box) async {
    final updatedBox = box.copyWith(updatedAt: DateTime.now());
    await _dbService.updateBox(updatedBox);
    final index = _boxes.indexWhere((b) => b.id == box.id);
    if (index != -1) {
      _boxes[index] = updatedBox;
    }
    notifyListeners();
  }

  Future<void> removeBox(String id) async {
    await _dbService.deleteBox(id);
    _boxes.removeWhere((b) => b.id == id);
    notifyListeners();
  }
}

extension WarehouseCopy on Warehouse {
  Warehouse copyWith({
    String? id,
    String? name,
    String? description,
    String? location,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Warehouse(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      location: location ?? this.location,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

extension ShelfCopy on Shelf {
  Shelf copyWith({
    String? id,
    String? warehouseId,
    String? name,
    int? rows,
    int? columns,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Shelf(
      id: id ?? this.id,
      warehouseId: warehouseId ?? this.warehouseId,
      name: name ?? this.name,
      rows: rows ?? this.rows,
      columns: columns ?? this.columns,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

extension BoxCopy on Box {
  Box copyWith({
    String? id,
    String? shelfId,
    int? row,
    int? column,
    String? code,
    String? contents,
    int? quantity,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Box(
      id: id ?? this.id,
      shelfId: shelfId ?? this.shelfId,
      row: row ?? this.row,
      column: column ?? this.column,
      code: code ?? this.code,
      contents: contents ?? this.contents,
      quantity: quantity ?? this.quantity,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
