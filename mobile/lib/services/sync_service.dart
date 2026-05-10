import 'package:flutter/foundation.dart';

import '../models/models.dart';
import 'api_service.dart';
import 'database_service.dart';

class SyncService extends ChangeNotifier {
  final DatabaseService _dbService = DatabaseService();
  bool _isSyncing = false;
  bool _hasPendingChanges = false;
  String _lastSyncTime = '';

  bool get isSyncing => _isSyncing;
  bool get hasPendingChanges => _hasPendingChanges;
  String get lastSyncTime => _lastSyncTime;

  Future<void> checkPendingChanges() async {
    final actions = await _dbService.getPendingSyncActions();
    _hasPendingChanges = actions.isNotEmpty;
    notifyListeners();
  }

  Future<void> sync() async {
    _isSyncing = true;
    notifyListeners();

    try {
      final pendingActions = await _dbService.getPendingSyncActions();
      
      if (pendingActions.isNotEmpty) {
        await ApiService.syncData(pendingActions);
        await _dbService.clearSyncActions();
      }

      final warehouses = await ApiService.fetchWarehouses();
      List<Shelf> allShelves = [];
      List<Box> allBoxes = [];

      for (var w in warehouses) {
        final shelves = await ApiService.fetchShelves(w.id);
        allShelves.addAll(shelves);
        for (var s in shelves) {
          final boxes = await ApiService.fetchBoxes(s.id);
          allBoxes.addAll(boxes);
        }
      }

      await _dbService.syncFromServer(warehouses, allShelves, allBoxes);
      
      _lastSyncTime = DateTime.now().toString();
      _hasPendingChanges = false;
    } catch (e) {
      print('Sync failed: $e');
      rethrow;
    } finally {
      _isSyncing = false;
      notifyListeners();
    }
  }

  Future<void> syncFromServer() async {
    _isSyncing = true;
    notifyListeners();

    try {
      final warehouses = await ApiService.fetchWarehouses();
      List<Shelf> allShelves = [];
      List<Box> allBoxes = [];

      for (var w in warehouses) {
        final shelves = await ApiService.fetchShelves(w.id);
        allShelves.addAll(shelves);
        for (var s in shelves) {
          final boxes = await ApiService.fetchBoxes(s.id);
          allBoxes.addAll(boxes);
        }
      }

      await _dbService.syncFromServer(warehouses, allShelves, allBoxes);
      _lastSyncTime = DateTime.now().toString();
    } catch (e) {
      print('Sync from server failed: $e');
      rethrow;
    } finally {
      _isSyncing = false;
      notifyListeners();
    }
  }
}
