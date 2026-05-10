import 'dart:convert';
import 'package:http/http.dart' as http;

import '../models/models.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:8080/api';
  static const String deviceId = 'mobile_device_001';

  static Future<List<Warehouse>> fetchWarehouses() async {
    final response = await http.get(Uri.parse('$baseUrl/warehouses'));
    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Warehouse.fromJson(json)).toList();
    }
    throw Exception('Failed to fetch warehouses');
  }

  static Future<Warehouse> createWarehouse(Warehouse warehouse) async {
    final response = await http.post(
      Uri.parse('$baseUrl/warehouses'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(warehouse.toJson()),
    );
    if (response.statusCode == 201) {
      return Warehouse.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to create warehouse');
  }

  static Future<Warehouse> updateWarehouse(Warehouse warehouse) async {
    final response = await http.put(
      Uri.parse('$baseUrl/warehouses/${warehouse.id}'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(warehouse.toJson()),
    );
    if (response.statusCode == 200) {
      return Warehouse.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to update warehouse');
  }

  static Future<void> deleteWarehouse(String id) async {
    final response = await http.delete(Uri.parse('$baseUrl/warehouses/$id'));
    if (response.statusCode != 204) {
      throw Exception('Failed to delete warehouse');
    }
  }

  static Future<List<Shelf>> fetchShelves(String warehouseId) async {
    final response = await http.get(Uri.parse('$baseUrl/shelves/warehouse/$warehouseId'));
    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Shelf.fromJson(json)).toList();
    }
    throw Exception('Failed to fetch shelves');
  }

  static Future<Shelf> createShelf(Shelf shelf) async {
    final response = await http.post(
      Uri.parse('$baseUrl/shelves'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(shelf.toJson()),
    );
    if (response.statusCode == 201) {
      return Shelf.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to create shelf');
  }

  static Future<Shelf> updateShelf(Shelf shelf) async {
    final response = await http.put(
      Uri.parse('$baseUrl/shelves/${shelf.id}'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(shelf.toJson()),
    );
    if (response.statusCode == 200) {
      return Shelf.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to update shelf');
  }

  static Future<void> deleteShelf(String id) async {
    final response = await http.delete(Uri.parse('$baseUrl/shelves/$id'));
    if (response.statusCode != 204) {
      throw Exception('Failed to delete shelf');
    }
  }

  static Future<List<Box>> fetchBoxes(String shelfId) async {
    final response = await http.get(Uri.parse('$baseUrl/boxes/shelf/$shelfId'));
    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((json) => Box.fromJson(json)).toList();
    }
    throw Exception('Failed to fetch boxes');
  }

  static Future<Box> createBox(Box box) async {
    final response = await http.post(
      Uri.parse('$baseUrl/boxes'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(box.toJson()),
    );
    if (response.statusCode == 201) {
      return Box.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to create box');
  }

  static Future<Box> updateBox(Box box) async {
    final response = await http.put(
      Uri.parse('$baseUrl/boxes/${box.id}'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(box.toJson()),
    );
    if (response.statusCode == 200) {
      return Box.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to update box');
  }

  static Future<void> deleteBox(String id) async {
    final response = await http.delete(Uri.parse('$baseUrl/boxes/$id'));
    if (response.statusCode != 204) {
      throw Exception('Failed to delete box');
    }
  }

  static Future<void> syncData(List<SyncAction> actions) async {
    final response = await http.post(
      Uri.parse('$baseUrl/sync'),
      headers: {
        'Content-Type': 'application/json',
        'X-Device-ID': deviceId,
      },
      body: json.encode(actions.map((a) => a.toJson()).toList()),
    );
    if (response.statusCode != 200) {
      throw Exception('Failed to sync data');
    }
  }
}
