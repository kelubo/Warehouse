import 'package:flutter/foundation.dart';

class Warehouse {
  final String id;
  final String name;
  final String description;
  final String location;
  final DateTime createdAt;
  final DateTime updatedAt;

  Warehouse({
    required this.id,
    required this.name,
    this.description = '',
    this.location = '',
    required this.createdAt,
    required this.updatedAt,
  });

  factory Warehouse.fromJson(Map<String, dynamic> json) {
    return Warehouse(
      id: json['id'],
      name: json['name'],
      description: json['description'] ?? '',
      location: json['location'] ?? '',
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'location': location,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class Shelf {
  final String id;
  final String warehouseId;
  final String name;
  final int rows;
  final int columns;
  final DateTime createdAt;
  final DateTime updatedAt;

  Shelf({
    required this.id,
    required this.warehouseId,
    required this.name,
    this.rows = 5,
    this.columns = 5,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Shelf.fromJson(Map<String, dynamic> json) {
    return Shelf(
      id: json['id'],
      warehouseId: json['warehouse_id'],
      name: json['name'],
      rows: json['rows'] ?? 5,
      columns: json['columns'] ?? 5,
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'warehouse_id': warehouseId,
      'name': name,
      'rows': rows,
      'columns': columns,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class Box {
  final String id;
  final String shelfId;
  final int row;
  final int column;
  final String code;
  final String contents;
  final int quantity;
  final DateTime createdAt;
  final DateTime updatedAt;

  Box({
    required this.id,
    required this.shelfId,
    required this.row,
    required this.column,
    required this.code,
    this.contents = '',
    this.quantity = 0,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Box.fromJson(Map<String, dynamic> json) {
    return Box(
      id: json['id'],
      shelfId: json['shelf_id'],
      row: json['row'],
      column: json['column'],
      code: json['code'],
      contents: json['contents'] ?? '',
      quantity: json['quantity'] ?? 0,
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'shelf_id': shelfId,
      'row': row,
      'column': column,
      'code': code,
      'contents': contents,
      'quantity': quantity,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class SyncAction {
  final String modelType;
  final String modelId;
  final String action;
  final Map<String, dynamic> data;

  SyncAction({
    required this.modelType,
    required this.modelId,
    required this.action,
    required this.data,
  });

  Map<String, dynamic> toJson() {
    return {
      'model_type': modelType,
      'model_id': modelId,
      'action': action,
      'data': data,
    };
  }
}

class SyncRecord {
  final String id;
  final String deviceId;
  final String modelType;
  final String modelId;
  final String action;
  final String data;
  final DateTime syncedAt;
  final DateTime createdAt;

  SyncRecord({
    required this.id,
    required this.deviceId,
    required this.modelType,
    required this.modelId,
    required this.action,
    required this.data,
    required this.syncedAt,
    required this.createdAt,
  });
}
